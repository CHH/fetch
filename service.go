package fetch

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net"
	"time"
)

type Config struct {
	MaxIdleConnections int
}

type Request struct {
	Url     string
	Method  string
	Headers map[string][]string
	Body    []byte
	Auth    []string
}

type Response struct {
	Status     string `json:"status"`
	StatusCode int `json:"statusCode"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte `json:"body"`
}

type Service struct {
	client *http.Client
}

func NewService(config Config) *Service {
	transport := &http.Transport{
		MaxIdleConnsPerHost: config.MaxIdleConnections,
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: 30*time.Second,
			KeepAlive: 15*time.Minute,
		}).Dial,
	}

	return &Service{
		client: &http.Client{Transport: transport},
	}
}

func (s *Service) Fetch(req *Request, resp *Response) error {
	var body *bytes.Buffer
	body = bytes.NewBuffer(req.Body)

	httpReq, err := http.NewRequest(req.Method, req.Url, body)

	if err != nil {
		return err
	}

	httpReq.Header = req.Headers
	httpReq.Close = false

	httpResp, err := s.client.Do(httpReq)

	if err != nil {
		return err
	}

	resp.Status = httpResp.Status
	resp.StatusCode = httpResp.StatusCode
	resp.Headers = httpResp.Header

	defer httpResp.Body.Close()
	resp.Body, _ = ioutil.ReadAll(httpResp.Body)

	return nil
}
