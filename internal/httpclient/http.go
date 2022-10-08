package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	httpTreadsCount              = 100
	httpLimiterWeight            = 1
	httpRequestTimeout           = 3 * time.Minute
	httpRequestRetryCount        = 5
	httpRequestRetryWaitInterval = 30 * time.Second
	httpRequestStatusTemplate    = "request attempt: %d, status: %v"
)

type RequestCallback func() (any, error)

type HttpClient struct {
	ctx     context.Context
	limiter *semaphore.Weighted
	client  *http.Client
}

func NewHTTPClient(ctx context.Context) (*HttpClient, error) {
	client := &HttpClient{
		ctx:     ctx,
		limiter: semaphore.NewWeighted(httpTreadsCount),
		client: &http.Client{
			Timeout: httpRequestTimeout,
		},
	}
	return client, nil
}

func (c *HttpClient) retryDoRequest(requestUrl string, requestHandler RequestCallback) (any, error) {
	var resp any
	var err error
	for tryIdx := 0; tryIdx < httpRequestRetryCount; tryIdx++ {
		resp, err = requestHandler()
		if err != nil {
			formatErr := requestErrorMessage(requestUrl, err)
			log.Printf(httpRequestStatusTemplate, tryIdx, formatErr)
			time.Sleep(httpRequestRetryWaitInterval)
		} else {
			break
		}
	}
	return resp, err
}

func (c *HttpClient) Get(requestURL string, headers http.Header) ([]byte, error) {
	resp, err := c.retryDoRequest(requestURL, func() (any, error) {
		return c.getOnce(requestURL, headers)
	})
	if err != nil {
		return nil, fmt.Errorf("request handle error: %v", err)
	}
	if typedResp, ok := resp.(*http.Response); ok {
		return readResponseBody(typedResp.Body)
	}
	return nil, fmt.Errorf("recieved invalid response")
}

func (c *HttpClient) Post(requestURL string, headers http.Header, payload any) ([]byte, error) {
	resp, err := c.retryDoRequest(requestURL, func() (any, error) {
		return c.postOnce(requestURL, headers, payload)
	})
	if err != nil {
		return nil, fmt.Errorf("request handle error: %v", err)
	}
	typedResp, castOk := resp.(*http.Response)
	if !castOk {
		return nil, fmt.Errorf("recieved invalid response")
	}
	return readResponseBody(typedResp.Body)
}

func (c *HttpClient) postOnce(requestUrl string, headers http.Header, payload any) (*http.Response, error) {
	err := c.limiter.Acquire(c.ctx, httpLimiterWeight)
	if err != nil {
		return nil, fmt.Errorf("cannot acquire the semaphore: %v", err)
	}
	defer c.limiter.Release(httpLimiterWeight)
	body, err := preparePostPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot read payload: %v", err)
	}
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, requestUrl, body)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %v", err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get response: %v", err)
	}
	return resp, nil
}

func (c *HttpClient) getOnce(requestUrl string, headers http.Header) (*http.Response, error) {
	err := c.limiter.Acquire(c.ctx, httpLimiterWeight)
	if err != nil {
		return nil, fmt.Errorf("cannot acquire the semaphore: %v", err)
	}
	defer c.limiter.Release(httpLimiterWeight)
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %v", err)
	}
	req.Header = headers
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get response: %v", err)
	}
	return resp, nil
}

func requestErrorMessage(url string, err error) error {
	return fmt.Errorf("cannot get response from '%s', recieved error: %v", url, err)
}

func preparePostPayload(payload any) (io.Reader, error) {
	var body io.Reader
	switch t := payload.(type) {
	case string:
		body = strings.NewReader(url.QueryEscape(payload.(string)))
	case []byte:
		body = bytes.NewReader(payload.([]byte))
	default:
		return nil, fmt.Errorf("unsupported payload type: %t", t)
	}
	return body, nil
}

func readResponseBody(body io.Reader) ([]byte, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %v", err)
	}
	return b, nil
}
