package v1

import (
	"fmt"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
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

func (chart Chart) ShowChartInfo(c *gin.Context) {
	name := c.Query("chart")
	if name == "" {
		app.RespErr(c, fmt.Errorf("chart name can not be empty"))
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
		app.RespErr(c, fmt.Errorf("bad info %s, chart info only support readme/values/chart", info))
		return
	}

	cp, err := client.ChartPathOptions.LocateChart(name, global.HelmClientSettings)
	if err != nil {
		app.RespErr(c, err)
		return
	}

	chrt, err := loader.Load(cp)
	if err != nil {
		app.RespErr(c, err)
		return
	}

	if client.OutputFormat == action.ShowChart {
		app.RespOK(c, chrt.Metadata)
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
		app.RespOK(c, values)
		return
	}
	if client.OutputFormat == action.ShowReadme {
		app.RespOK(c, string(findReadme(chrt.Files).Data))
		return
	}
}
