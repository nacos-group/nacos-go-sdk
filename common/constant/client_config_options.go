/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package constant

import (
	"os"
	"time"

	"github.com/nacos-group/nacos-sdk-go/common/file"
)

func NewClientConfig(opts ...ClientOption) *ClientConfig {
	clientConfig := &ClientConfig{
		TimeoutMs:            10 * 1000,
		BeatInterval:         5 * 1000,
		OpenKMS:              false,
		CacheDir:             file.GetCurrentPath() + string(os.PathSeparator) + "cache",
		UpdateThreadNum:      20,
		NotLoadCacheAtStart:  false,
		UpdateCacheWhenEmpty: false,
		LogDir:               file.GetCurrentPath() + string(os.PathSeparator) + "log",
		RotateTime:           "24h",
		MaxAge:               3,
		LogLevel:             "info",
	}

	for _, opt := range opts {
		opt(clientConfig)
	}

	return clientConfig
}

// ClientOption ...
type ClientOption func(*ClientConfig)

// WithTimeoutMs ...
func WithTimeoutMs(timeoutMs uint64) ClientOption {
	return func(config *ClientConfig) {
		config.TimeoutMs = timeoutMs
	}
}

// WithBeatInterval ...
func WithBeatInterval(beatInterval int64) ClientOption {
	return func(config *ClientConfig) {
		config.BeatInterval = beatInterval
	}
}

// WithNamespaceId ...
func WithNamespaceId(namespaceId string) ClientOption {
	return func(config *ClientConfig) {
		config.NamespaceId = namespaceId
	}
}

// WithEndpoint ...
func WithEndpoint(endpoint string) ClientOption {
	return func(config *ClientConfig) {
		config.Endpoint = endpoint
	}
}

// WithRegionId ...
func WithRegionId(regionId string) ClientOption {
	return func(config *ClientConfig) {
		config.RegionId = regionId
	}
}

// WithAccessKey ...
func WithAccessKey(accessKey string) ClientOption {
	return func(config *ClientConfig) {
		config.AccessKey = accessKey
	}
}

// WithSecretKey ...
func WithSecretKey(secretKey string) ClientOption {
	return func(config *ClientConfig) {
		config.SecretKey = secretKey
	}
}

// WithOpenKMS ...
func WithOpenKMS(openKMS bool) ClientOption {
	return func(config *ClientConfig) {
		config.OpenKMS = openKMS
	}
}

// WithCacheDir ...
func WithCacheDir(cacheDir string) ClientOption {
	return func(config *ClientConfig) {
		config.CacheDir = cacheDir
	}
}

// WithUpdateThreadNum ...
func WithUpdateThreadNum(updateThreadNum int) ClientOption {
	return func(config *ClientConfig) {
		config.UpdateThreadNum = updateThreadNum
	}
}

// WithNotLoadCacheAtStart ...
func WithNotLoadCacheAtStart(notLoadCacheAtStart bool) ClientOption {
	return func(config *ClientConfig) {
		config.NotLoadCacheAtStart = notLoadCacheAtStart
	}
}

// WithUpdateCacheWhenEmpty ...
func WithUpdateCacheWhenEmpty(updateCacheWhenEmpty bool) ClientOption {
	return func(config *ClientConfig) {
		config.UpdateCacheWhenEmpty = updateCacheWhenEmpty
	}
}

// WithUsername ...
func WithUsername(username string) ClientOption {
	return func(config *ClientConfig) {
		config.Username = username
	}
}

// WithPassword ...
func WithPassword(password string) ClientOption {
	return func(config *ClientConfig) {
		config.Password = password
	}
}

// WithLogDir ...
func WithLogDir(logDir string) ClientOption {
	return func(config *ClientConfig) {
		config.LogDir = logDir
	}
}

// WithRotateTime ...
func WithRotateTime(rotateTime string) ClientOption {
	return func(config *ClientConfig) {
		config.RotateTime = rotateTime
	}
}

// WithMaxAge ...
func WithMaxAge(maxAge int64) ClientOption {
	return func(config *ClientConfig) {
		config.MaxAge = maxAge
	}
}

// WithLogLevel ...
func WithLogLevel(logLevel string) ClientOption {
	return func(config *ClientConfig) {
		config.LogLevel = logLevel
	}
}

// WithLogSampling ...
func WithLogSampling(tick time.Duration, initial int, thereafter int) ClientOption {
	return func(config *ClientConfig) {
		config.LogSampling = &ClientLogSamplingConfig{initial, thereafter, tick}
	}
}
