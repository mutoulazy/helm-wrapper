package v1

import (
	"context"
	"fmt"
	"helm-wrapper/global"
	"helm-wrapper/pkg/app"
	"helm-wrapper/pkg/errcode"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/flock"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

const searchMaxScore = 25

type Repository struct {
}

func NewRepository() Repository {
	return Repository{}
}

type repoChartElement struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	AppVersion  string `json:"app_version"`
	Description string `json:"description"`
}

type repoChartList []repoChartElement

func applyConstraint(version string, versions bool, res []*search.Result) ([]*search.Result, error) {
	if len(version) == 0 {
		return res, nil
	}

	constraint, err := semver.NewConstraint(version)
	if err != nil {
		return res, errors.Wrap(err, "an invalid version/constraint format")
	}

	data := res[:0]
	foundNames := map[string]bool{}
	for _, r := range res {
		if _, found := foundNames[r.Name]; found {
			continue
		}
		v, err := semver.NewVersion(r.Chart.Version)
		if err != nil || constraint.Check(v) {
			data = append(data, r)
			if !versions {
				foundNames[r.Name] = true // If user hasn't requested all versions, only show the latest that matches
			}
		}
	}

	return data, nil
}

func buildSearchIndex(version string) (*search.Index, error) {
	i := search.NewIndex()
	for _, re := range global.MyHelmConfig.HelmRepos {
		n := re.Name
		f := filepath.Join(global.HelmClientSettings.RepositoryCache, helmpath.CacheIndexFile(n))
		ind, err := repo.LoadIndexFile(f)
		if err != nil {
			glog.Warningf("WARNING: Repo %q is corrupt or missing. Try 'helm repo update'.", n)
			continue
		}

		i.AddRepo(n, ind, len(version) > 0)
	}
	return i, nil
}

func (rep Repository) InitRepository(c *repo.Entry) error {
	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(global.HelmClientSettings.RepositoryConfig), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(global.HelmClientSettings.RepositoryConfig, filepath.Ext(global.HelmClientSettings.RepositoryConfig), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		SafeCloser(fileLock, &err)
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(global.HelmClientSettings.RepositoryConfig)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	r, err := repo.NewChartRepository(c, getter.All(global.HelmClientSettings))
	if err != nil {
		return err
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		return err
	}

	f.Update(c)

	if err := f.WriteFile(global.HelmClientSettings.RepositoryConfig, 0644); err != nil {
		return err
	}

	return nil
}

func updateChart(c *repo.Entry) error {
	r, err := repo.NewChartRepository(c, getter.All(global.HelmClientSettings))
	if err != nil {
		return err
	}
	_, err = r.DownloadIndexFile()
	if err != nil {
		return err
	}

	return nil
}

// @Summary			更新chart镜像库
// @Description 	更新chart仓库信息
// @Tags			Repository
// @Success 		200 {object} app.ResponseBody
// @Router 			/api/v1/repositories [put]
func (rep Repository) UpdateRepositories(c *gin.Context) {
	response := app.NewResponse(c)
	type errRepo struct {
		Name string
		Err  string
	}
	errRepoList := []errRepo{}

	var wg sync.WaitGroup
	for _, c := range global.MyHelmConfig.HelmRepos {
		wg.Add(1)
		go func(c *repo.Entry) {
			defer wg.Done()
			err := updateChart(c)
			if err != nil {
				errRepoList = append(errRepoList, errRepo{
					Name: c.Name,
					Err:  err.Error(),
				})
			}
		}(c)
	}
	wg.Wait()

	if len(errRepoList) > 0 {
		response.ToErrorResponse(errcode.ErrorUpdateRepositoriesFail.WithDetails(fmt.Sprintf("error list: %v", errRepoList)))
		return
	}

	response.ToResponse(gin.H{"msg": "success"})
	return
}

// @Summary 		查找chart/列出本地库中的所有chart
// @Description 	在本地库中查找chart，如没有keyword则列出所有chart
// @Tags			Repository
// @Param 			keyword query string false "搜索关键字"
// @Param   		version query string false "chart版本"
// @Param   		versions query bool false "如果true，查询出每个chart的所有版本；false，只列出每个chart最新版"
// @Success 		200 {object} app.ResponseBody
// @Router 			/api/v1/repositories/charts [get]
func (rep Repository) ListRepoCharts(c *gin.Context) {
	response := app.NewResponse(c)
	version := c.Query("version")   // chart version
	versions := c.Query("versions") // if "true", all versions
	keyword := c.Query("keyword")   // search keyword

	// default stable
	if version == "" {
		version = ">0.0.0"
	}

	index, err := buildSearchIndex(version)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorListRepoChartsFail.WithDetails(err.Error()))
		return
	}

	var res []*search.Result
	if keyword == "" {
		res = index.All()
	} else {
		res, err = index.Search(keyword, searchMaxScore, false)
		if err != nil {
			response.ToErrorResponse(errcode.ErrorListRepoChartsFail.WithDetails(err.Error()))
			return
		}
	}

	search.SortScore(res)
	var versionsB bool
	if versions == "true" {
		versionsB = true
	}
	data, err := applyConstraint(version, versionsB, res)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorListRepoChartsFail.WithDetails(err.Error()))
		return
	}
	chartList := make(repoChartList, 0, len(data))
	for _, v := range data {
		chartList = append(chartList, repoChartElement{
			Name:        v.Name,
			Version:     v.Chart.Version,
			AppVersion:  v.Chart.AppVersion,
			Description: v.Chart.Description,
		})
	}

	response.ToResponse(gin.H{"chartList": chartList})
}

func SafeCloser(fileLock *flock.Flock, err *error) {
	if fileErr := fileLock.Unlock(); fileErr != nil && *err == nil {
		*err = fileErr
		glog.Error(fileErr)
	}
}
