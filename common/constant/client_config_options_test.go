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
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/common/file"

	"github.com/stretchr/testify/assert"
)

func TestNewClientConfig(t *testing.T) {
	config := NewClientConfig()

	assert.Equal(t, config.TimeoutMs, uint64(10000))
	assert.Equal(t, config.Endpoint, "")
	assert.Equal(t, config.LogLevel, "info")
	assert.Equal(t, config.BeatInterval, int64(5000))
	assert.Equal(t, config.UpdateThreadNum, 20)
	assert.Equal(t, config.RotateTime, "24h")

	assert.Equal(t, config.LogDir, file.GetCurrentPath()+string(os.PathSeparator)+"log")
	assert.Equal(t, config.CacheDir, file.GetCurrentPath()+string(os.PathSeparator)+"cache")

	assert.Equal(t, config.MaxAge, int64(3))
	assert.Equal(t, config.NotLoadCacheAtStart, false)
	assert.Equal(t, config.UpdateCacheWhenEmpty, false)

	assert.Equal(t, config.Username, "")
	assert.Equal(t, config.Password, "")
	assert.Equal(t, config.OpenKMS, false)
	assert.Equal(t, config.NamespaceId, "")
	assert.Equal(t, config.Username, "")
	assert.Equal(t, config.RegionId, "")
	assert.Equal(t, config.AccessKey, "")
	assert.Equal(t, config.SecretKey, "")
	assert.Nil(t, config.LogSampling)
}

func TestNewClientConfigWithOptions(t *testing.T) {
	config := NewClientConfig(
		WithTimeoutMs(uint64(20000)),
		WithEndpoint("http://console.nacos.io:80"),
		WithLogLevel("error"),
		WithBeatInterval(int64(2000)),
		WithUpdateThreadNum(30),
		WithRotateTime("16h"),

		WithLogDir("/tmp/nacos/log"),
		WithCacheDir("/tmp/nacos/cache"),

		WithMaxAge(int64(6)),
		WithNotLoadCacheAtStart(true),
		WithUpdateCacheWhenEmpty(true),

		WithUsername("nacos"),
		WithPassword("nacos"),
		WithOpenKMS(true),
		WithRegionId("shanghai"),
		WithNamespaceId("namespace_1"),
		WithAccessKey("accessKey_1"),
		WithSecretKey("secretKey_1"),

		WithLogSampling(time.Second*10, 5, 10),
	)

	assert.Equal(t, config.TimeoutMs, uint64(20000))
	assert.Equal(t, config.Endpoint, "http://console.nacos.io:80")
	assert.Equal(t, config.LogLevel, "error")
	assert.Equal(t, config.BeatInterval, int64(2000))
	assert.Equal(t, config.UpdateThreadNum, 30)
	assert.Equal(t, config.RotateTime, "16h")

	assert.Equal(t, config.LogDir, "/tmp/nacos/log")
	assert.Equal(t, config.CacheDir, "/tmp/nacos/cache")

	assert.Equal(t, config.MaxAge, int64(6))
	assert.Equal(t, config.NotLoadCacheAtStart, true)
	assert.Equal(t, config.UpdateCacheWhenEmpty, true)

	assert.Equal(t, config.Username, "nacos")
	assert.Equal(t, config.Password, "nacos")
	assert.Equal(t, config.OpenKMS, true)
	assert.Equal(t, config.RegionId, "shanghai")
	assert.Equal(t, config.NamespaceId, "namespace_1")
	assert.Equal(t, config.AccessKey, "accessKey_1")
	assert.Equal(t, config.SecretKey, "secretKey_1")

	assert.Equal(t, config.LogSampling.Tick, time.Second*10)
	assert.Equal(t, config.LogSampling.Initial, 5)
	assert.Equal(t, config.LogSampling.Thereafter, 10)
}
