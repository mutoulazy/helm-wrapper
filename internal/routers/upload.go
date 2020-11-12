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

// @Summary 		上传chart脚本到服务器
// @Description 	上传chart脚本压缩包到服务器
// @Tags			Chart
// @Param 			chart formData file true "chart文件"
// @Success 		200 {object} app.ResponseBody
// @Router 			/api/v1/charts/upload [post]
func UploadChart(c *gin.Context) {
	response := app.NewResponse(c)
	file, header, err := c.Request.FormFile("chart")
	if err != nil {
		global.Logger.Error(c, "no chart file")
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	filename := header.Filename
	t := strings.Split(filename, ".")
	if t[len(t)-1] != "tgz" {
		global.Logger.Error(c, "chart file suffix must .tgz")
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

// @Summary 		查询上传chart脚本
// @Description 	查询已经上传chart脚本列表
// @Tags			Chart
// @Success 		200 {object} app.ResponseBody
// @Router 			/api/v1/charts/upload [get]
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
