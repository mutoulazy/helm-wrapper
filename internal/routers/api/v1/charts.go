package v1

import (
	"fmt"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
	"helm-wrapper/pkg/errcode"
	"strings"

	"github.com/gin-gonic/gin"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

type Chart struct {
}

func NewChart() Chart {
	return Chart{}
}

var readmeFileNames = []string{"readme.md", "readme.txt", "readme"}

type file struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func findReadme(files []*chart.File) (file *chart.File) {
	for _, file := range files {
		for _, n := range readmeFileNames {
			if strings.EqualFold(file.Name, n) {
				return file
			}
		}
	}
	return nil
}

// @Summary 		获取chart详细信息
// @Description 	根据chart名称，获取chart的readme、values、chart、template信息
// @Tags			Chart
// @Param 			chart query string true "chart名称"
// @Param   		version query string false "chart版本"
// @Param   		info query string false "Enums(all, readme, values, chart、template)"
// @Success 		200 {object} app.ResponseBody
// @Router 			/api/v1/charts [get]
func (chart Chart) ShowChartInfo(c *gin.Context) {
	response := app.NewResponse(c)
	name := c.Query("chart")
	if name == "" {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails("chart name can not be empty"))
		return
	}
	// local charts with abs path *.tgz
	splitChart := strings.Split(name, ".")
	if splitChart[len(splitChart)-1] == "tgz" {
		name = global.MyHelmConfig.UploadPath + "/" + name
	}

	info := c.Query("info") // readme, values, chart
	version := c.Query("version")

	client := action.NewShow(action.ShowAll)
	client.Version = version
	if info == string(action.ShowChart) {
		client.OutputFormat = action.ShowChart
	} else if info == string(action.ShowReadme) {
		client.OutputFormat = action.ShowReadme
	} else if info == string(action.ShowValues) {
		client.OutputFormat = action.ShowValues
	} else {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(fmt.Sprintf("bad info %s, chart info only support readme/values/chart", info)))
		return
	}

	cp, err := client.ChartPathOptions.LocateChart(name, global.HelmClientSettings)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorShowChartInfoFail.WithDetails(err.Error()))
		return
	}

	chrt, err := loader.Load(cp)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorShowChartInfoFail.WithDetails(err.Error()))
		return
	}

	if client.OutputFormat == action.ShowChart {
		response.ToResponse(gin.H{"chart": chrt.Metadata})
		return
	}
	if client.OutputFormat == action.ShowValues {
		values := make([]*file, 0, len(chrt.Raw))
		for _, v := range chrt.Raw {
			values = append(values, &file{
				Name: v.Name,
				Data: string(v.Data),
			})
		}
		response.ToResponse(gin.H{"values": values})
		return
	}
	if client.OutputFormat == action.ShowReadme {
		response.ToResponse(gin.H{"readme": string(findReadme(chrt.Files).Data)})
		return
	}
}
