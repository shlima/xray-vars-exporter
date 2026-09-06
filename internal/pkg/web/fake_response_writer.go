package web

import (
	"io"
	"net/http"
)

type FakeResponseWriter struct {
	body   io.Writer
	header http.Header
	status int
}

func (f *FakeResponseWriter) Header() http.Header {
	return f.header
}

func (f *FakeResponseWriter) WriteHeader(statusCode int) {
	f.status = statusCode
}

func (f *FakeResponseWriter) Write(data []byte) (int, error) {
	return f.body.Write(data)
}

func NewFakeResponseWriter(writer io.Writer) *FakeResponseWriter {
	return &FakeResponseWriter{
		body:   writer,
		header: make(http.Header),
		status: http.StatusOK,
	}
}
