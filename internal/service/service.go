package service

import "net/http"

type Starter interface {
	Start() error
}

type Stopper interface {
	Stop() error
}

type HTTPRoute struct {
	Method  string
	Path    string
	Handler http.Handler
}
