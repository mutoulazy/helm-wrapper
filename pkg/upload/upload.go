package upload

import (
	"fmt"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadChart(c *gin.Context) {
	file, header, err := c.Request.FormFile("chart")
	if err != nil {
		app.RespErr(c, err)
		return
	}

	filename := header.Filename
	t := strings.Split(filename, ".")
	if t[len(t)-1] != "tgz" {
		app.RespErr(c, fmt.Errorf("chart file suffix must .tgz"))
		return
	}

	out, err := os.Create(global.MyHelmConfig.UploadPath + "/" + filename)
	if err != nil {
		app.RespErr(c, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		app.RespErr(c, err)
		return
	}

	app.RespOK(c, nil)
}

func ListUploadedCharts(c *gin.Context) {
	charts := []string{}
	files, err := ioutil.ReadDir(global.MyHelmConfig.UploadPath)
	if err != nil {
		app.RespErr(c, err)
		return
	}
	for _, f := range files {
		t := strings.Split(f.Name(), ".")
		if t[len(t)-1] == "tgz" {
			charts = append(charts, f.Name())
		}
	}

	app.RespOK(c, charts)
}
