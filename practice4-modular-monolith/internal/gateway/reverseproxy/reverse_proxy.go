package reverseproxy

import (
	"net/http/httputil"
	"net/url"
)

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {

	targetURL, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return proxy, nil
}
