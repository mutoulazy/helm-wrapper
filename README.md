# A [Helm3](https://github.com/helm/helm) HTTP Wrapper With Go SDK

通过 Go [Gin](https://github.com/gin-gonic/gin) Web 框架，结合 Helm Go SDK 封装的 HTTP Server，
让 Helm 相关的日常命令操作可以通过 Restful API 的方式来实现命令行同样的操作。
> __注：__  基于helm-wrapper原版进行了一些工程化的改造

## Support API

* 如果某些API需要支持多个集群，则可以使用以下参数

| Params | Description |
| :- | :- |
| kube_context | 支持指定kube_context来区分不同集群 |

支持helm install的基本增删回查的命令，详情可以看swagger文档
[swagger文档](http://localhost:8080/swagger/index.html)

### 响应

为了简化，所有请求统一返回 200 状态码，通过返回 Body 中的 Code 值来判断响应是否正常：

``` go
type ResponseBody struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Details []string    `json:"details,omitempty"`
}
```

并且构建了统一错误码和日志文件用于快速定位问题

### 监控指标
采用[gin-metrics](github.com/penglongli/gin-metrics)组件为应用提供业务相关监控指标
http://localhost:8080/metrics

| Metric                  | Type      | Description                                         |
| ----------------------- | --------- | --------------------------------------------------- |
| gin_request_total       | Counter   | 服务接收到的请求总数                |
| gin_request_uv          | Counter   | 服务接收到的 IP 总数                     |
| gin_uri_request_total   | Counter   | 每个 URI 接收到的服务请求数 |
| gin_request_body_total  | Counter   | 服务接收到的请求量，单位: 字节   |
| gin_response_body_total | Counter   | 服务返回的请求量，单位: 字节      |
| gin_request_duration    | Histogram | 服务处理请求使用的时间         |
| gin_slow_request_total  | Counter   | 服务接收到的慢请求计数     |


## Build & Run 

### Build

源码提供了简单的 `Makefile` 文件，如果要构建二进制，只需要通过以下方式构建即可。

```
make build          // 构建当前主机架构的二进制版本
make build-linux    // 构建 Linux 版本的二进制
make build-docker   // 构建 Docker 镜像
```

直接构建会生成名为 `helm-wrapper` 的二进制程序，你可以通过如下方式获取帮助：

```
$ helm-wrapper -h
Usage of helm-wrapper:
      --addr string                      server listen addr (default "0.0.0.0")
      --alsologtostderr                  log to standard error as well as files
      --config string                    helm wrapper config (default "config.yaml")
      --debug                            enable verbose output
      --kube-context string              name of the kubeconfig context to use
      --kubeconfig string                path to the kubeconfig file
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files (default true)
  -n, --namespace string                 namespace scope for this request
      --port string                      server listen port (default "8080")
      --registry-config string           path to the registry config file (default "/root/.config/helm/registry.json")
      --repository-cache string          path to the file containing cached repository indexes (default "/root/.cache/helm/repository")
      --repository-config string         path to the file containing repository names and URLs (default "/root/.config/helm/repositories.yaml")
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
pflag: help requested
```

关键性的选项说明一下：

+ `--config` helm-wrapper 的配置项，内容如下，主要是指定 Helm Repo 命名和 URL，用于 Repo 初始化。

```
$ cat config-example.yaml
logSavePath: storage/logs
logFileName: app
logFileExt: .log
writeTimeout: 60
readTimeout: 60
defaultContextTimeout: 60
#runMode: release debug test
runMode: debug
uploadPath: /tmp/charts
helmRepos:
#  - name: bitnami
#    url: https://charts.bitnami.com/bitnami
```
+ `--kubeconfig` 默认如果你不指定的话，使用默认的路径，一般是 `~/.kube/config`。这个配置是必须的，这指明了你要操作的 Kubernetes 集群地址以及访问方式。`kubeconfig` 文件如何生成，这里不过多介绍，具体可以详见 [Configure Access to Multiple Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)

### Run

运行比较简单，如果你本地已经有默认的 `kubeconfig` 文件，只需要把 helm-wrapper 需要的 repo 配置文件配置好即可，然后执行以下命令即可运行，示例如下：

```
$ ./helm-wrapper --config </path/to/config.yaml> --kubeconfig </path/to/kubeconfig>
```

> 启动时会先初始化 repo，因此根据 repo 本身的大小或者网络因素，会耗费些时间

#### 运行在 Kubernetes 集群中

替换 `deployment/deployment.yaml` 中 image 字段为你正确的 helm-wrapper 镜像地址即可，然后执行命令部署：

```
kubectl create -f ./deployment
```

> __注：__ 以上操作会创建 RBAC 相关，因此不需要在构建镜像的时候额外添加 kubeconfig 文件，默认会拥有相关的权限
