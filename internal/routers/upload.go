package routers

import (
	"fmt"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
	"helm-wrapper/pkg/errcode"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadChart(c *gin.Context) {
	response := app.NewResponse(c)
	file, header, err := c.Request.FormFile("chart")
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	filename := header.Filename
	t := strings.Split(filename, ".")
	if t[len(t)-1] != "tgz" {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(fmt.Sprintf("chart file suffix must .tgz")))
		return
	}

	out, err := os.Create(global.MyHelmConfig.UploadPath + "/" + filename)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{"msg": "success"})
	return
}

func ListUploadedCharts(c *gin.Context) {
	response := app.NewResponse(c)
	charts := []string{}
	files, err := ioutil.ReadDir(global.MyHelmConfig.UploadPath)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorDListUploadedChartsFail.WithDetails(err.Error()))
		return
	}
	for _, f := range files {
		t := strings.Split(f.Name(), ".")
		if t[len(t)-1] == "tgz" {
			charts = append(charts, f.Name())
		}
	}

	response.ToResponse(gin.H{"charts": charts})
	return
}
