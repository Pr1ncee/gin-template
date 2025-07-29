/*
Package utils provide utilities for different kind of operations.

Specifically, this file implements utils for caching middleware.
*/
package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

// ResponseBodyWriter wraps the actual response in this object to be able to analyze it further.
type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

// Write overrides Write method of a bytes.Buffer object.
func (r *ResponseBodyWriter) Write(b []byte) (int, error) {
	r.Body.Write(b)
	return r.ResponseWriter.Write(b)
}
