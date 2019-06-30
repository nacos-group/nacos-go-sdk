package nacos_server

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/utils"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NacosServer struct {
	sync.RWMutex
	serverList          []constant.ServerConfig
	httpAgent           http_agent.IHttpAgent
	timeoutMs           uint64
	endpoint            string
	lastSrvRefTime      int64
	vipSrvRefInterMills int64
}

func NewNacosServer(serverList []constant.ServerConfig, httpAgent http_agent.IHttpAgent, timeoutMs uint64, endpoint string) (NacosServer, error) {
	if len(serverList) == 0 && endpoint == "" {
		return NacosServer{}, errors.New("both serverlist  and  endpoint are empty")
	}
	ns := NacosServer{
		serverList:          serverList,
		httpAgent:           httpAgent,
		timeoutMs:           timeoutMs,
		endpoint:            endpoint,
		vipSrvRefInterMills: 10000,
	}
	ns.initRefreshSrvIfNeed()
	return ns, nil
}

func (server *NacosServer) callServer(api string, params map[string]string, method string, curServer string, contextPath string) (result string, err error) {
	if contextPath == "" {
		contextPath = constant.WEB_CONTEXT
	}
	url := "http://" + curServer + contextPath + api
	headers := map[string][]string{}
	headers["Client-Version"] = []string{constant.CLIENT_VERSION}
	headers["User-Agent"] = []string{constant.CLIENT_VERSION}
	headers["Accept-Encoding"] = []string{"gzip,deflate,sdch"}
	headers["Connection"] = []string{"Keep-Alive"}
	headers["RequestId"] = []string{uuid.NewV4().String()}
	headers["Request-Module"] = []string{"Naming"}
	headers["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	var response *http.Response
	response, err = server.httpAgent.Request(method, url, headers, server.timeoutMs, params)
	if err != nil {
		return
	}
	var bytes []byte
	bytes, err = ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return
	}
	result = string(bytes)
	if response.StatusCode == 200 {
		return
	} else {
		err = errors.New(fmt.Sprintf("request return error code %d", response.StatusCode))
		return
	}
}

func (server *NacosServer) ReqApi(api string, params map[string]string, method string) (string, error) {
	srvs := server.serverList
	if srvs == nil || len(srvs) == 0 {
		return "", errors.New("server list is empty")
	}
	//only one server,retry request when error
	if len(srvs) == 1 {
		for i := 0; i < constant.REQUEST_DOMAIN_RETRY_TIME; i++ {
			result, err := server.callServer(api, params, method, getAddress(srvs[0]), srvs[0].ContextPath)
			if err == nil {
				return result, nil
			}
			log.Printf("[ERROR] api<%s>,method:<%s>, params:<%s>, call domain error:<%s> , result:<%s> \n", api, method, utils.ToJsonString(params), err.Error(), result)
		}
		return "", errors.New("retry " + strconv.Itoa(constant.REQUEST_DOMAIN_RETRY_TIME) + " times request failed!")
	} else {
		index := rand.Intn(len(srvs))
		for i := 1; i <= len(srvs); i++ {
			curServer := srvs[index]
			result, err := server.callServer(api, params, method, getAddress(curServer), curServer.ContextPath)
			if err == nil {
				return result, nil
			}
			log.Printf("[ERROR] api<%s>,method:<%s>, params:<%s>, call domain error:<%s> , result:<%s> \n", api, method, utils.ToJsonString(params), err.Error(), result)
			index = (index + i) % len(srvs)
		}
		return "", errors.New("retry " + strconv.Itoa(constant.REQUEST_DOMAIN_RETRY_TIME) + " times request failed!")
	}
}

func (server *NacosServer) initRefreshSrvIfNeed() {
	if server.endpoint == "" {
		return
	}
	server.refreshServerSrvIfNeed()
	go func() {
		time.Sleep(time.Duration(1) * time.Second)
		server.refreshServerSrvIfNeed()
	}()

}

func (server *NacosServer) refreshServerSrvIfNeed() {
	if utils.CurrentMillis()-server.lastSrvRefTime < server.vipSrvRefInterMills && len(server.serverList) > 0 {
		return
	}

	var list []string
	urlString := "http://" + server.endpoint + "/nacos/serverlist"
	result := server.httpAgent.RequestOnlyResult(http.MethodGet, urlString, nil, server.timeoutMs, nil)
	list = strings.Split(result, "\n")
	log.Printf("[info] http nacos server list: <%s> \n", result)

	var servers []constant.ServerConfig
	for _, line := range list {
		if line != "" {
			splitLine := strings.Split(strings.TrimSpace(line), ":")
			port := 80
			var err error
			if len(splitLine) == 2 {
				port, err = strconv.Atoi(splitLine[1])
				if err != nil {
					log.Printf("[ERROR] get port from server:<%s>  error: <%s> \n", line, err.Error())
					continue
				}
			}
			servers = append(servers, constant.ServerConfig{IpAddr: splitLine[0], Port: uint64(port), ContextPath: constant.WEB_CONTEXT})
		}
	}

	if len(servers) > 0 {
		server.Lock()
		if !reflect.DeepEqual(server.serverList, servers) {
			log.Printf("[info] server list is updated, old: <%v>,new:<%v> \n", server.serverList, servers)
		}
		server.serverList = servers
		server.lastSrvRefTime = utils.CurrentMillis()
		server.Unlock()
	}

	return
}

func getAddress(cfg constant.ServerConfig) string {
	return cfg.IpAddr + ":" + strconv.Itoa(int(cfg.Port))
}
