# Nacos-sdk-go [中文](./README_CN.md) #

[![Build Status](https://travis-ci.org/nacos-group/nacos-sdk-go.svg?branch=master)](https://travis-ci.org/nacos-group/nacos-sdk-go) [![Go Report Card](https://goreportcard.com/badge/github.com/nacos-group/nacos-sdk-go)](https://goreportcard.com/report/github.com/nacos-group/nacos-sdk-go) ![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

---

## Nacos-sdk-go

Nacos-sdk-go for Go client allows you to access Nacos service,it supports service discovery and dynamic configuration.

## Requirements
Supported Go version over 1.12

Supported Nacos version over 1.x

## Installation
Use `go get` to install SDK：
```sh
$ go get -u github.com/nacos-group/nacos-sdk-go
```
## Quick Examples
* ClientConfig

```go
constant.ClientConfig{
	TimeoutMs            uint64 // timeout for requesting Nacos server, default value is 10000ms
	NamespaceId          string // the namespaceId of Nacos
	Endpoint             string // the endpoint for ACM. https://help.aliyun.com/document_detail/130146.html
	RegionId             string // the regionId for ACM & KMS
	AccessKey            string // the AccessKey for ACM & KMS
	SecretKey            string // the SecretKey for ACM & KMS
	OpenKMS              bool   // it's to open KMS, default is false. https://help.aliyun.com/product/28933.html
	                            // , to enable encrypt/decrypt, DataId should be start with "cipher-"
	CacheDir             string // the directory for persist nacos service info,default value is current path
	UpdateThreadNum      int    // the number of goroutine for update nacos service info,default value is 20
	NotLoadCacheAtStart  bool   // not to load persistent nacos service info in CacheDir at start time
	UpdateCacheWhenEmpty bool   // update cache when get empty service instance from server
	Username             string // the username for nacos auth
	Password             string // the password for nacos auth
	LogDir               string // the directory for log, default is current path
	RotateTime           string // the rotate time for log, eg: 30m, 1h, 24h, default is 24h
	MaxAge               int64  // the max age of a log file, default value is 3
	LogLevel             string // the level of log, it's must be debug,info,warn,error, default value is info
}
```


* ServerConfig

```go
constant.ServerConfig{
	ContextPath string // the nacos server context path
	IpAddr      string // the nacos server address
	Port        uint64 // the nacos server port
	Scheme      string // the nacos server scheme
}
```

<b>Note：We can config multiple ServerConfig,the client will rotate request the servers</b>

### Create client

```go
//create clientConfig
clientConfig := constant.ClientConfig{
	NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         
	TimeoutMs:           5000,
	NotLoadCacheAtStart: true,
	LogDir:              "/tmp/nacos/log",
	CacheDir:            "/tmp/nacos/cache",
	RotateTime:          "1h",
	MaxAge:              3,
	LogLevel:            "debug",
}
//Another way of create clientConfig
clientConfig := *constant.NewClientConfig(
    constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"), //When namespace is public, fill in the blank string here.
    constant.WithTimeoutMs(5000),
    constant.WithNotLoadCacheAtStart(true),
    constant.WithLogDir("/tmp/nacos/log"),
    constant.WithCacheDir("/tmp/nacos/cache"),
    constant.WithRotateTime("1h"),
    constant.WithMaxAge(3),
    constant.WithLogLevel("debug"),
)

// At least one ServerConfig
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
//Another way of create serverConfigs
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

// Create naming client for service discovery
_, _ := clients.CreateNamingClient(map[string]interface{}{
	"serverConfigs": serverConfigs,
	"clientConfig":  clientConfig,
})

// Create config client for dynamic configuration
_, _ := clients.CreateConfigClient(map[string]interface{}{
	"serverConfigs": serverConfigs,
	"clientConfig":  clientConfig,
})

// Another way of create naming client for service discovery (recommend)
namingClient, err := clients.NewNamingClient(
    vo.NacosClientParam{
        ClientConfig:  &clientConfig,
        ServerConfigs: serverConfigs,
    },
)

// Another way of create config client for dynamic configuration (recommend)
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


### Service Discovery

* Register instance：RegisterInstance

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
    ClusterName: "cluster-a", // default value is DEFAULT
    GroupName:   "group-a",  // default value is DEFAULT_GROUP
})

```

* Deregister instance：DeregisterInstance

```go

success, err := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
    Ip:          "10.0.0.11",
    Port:        8848,
    ServiceName: "demo.go",
    Ephemeral:   true,
    Cluster:     "cluster-a", // default value is DEFAULT
    GroupName:   "group-a",   // default value is DEFAULT_GROUP
})

```

* Get service：GetService

```go

services, err := namingClient.GetService(vo.GetServiceParam{
    ServiceName: "demo.go",
    Clusters:    []string{"cluster-a"}, // default value is DEFAULT
    GroupName:   "group-a",             // default value is DEFAULT_GROUP
})

```

* Get all instances：SelectAllInstances

```go
// SelectAllInstance return all instances,include healthy=false,enable=false,weight<=0
instances, err := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // default value is DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // default value is DEFAULT
})

```

* Get instances ：SelectInstances

```go
// SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
instances, err := namingClient.SelectInstances(vo.SelectInstancesParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // default value is DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // default value is DEFAULT
    HealthyOnly: true,
})

```

* Get one healthy instance（WRR）：SelectOneHealthyInstance

```go
// SelectOneHealthyInstance return one instance by WRR strategy for load balance
// And the instance should be health=true,enable=true and weight>0
instance, err := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // default value is DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // default value is DEFAULT
})

```

* Listen service change event：Subscribe

```go

// Subscribe key = serviceName+groupName+cluster
// Note: We call add multiple SubscribeCallback with the same key.
err := namingClient.Subscribe(vo.SubscribeParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // default value is DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // default value is DEFAULT
    SubscribeCallback: func(services []model.Instance, err error) {
        log.Printf("\n\n callback return services:%s \n\n", utils.ToJsonString(services))
    },
})

```

* Cancel listen of service change event：Unsubscribe

```go

err := namingClient.Unsubscribe(vo.SubscribeParam{
    ServiceName: "demo.go",
    GroupName:   "group-a",             // default value is DEFAULT_GROUP
    Clusters:    []string{"cluster-a"}, // default value is DEFAULT
    SubscribeCallback: func(services []model.Instance, err error) {
        log.Printf("\n\n callback return services:%s \n\n", utils.ToJsonString(services))
    },
})

```

* Get all services name:GetAllServicesInfo
```go

serviceInfos, err := namingClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
    NameSpace: "0e83cc81-9d8c-4bb8-a28a-ff703187543f",
    PageNo:   1,
    PageSize: 10,
    }),

```

### Dynamic configuration

* publish config：PublishConfig

```go

success, err := configClient.PublishConfig(vo.ConfigParam{
    DataId:  "dataId",
    Group:   "group",
    Content: "hello world!222222"})

```

* delete config：DeleteConfig

```go

success, err = configClient.DeleteConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group"})

```

* get config info：GetConfig

```go

content, err := configClient.GetConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group"})

```

* Listen config change event：ListenConfig

```go

err := configClient.ListenConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group",
    OnChange: func(namespace, group, dataId, data string) {
        fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	},
})

```
* Cancel the listening of config change event：CancelListenConfig

```go

err := configClient.CancelListenConfig(vo.ConfigParam{
    DataId: "dataId",
    Group:  "group",
})

```

* Search config: SearchConfig
```go
configPage, err := configClient.SearchConfig(vo.SearchConfigParam{
    Search:   "blur",
    DataId:   "",
    Group:    "",
    PageNo:   1,
    PageSize: 10,
})
```
## Example
We can run example to learn how to use nacos go client.
* [Config Example](./example/config)
* [Naming Example](./example/service)

## Documentation
You can view the open-api documentation from the [Nacos open-api wepsite](https://nacos.io/en-us/docs/open-api.html).

You can view the full documentation from the [Nacos website](https://nacos.io/en-us/docs/what-is-nacos.html).

## Contributing
Contributors are welcomed to join Nacos-sdk-go project. Please check [CONTRIBUTING.md](./CONTRIBUTING.md) about how to contribute to this project.

## Contact
* Join us from DingDing Group(23191211).
* [Gitter](https://gitter.im/alibaba/nacos): Nacos's IM tool for community messaging, collaboration and discovery.
* [Twitter](https://twitter.com/nacos2): Follow along for latest nacos news on Twitter.
* [Weibo](https://weibo.com/u/6574374908): Follow along for latest nacos news on Weibo (Twitter of China version).
* [Nacos SegmentFault](https://segmentfault.com/t/nacos): Get the latest notice and prompt help from SegmentFault.
* Email Group:
     * users-nacos@googlegroups.com: Nacos usage general discussion.
     * dev-nacos@googlegroups.com: Nacos developer discussion (APIs, feature design, etc).
     * commits-nacos@googlegroups.com: Commits notice, very high frequency.

