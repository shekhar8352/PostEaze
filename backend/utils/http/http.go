package httpclient

import (
	"context"
	"io"
	"net/http"
	"sync"
)

type client struct {
	RequestConfigStore map[string]*requestMapping
	mu                 sync.RWMutex
}

type requestMapping struct {
	httpClient    http.Client
	RequestConfig RequestConfig
}

var httpClient client

// InitHttp initializes and stores request configs globally
func InitHttp(configs ...*RequestConfig) {
	requestConfigStore := make(map[string]*requestMapping)
	for _, cfg := range configs {
		if cfg != nil && cfg.name != "" {
			clientRequestMapping :=
				requestMapping{
					httpClient:    http.Client{Timeout: cfg.timeout},
					RequestConfig: *cfg,
				}
			requestConfigStore[cfg.name] = &clientRequestMapping
		}
	}
	httpClient = client{
		RequestConfigStore: requestConfigStore,
	}
}

func GetClient() *client {
	return &httpClient
}

type Request struct {
	name         string
	ctx          context.Context
	method       string
	url          string
	queryParams  map[string]string
	headerParams map[string]string
	body         io.Reader
}

func (c *client) getRequestDetails(requestName string) *Request {
	c.mu.RLock()
	config := httpClient.RequestConfigStore[requestName].RequestConfig
	c.mu.RUnlock()
	return &Request{
		name:         config.name,
		method:       config.method,
		url:          config.url,
		headerParams: config.headers,
	}
}

func (c *client) createRequest(ctx context.Context, request *Request) (*http.Request, error) {
	return getRequest(ctx, request.method, request.url, request.queryParams, request.headerParams, request.body)
}

func getRequest(ctx context.Context, method string, url string, queryParams map[string]string,
	headerParams map[string]string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if ctx != nil {
		request = request.WithContext(ctx)
	}

	if queryParams != nil {
		q := request.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}

	for k, v := range headerParams {
		request.Header.Add(k, v)
	}

	return request, err
}

// SetContext is used to set the context for the request
func (req *Request) SetContext(ctx context.Context) *Request {
	req.ctx = ctx
	return req
}

// SetMethod is used to set the method for the request
func (req *Request) SetMethod(method string) *Request {
	req.method = method
	return req
}

// SetURL is used to set the url for the request
// if not done, then the url already configured will be used
func (req *Request) SetURL(url string) *Request {
	req.url = url
	return req
}

// SetQueryParam is used to set a query param key value pair
// These will be passed in query param while executing HTTP request
func (req *Request) SetQueryParam(param, value string) *Request {
	if req.queryParams == nil {
		req.queryParams = make(map[string]string)
	}
	req.queryParams[param] = value
	return req
}

// SetQueryParams is used to set multiple query params - map of key-value pair
// These will be passed in query param while executing HTTP request
func (req *Request) SetQueryParams(queryParams map[string]string) *Request {
	if req.queryParams == nil {
		req.queryParams = make(map[string]string)
	}
	for k, v := range queryParams {
		req.queryParams[k] = v
	}
	return req
}

// SetHeaderParam is used to set a header - key-value pair
// These will be passed in header while executing HTTP request
func (req *Request) SetHeaderParam(param, value string) *Request {
	if req.headerParams == nil {
		req.headerParams = make(map[string]string)
	}
	req.headerParams[param] = value
	return req
}

// SetHeaderParams is used to set multiple headers -  map of key-value pair
// These will be passed in header while executing HTTP request
func (req *Request) SetHeaderParams(headerParams map[string]string) *Request {
	if req.headerParams == nil {
		req.headerParams = make(map[string]string)
	}
	for k, v := range headerParams {
		req.headerParams[k] = v
	}
	return req
}

// SetBody is used to set request body to pass in http request
func (req *Request) SetBody(body io.Reader) *Request {
	req.body = body
	return req
}
