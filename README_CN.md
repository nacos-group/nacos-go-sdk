# Nacos-sdk-go [English](./README.md) #

[![Build Status](https://travis-ci.org/nacos-group/nacos-sdk-go.svg?branch=master)](https://travis-ci.org/nacos-group/nacos-sdk-go) [![Go Report Card](https://goreportcard.com/badge/github.com/nacos-group/nacos-sdk-go)](https://goreportcard.com/report/github.com/nacos-group/nacos-sdk-go) ![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

---

## Nacos-sdk-go

Nacos-sdk-go是Nacos的Go语言客户端，它实现了服务发现和动态配置的功能

## 使用限制
支持Go>v1.12版本

支持Nacos>1.x版本

## 安装
使用`go get`安装SDK：
```sh
$ go get -u github.com/nacos-group/nacos-sdk-go
```
## 快速使用
* ClientConfig

```go
constant.ClientConfig{
	TimeoutMs            uint64 // 请求Nacos服务端的超时时间，默认是10000ms
	NamespaceId          string // ACM的命名空间Id
	Endpoint             string // 当使用ACM时，需要该配置. https://help.aliyun.com/document_detail/130146.html
	RegionId             string // ACM&KMS的regionId，用于配置中心的鉴权
	AccessKey            string // ACM&KMS的AccessKey，用于配置中心的鉴权
	SecretKey            string // ACM&KMS的SecretKey，用于配置中心的鉴权
	OpenKMS              bool   // 是否开启kms，默认不开启，kms可以参考文档 https://help.aliyun.com/product/28933.html
	                            // 同时DataId必须以"cipher-"作为前缀才会启动加解密逻辑
	CacheDir             string // 缓存service信息的目录，默认是当前运行目录
	UpdateThreadNum      int    // 监听service变化的并发数，默认20
	NotLoadCacheAtStart  bool   // 在启动的时候不读取缓存在CacheDir的service信息
	UpdateCacheWhenEmpty bool   // 当service返回的实例列表为空时，不更新缓存，用于推空保护
	Username             string // Nacos服务端的API鉴权Username
	Password             string // Nacos服务端的API鉴权Password
	LogDir               string // 日志存储路径
	RotateTime           string // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
	MaxAge               int64  // 日志最大文件数，默认3
	LogLevel             string // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
}
```

* ServerConfig

```go
constant.ServerConfig{
	ContextPath string // Nacos的ContextPath
	IpAddr      string // Nacos的服务地址
	Port        uint64 // Nacos的服务端口
	Scheme      string // Nacos的服务地址前缀
}
```

<b>Note：我们可以配置多个ServerConfig，客户端会对这些服务端做轮询请求</b>

### Create client

```go
// 创建clientConfig
clientConfig := constant.ClientConfig{
	NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
	TimeoutMs:           5000,
	NotLoadCacheAtStart: true,
	LogDir:              "/tmp/nacos/log",
	CacheDir:            "/tmp/nacos/cache",
	RotateTime:          "1h",
	MaxAge:              3,
	LogLevel:            "debug",
}

// 创建clientConfig的另一种方式
clientConfig := *constant.NewClientConfig(
    constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"), //当namespace是public时，此处填空字符串。
    constant.WithTimeoutMs(5000),
    constant.WithNotLoadCacheAtStart(true),
    constant.WithLogDir("/tmp/nacos/log"),
    constant.WithCacheDir("/tmp/nacos/cache"),
    constant.WithRotateTime("1h"),
    constant.WithMaxAge(3),
    constant.WithLogLevel("debug"),
)

// 至少一个ServerConfig
serverConfigs := []constant.ServerConfig{
    {
        IpAddr:      "console1.nacos.io",
        ContextPath: "/nacos",
        Port:        80,
        Scheme:      "http",
    },
    {
    	IpAddr:      "console2.nacos.io",
    	ContextPath: "/nacos",
    	Port:        80,
        Scheme:      "http",
    },
}

// 创建serverConfig的另一种方式
serverConfigs := []constant.ServerConfig{
    *constant.NewServerConfig(
        "console1.nacos.io",
        80,
        constant.WithScheme("http"),
        constant.WithContextPath("/nacos")
    ),
    *constant.NewServerConfig(
        "console2.nacos.io",
        80,
        constant.WithScheme("http"),
        constant.WithContextPath("/nacos")
    ),
}

// 创建服务发现客户端
_, _ := clients.CreateNamingClient(map[string]interface{}{
	"serverConfigs": serverConfigs,
	"clientConfig":  clientConfig,
})

// 创建动态配置客户端
_, _ := clients.CreateConfigClient(map[string]interface{}{
	"serverConfigs": serverConfigs,
	"clientConfig":  clientConfig,
})

// 创建服务发现客户端的另一种方式 (推荐)
namingClient, err := clients.NewNamingClient(
    vo.NacosClientParam{
        ClientConfig:  &clientConfig,
        ServerConfigs: serverConfigs,
    },
)

// 创建动态配置客户端的另一种方式 (推荐)
configClient, err := clients.NewConfigClient(
    vo.NacosClientParam{
        ClientConfig:  &clientConfig,
        ServerConfigs: serverConfigs,
    },
)
```

### Create client for ACM
https://help.aliyun.com/document_detail/130146.html

```go
cc := constant.ClientConfig{
  Endpoint:    "acm.aliyun.com:8080",
  NamespaceId: "e525eafa-f7d7-4029-83d9-008937f9d468",
  RegionId:    "cn-shanghai",
  AccessKey:   "LTAI4G8KxxxxxxxxxxxxxbwZLBr",
  SecretKey:   "n5jTL9YxxxxxxxxxxxxaxmPLZV9",
  OpenKMS:     true,
  TimeoutMs:   5000,
  LogLevel:    "debug",
}

// a more graceful way to create config client
client, err := clients.NewConfigClient(
  vo.NacosClientParam{
    ClientConfig: &cc,
  },
)
```


### 服务发现

* 注册实例：RegisterInstance

```go

success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
    Ip:          "10.0.0.11",
    Port:        8848,
    ServiceName: "demo.go",
    Weight:      10,
    Enable:      true,
    Healthy:     true,
    Ephemeral:   true,
    Metadata:    map[string]string{"idc":"shanghai"},
    ClusterName: "cluster-a", // 默认值DEFAULT
    GroupName:   "group-a",   // 默认值DEFAULT_GROUP
})

```

* 注销实例：DeregisterInstance

```go

success, err := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
    Ip:          "10.0.0.11",
    Port:        8848,
    ServiceName: "demo.go",
    Ephemeral:   true,
    Cluster:     "cluster-a", // 默认值DEFAULT
    GroupName:   "group-a",   // 默认值DEFAULT_GROUP
})

```

* 获取服务信息：GetService

```go

services, err := namingClient.GetService(vo.GetServiceParam{
    ServiceName: "demo.go",
    Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
    GroupName:   "group-a",             // 默认值DEFAULT_GROUP
})

```

* 获取所有的实例列表：SelectAllInstances

```go
// SelectAllInstance可以返回全部实例列表,包括healthy=false,enable=false,weight<=0
instances, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // 默认值DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
})

```

* 获取实例列表 ：SelectInstances

```go
// SelectInstances 只返回满足这些条件的实例列表：healthy=${HealthyOnly},enable=true 和weight>0
instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // 默认值DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
    HealthyOnly: true,
})

```

* 获取一个健康的实例（加权随机轮询）：SelectOneHealthyInstance

```go
// SelectOneHealthyInstance将会按加权随机轮询的负载均衡策略返回一个健康的实例
// 实例必须满足的条件：health=true,enable=true and weight>0
instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // 默认值DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
})

```

* 监听服务变化：Subscribe

```go

// Subscribe key=serviceName+groupName+cluster
// 注意:我们可以在相同的key添加多个SubscribeCallback.
err := namingClient.Subscribe(vo.SubscribeParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // 默认值DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
    SubscribeCallback: func(services []model.Instance, err error) {
        log.Printf("\n\n callback return services:%s \n\n", utils.ToJsonString(services))
    },
})

```

* 取消服务监听：Unsubscribe

```go

err := namingClient.Unsubscribe(vo.SubscribeParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // 默认值DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
    SubscribeCallback: func(services []model.Instance, err error) {
        log.Printf("\n\n callback return services:%s \n\n", utils.ToJsonString(services))
    },
})

```

* 获取服务名列表:GetAllServicesInfo
```go

serviceInfos, err := namingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
    NameSpace: "0e83cc81-9d8c-4bb8-a28a-ff703187543f",
    PageNo:   1,
    PageSize: 10,
	}),

```

### 动态配置

* 发布配置：PublishConfig

```go

success, err := configClient.PublishConfig(vo.ConfigParam{
    DataId:  "dataId",
    Group:   "group",
    Content: "hello world!222222"})

```

* 删除配置：DeleteConfig

```go

success, err = configClient.DeleteConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group"})

```

* 获取配置：GetConfig

```go

content, err := configClient.GetConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group"})

```

* 监听配置变化：ListenConfig

```go

err := configClient.ListenConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group",
    OnChange: func(namespace, group, dataId, data string) {
        fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	},
})

```
* 取消配置监听：CancelListenConfig

```go

err := configClient.CancelListenConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group",
})

```

* 搜索配置: SearchConfig
```go
configPage,err := configClient.SearchConfig(vo.SearchConfigParam{
    Search:   "blur",
    DataId:   "",
    Group:    "",
    PageNo:   1,
    PageSize: 10,
})
```
## 例子
我们能从示例中学习如何使用Nacos go客户端
* [动态配置示例](./example/config)
* [服务发现示例](./example/service)

## 文档
Nacos open-api相关信息可以查看文档 [Nacos open-api wepsite](https://nacos.io/en-us/docs/open-api.html).

Nacos产品了解可以查看 [Nacos website](https://nacos.io/en-us/docs/what-is-nacos.html).

## 贡献代码
我们非常欢迎大家为Nacos-sdk-go贡献代码. 贡献前请查看[CONTRIBUTING.md](./CONTRIBUTING.md)

## 联系我们
* 加入Nacos-sdk-go钉钉群(23191211).
* [Gitter](https://gitter.im/alibaba/nacos): Nacos即时聊天工具.
* [Twitter](https://twitter.com/nacos2): 在Twitter上关注Nacos的最新动态.
* [Weibo](https://weibo.com/u/6574374908): 在微博上关注Nacos的最新动态.
* [Nacos SegmentFault](https://segmentfault.com/t/nacos): SegmentFault可以获得最新的推送和帮助.
* Email Group:
     * users-nacos@googlegroups.com: Nacos用户讨论组.
     * dev-nacos@googlegroups.com: Nacos开发者讨论组 (APIs, feature design, etc).
     * commits-nacos@googlegroups.com: Nacos commit提醒.


