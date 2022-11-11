package main

import (
	"log"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}

// our http.ResponseWriter implementation
type loggingResponseWriter struct {
	http.ResponseWriter // compose original http.ResponseWriter
	responseData        *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}

func LogginMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var (
			start        = time.Now()
			responseData = &responseData{}
		)

		h.ServeHTTP(&loggingResponseWriter{
			ResponseWriter: rw,
			responseData:   responseData,
		}, req)

		duration := time.Since(start)

		log.Printf("%s[%v] uri:%s duration:%v size:%d", req.Method, responseData.status, req.RequestURI, duration, responseData.size)
	})
}
