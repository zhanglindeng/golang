package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"time"
)

type myResponse struct {
	data []byte
}

type myResponseWriter struct {
	gin.ResponseWriter
	status  int
	written bool
	*myResponse
}

func (w *myResponseWriter) WriteHeader(code int) {
	w.status = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *myResponseWriter) Status() int {
	return w.ResponseWriter.Status()
}

func (w *myResponseWriter) Written() bool {
	return w.written
}

func (w *myResponseWriter) Write(data []byte) (int, error) {
	ret, err := w.ResponseWriter.Write(data)
	if err == nil {
		//save response
		w.data = append(w.data, data[0:len(data)]...)
		if err != nil {
			// need logger
			log.Println("Write", err)
		}
	}
	return ret, err
}

func (w *myResponseWriter) WriteString(str string) (int, error) {
	ret, err := w.ResponseWriter.WriteString(str)
	if err == nil {
		//save response
		data := []byte(str)
		w.data = append(w.data, data[0:len(data)]...)
		if err != nil {
			// need logger
			log.Println("WriteString", err)
		}
	}
	return ret, err
}

func newMyResponseWriter(writer gin.ResponseWriter, response *myResponse) *myResponseWriter {
	return &myResponseWriter{writer, 0, false, response}
}

func setStartTime(c *gin.Context) {
	c.Set("startTime", time.Now().UnixNano())
}

func replaceResponseWriter(c *gin.Context) {
	response := &myResponse{}
	c.Writer = newMyResponseWriter(c.Writer, response)
	c.Set("myResponse", response)
}

func afterRequest(c *gin.Context) {
	c.Next()
	startTime := c.MustGet("startTime").(int64)
	exeTime := time.Now().UnixNano() - startTime
	log.Println("exeTime", exeTime)
	log.Println(c.Writer.Size(), c.Writer.Status(), c.Writer.Header().Get("Content-Type"))
	response := c.MustGet("myResponse").(*myResponse)
	log.Println(string(response.data))
}
