package naming_client

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
	"time"
)

func TestEventDispatcher_AddCallbackFuncs(t *testing.T) {
	service := model.Service{
		Dom:         "public@@Test",
		Clusters:    strings.Join([]string{"default"}, ","),
		CacheMillis: 10000,
		Checksum:    "abcd",
		LastRefTime: uint64(time.Now().Unix()),
	}
	var hosts []model.Instance
	host := model.Instance{
		Valid:       true,
		Enable:      true,
		InstanceId:  "123",
		Port:        8080,
		Ip:          "127.0.0.1",
		Weight:      10,
		ServiceName: "public@@Test",
		ClusterName: strings.Join([]string{"default"}, ","),
	}
	hosts = append(hosts, host)
	service.Hosts = hosts

	ed := NewSubscribeCallback()
	param := vo.SubscribeParam{
		ServiceName: "Test",
		Clusters:    []string{"default"},
		GroupName:   "public",
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Println(utils.ToJsonString(ed.callbackFuncsMap))
		},
	}

	clusters := param.Clusters

	if clusters == nil || len(clusters) == 0 {
		clusters = []string{constant.STRING_EMPTY}
	}

	for index := range clusters {
		ed.AddCallbackFuncs(clusters[index], utils.GetGroupName(param.ServiceName, param.GroupName), param.CallbackFuncId, &param.SubscribeCallback)
	}
	key := utils.GetServiceCacheKey(utils.GetGroupName(param.ServiceName, param.GroupName), strings.Join(param.Clusters, ","))
	for k, v := range ed.callbackFuncsMap.Items() {
		assert.Equal(t, key, k, "key should be equal!")
		funcs := v.([]*func(services []model.SubscribeService, err error))
		assert.Equal(t, len(funcs), 1)
		assert.Equal(t, funcs[0], &param.SubscribeCallback, "callback function must be equal!")

	}
}

func TestEventDispatcher_RemoveCallbackFuncs(t *testing.T) {
	service := model.Service{
		Dom:         "public@@Test",
		Clusters:    strings.Join([]string{"default"}, ","),
		CacheMillis: 10000,
		Checksum:    "abcd",
		LastRefTime: uint64(time.Now().Unix()),
	}
	var hosts []model.Instance
	host := model.Instance{
		Valid:       true,
		Enable:      true,
		InstanceId:  "123",
		Port:        8080,
		Ip:          "127.0.0.1",
		Weight:      10,
		ServiceName: "public@@Test",
		ClusterName: strings.Join([]string{"default"}, ","),
	}
	hosts = append(hosts, host)
	service.Hosts = hosts

	ed := NewSubscribeCallback()
	param := vo.SubscribeParam{
		ServiceName: "Test",
		Clusters:    []string{"default"},
		GroupName:   "public",
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("func1:%s \n", utils.ToJsonString(services))
		},
	}
	clusters := param.Clusters

	if clusters == nil || len(clusters) == 0 {
		clusters = []string{constant.STRING_EMPTY}
	}

	for index := range clusters {
		ed.AddCallbackFuncs(clusters[index], utils.GetGroupName(param.ServiceName, param.GroupName), param.CallbackFuncId, &param.SubscribeCallback)
	}
	assert.Equal(t, len(ed.callbackFuncsMap.Items()), 1, "callback funcs map length should be 1")

	param2 := vo.SubscribeParam{
		ServiceName: "Test",
		Clusters:    []string{"default"},
		GroupName:   "public",
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("func2:%s \n", utils.ToJsonString(services))
		},
	}
	clusters2 := param2.Clusters

	if clusters2 == nil || len(clusters2) == 0 {
		clusters2 = []string{constant.STRING_EMPTY}
	}

	for index2 := range clusters2 {
		ed.AddCallbackFuncs(clusters2[index2], utils.GetGroupName(param2.ServiceName, param2.GroupName), param2.CallbackFuncId, &param2.SubscribeCallback)
	}
	assert.Equal(t, len(ed.callbackFuncsMap.Items()), 1, "callback funcs map length should be 2")

	for k, v := range ed.callbackFuncsMap.Items() {
		log.Printf("key:%s,%d", k, len(v.([]*func(services []model.SubscribeService, err error))))
	}

	for index2 := range clusters2 {
		ed.RemoveCallbackFuncs(clusters2[index2], utils.GetGroupName(param2.ServiceName, param2.GroupName), param2.CallbackFuncId, &param2.SubscribeCallback)
	}

	key := utils.GetServiceCacheKey(utils.GetGroupName(param.ServiceName, param.GroupName), strings.Join(param.Clusters, ","))
	for k, v := range ed.callbackFuncsMap.Items() {
		assert.Equal(t, key, k, "key should be equal!")
		funcs := v.([]*func(services []model.SubscribeService, err error))
		assert.Equal(t, len(funcs), 1)
		assert.Equal(t, funcs[0], &param.SubscribeCallback, "callback function must be equal!")

	}
}

func TestSubscribeCallback_ServiceChanged(t *testing.T) {
	service := model.Service{
		Name:        "public@@Test",
		Clusters:    strings.Join([]string{"default"}, ","),
		CacheMillis: 10000,
		Checksum:    "abcd",
		LastRefTime: uint64(time.Now().Unix()),
	}
	var hosts []model.Instance
	host := model.Instance{
		Valid:       true,
		Enable:      true,
		InstanceId:  "123",
		Port:        8080,
		Ip:          "127.0.0.1",
		Weight:      10,
		ServiceName: "public@@Test",
		ClusterName: strings.Join([]string{"default"}, ","),
	}
	hosts = append(hosts, host)
	service.Hosts = hosts

	ed := NewSubscribeCallback()
	param := vo.SubscribeParam{
		ServiceName: "Test",
		Clusters:    []string{"default"},
		GroupName:   "public",
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			log.Printf("func1:%s \n", utils.ToJsonString(services))
		},
	}
	clusters := param.Clusters

	if clusters == nil || len(clusters) == 0 {
		clusters = []string{constant.STRING_EMPTY}
	}

	for index := range clusters {
		ed.AddCallbackFuncs(clusters[index], utils.GetGroupName(param.ServiceName, param.GroupName), param.CallbackFuncId, &param.SubscribeCallback)
	}

	param2 := vo.SubscribeParam{
		ServiceName: "Test",
		Clusters:    []string{"default"},
		GroupName:   "public",
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			log.Printf("func2:%s \n", utils.ToJsonString(services))

		},
	}

	clusters2 := param2.Clusters

	if clusters2 == nil || len(clusters2) == 0 {
		clusters2 = []string{constant.STRING_EMPTY}
	}

	for index2 := range clusters2 {
		ed.AddCallbackFuncs(clusters2[index2], utils.GetGroupName(param2.ServiceName, param2.GroupName), param2.CallbackFuncId, &param2.SubscribeCallback)
	}

	ed.ServiceChanged(&service)
}
