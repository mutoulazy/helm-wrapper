package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"gopkg.in/natefinch/lumberjack.v2"
	"helm-wrapper/global"
	"helm-wrapper/internal/routers"
	"helm-wrapper/internal/routers/api/v1"
	"helm-wrapper/pkg/logger"
	"io/ioutil"
	"log"
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
	isVersion    bool
	listenHost   string
	listenPort   string
	config       string
	buildTime    string
	buildVersion string
	gitCommitID  string
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
	err = setupLogger()
	if err != nil {
		glog.Fatalf("加载日志文件失败: %v", err)
	}
}

// @title helm3代理
// @version v1.0
// @description golang helm-wrapper
// @termsOfService https://github.com/mutoulazy/helm-wrapper
func main() {
	// 输出编译信息
	if isVersion {
		fmt.Printf("build_time: %s\n", buildTime)
		fmt.Printf("build_version: %s\n", buildVersion)
		fmt.Printf("git_commit_id: %s\n", gitCommitID)
		return
	}

	gin.SetMode(global.MyHelmConfig.RunMode)

	// router
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome helm wrapper server")
	})

	// register router
	routers.RegisterRouter(router)

	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%s", listenHost, listenPort),
		Handler:        router,
		ReadTimeout:    global.MyHelmConfig.ReadTimeout * time.Second,
		WriteTimeout:   global.MyHelmConfig.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
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
	pflag.CommandLine.BoolVar(&isVersion, "version", false, "编译信息")
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

func setupLogger() error {
	fileName := global.MyHelmConfig.LogSavePath + "/" + global.MyHelmConfig.LogFileName + global.MyHelmConfig.LogFileExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}
