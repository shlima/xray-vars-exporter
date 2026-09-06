package web

import (
	"resty.dev/v3"
)

type Web struct {
	proxy        *resty.Client
	proxyAddress string
}

func NewWeb() *Web {
	return &Web{
		proxy: resty.New(),
	}
}

func (w *Web) SetProxyAddress(input string) {
	w.proxyAddress = input
}
