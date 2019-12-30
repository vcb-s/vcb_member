package helper

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
)

// RequestParam 请求参数
type RequestParam struct {
	Method string
	URL    string
	Body   []byte
	Header http.Header
	Cert   string
	Key    string
	IP     string
}

// Fetch 请求
func Fetch(args RequestParam) ([]byte, error) {
	client := &http.Client{}
	// set tls
	if args.Cert != "" && args.Key != "" {
		client = setTLS(args)
	}
	// set proxy
	if args.IP != "" {
		client = setProxy(args)
	}
	// set request
	req, err := http.NewRequest(args.Method, args.URL, bytes.NewReader(args.Body))
	if err != nil {
		return nil, nil
	}
	req.Close = true
	req.Header = args.Header
	// get response
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func setTLS(args RequestParam) *http.Client {
	cert, err := tls.LoadX509KeyPair(args.Cert, args.Key)
	if err != nil {
		return nil
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},
	}
	transport := &http.Transport{
		TLSClientConfig: config,
	}
	return &http.Client{
		Transport: transport,
	}
}

func setProxy(args RequestParam) *http.Client {
	proxy := func(*http.Request) (*url.URL, error) {
		return url.Parse(args.IP)
	}
	transport := &http.Transport{
		Proxy: proxy,
	}
	return &http.Client{
		Transport: transport,
	}
}
