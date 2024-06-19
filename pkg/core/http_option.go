package core

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type HTTPServerOption func(s *HTTPServer)

// HTTPServerWithName set name of http server
func HTTPServerWithName(name string) HTTPServerOption {
	return func(s *HTTPServer) {
		s.name = name
	}
}

func getServerAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

// HTTPServerWithAddress address of http server
func HTTPServerWithAddress(host string, port int) HTTPServerOption {
	return func(s *HTTPServer) {
		s.address = getServerAddress(host, port)
	}
}

func HTTPServerWithHandler(
	router http.Handler) HTTPServerOption {
	return func(s *HTTPServer) {
		// assign handler
		s.handler = router
	}
}

func HTTPServeWithTimeOutConfig(readTimeout time.Duration, readHeaderTimeout time.Duration, writeTimeout time.Duration, idleTimeout time.Duration) HTTPServerOption {
	return func(s *HTTPServer) {
		s.readTimeout = readTimeout
		s.writeTimeout = writeTimeout
		s.idleTimeout = idleTimeout
		s.readHeaderTimeout = readHeaderTimeout
	}
}

func HTTPServeWithTlsConfig(tlsConfig *tls.Config) HTTPServerOption {
	return func(s *HTTPServer) {
		s.tlsConfig = tlsConfig
	}
}
