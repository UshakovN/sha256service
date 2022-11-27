package httpclient

import (
  "context"
  "fmt"
  "golang.org/x/sync/semaphore"
  "io"
  "log"
  "net/http"
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

func NewHTTPClient(ctx context.Context) *HttpClient {
  client := &HttpClient{
    ctx:     ctx,
    limiter: semaphore.NewWeighted(httpTreadsCount),
    client: &http.Client{
      Timeout: httpRequestTimeout,
    },
  }
  return client
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

func readResponseBody(body io.Reader) ([]byte, error) {
  b, err := io.ReadAll(body)
  if err != nil {
    return nil, fmt.Errorf("cannot read response body: %v", err)
  }
  return b, nil
}
