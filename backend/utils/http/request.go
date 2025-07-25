package httpclient

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

type RequestConfig struct {
	name                  string
	method                string
	url                   string
	timeout               time.Duration
	connectTimeout        time.Duration
	keepAlive             time.Duration
	maxIdleConnections    int
	idleConnectionTimeout time.Duration
	tlsHandshakeTimeout   time.Duration
	proxyURL              string
	retryCount            int
	headers               map[string]string
}

func NewRequestConfig(name string, configMap map[string]interface{}) *RequestConfig {
	rc := &RequestConfig{
		name: name,
	}

	if configMap != nil {
		rc.method, _ = getConfigOptionString(configMap, "method")
		rc.url, _ = getConfigOptionString(configMap, "url")

		if timeoutMs, err := getConfigOptionInt(configMap, "timeoutinmillis"); err == nil {
			rc.timeout = time.Duration(timeoutMs) * time.Millisecond
		}

		if retryCount, err := getConfigOptionInt(configMap, "retrycount"); err == nil {
			rc.retryCount = retryCount
		}

		if headers, err := getConfigOptionMap(configMap, "headers"); err == nil {
			rc.headers = cast.ToStringMapString(headers)
		}
	}
	return rc
}

func getConfigOptionInt(options map[string]interface{}, key string) (int, error) {
	var val interface{}
	var ok bool
	var s int
	if val, ok = options[key]; ok {
		return cast.ToIntE(val)
	} else {
		return s, fmt.Errorf("missing %s", key)
	}
}

func getConfigOptionFloat(options map[string]interface{}, key string) (float64, error) {
	var val interface{}
	var ok bool
	var s float64
	if val, ok = options[key]; ok {
		return cast.ToFloat64E(val)
	} else {
		return s, fmt.Errorf("missing %s", key)
	}
}

func getConfigOptionMap(options map[string]interface{}, key string) (map[string]interface{}, error) {
	var val interface{}
	var ok bool
	var s map[string]interface{}
	if val, ok = options[key]; ok {
		return cast.ToStringMapE(val)
	} else {
		return s, fmt.Errorf("missing %s", key)
	}
}

func getConfigOptionString(options map[string]interface{}, key string) (string, error) {
	var val interface{}
	var ok bool
	var s string
	if val, ok = options[key]; ok {
		return cast.ToStringE(val)
	} else {
		return s, fmt.Errorf("missing %s", key)
	}
}

func getConfigOptionBool(options map[string]interface{}, key string) (bool, error) {
	var val interface{}
	var b, ok bool
	if val, ok = options[key]; ok {
		return cast.ToBoolE(val)
	} else {
		return b, fmt.Errorf("missing %s", key)
	}
}
