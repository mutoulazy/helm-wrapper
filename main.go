package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"helm-wrapper/global"
	"helm-wrapper/internal/routers"
	"helm-wrapper/internal/routers/api/v1"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"sigs.k8s.io/yaml"
)

var (
	listenHost string
	listenPort string
	config     string
)

func init() {
	err := setupFlag()
	if err != nil {
		glog.Fatalf("加载Flag失败: %v", err)
	}
	err = setupConfig()
	if err != nil {
		glog.Fatalf("加载配置文件失败: %v", err)
	}
}

func main() {
	// router
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome helm wrapper server")
	})

	// register router
	routers.RegisterRouter(router)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", listenHost, listenPort),
		Handler:      router,
		ReadTimeout:  global.MyHelmConfig.ReadTimeout * time.Second,
		WriteTimeout: global.MyHelmConfig.WriteTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.Fatalf("helm wrapper listen err: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	glog.Infoln("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		glog.Fatal("Server forced to shutdown:", err)
	}

	glog.Infoln("Server exiting")
}

func setupFlag() error {
	err := flag.Set("logtostderr", "true")
	if err != nil {
		glog.Fatalf("setupFlag set logtostderr err: %v", err)
	}
	pflag.CommandLine.StringVar(&listenHost, "addr", "0.0.0.0", "server listen addr")
	pflag.CommandLine.StringVar(&listenPort, "port", "8080", "server listen port")
	pflag.CommandLine.StringVar(&config, "config", "config/config.yaml", "helm wrapper config")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	global.HelmClientSettings.AddFlags(pflag.CommandLine)
	pflag.Parse()
	defer glog.Flush()

	return err
}

func setupConfig() error {
	configBody, err := ioutil.ReadFile(config)
	if err != nil {
		glog.Fatalf("setupConfig ReadFile config file err: %v", err)
	}
	err = yaml.Unmarshal(configBody, global.MyHelmConfig)
	if err != nil {
		glog.Fatalf("setupConfig Unmarshal config err: %v", err)
	}

	// 初始化上传路径
	if global.MyHelmConfig.UploadPath == "" {
		global.MyHelmConfig.UploadPath = global.DefaultUploadPath
	} else {
		if !filepath.IsAbs(global.MyHelmConfig.UploadPath) {
			glog.Fatalln("charts upload path is not absolute")
		}
	}
	_, err = os.Stat(global.MyHelmConfig.UploadPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(global.MyHelmConfig.UploadPath, 0755)
			if err != nil {
				glog.Fatalf("setupConfig Mkdir upload err: %v", err)
			}
		} else {
			glog.Fatalf("setupConfig upload filepath stat err: %v", err)
		}
	}

	// 初始化chart repo
	for _, c := range global.MyHelmConfig.HelmRepos {
		err = v1.NewRepository().InitRepository(c)
		if err != nil {
			glog.Fatalf("setupConfig initRepository err: %v", err)
		}
	}

	return err
}
