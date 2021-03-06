package ginjsonp

import (
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	noWritten	= -1
	defaultStatus	= 200
)

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		callback := c.DefaultQuery("callback", "")
		if callback == "" {
			callback = c.DefaultQuery("jsonp", "")
			if callback == "" {
				c.Next()
				return
			}
		}
		var wb *responseBuffer
		if w, ok := c.Writer.(gin.ResponseWriter); ok {
			wb = NewResponseBuffer(w)
			c.Writer = wb
			c.Next()
		}else {
			c.Next()
			return
		}
		

		if strings.Index(wb.Header().Get("Content-Type"), "/json") >= 0 {
			status := wb.status
			data := wb.Body.Bytes()
			wb.Body.Reset()

			resp := &jsonpResponse{
				Meta: map[string]interface{}{"status": status},
				Data: data,
			}
			for k, v := range wb.Header() {
				resp.Meta[strings.ToLower(k)] = v[0]
			}
			resp.Meta["content-length"] = len(data)

			body, err := json.Marshal(resp)
			if err != nil {
				panic(err.Error())
			}

			wb.Body.Write([]byte(callback + "("))
			wb.Body.Write(body)
			wb.Body.Write([]byte(")"))

			wb.Header().Set("Content-Type", "application/javascript")
			wb.Header().Set("Content-Length", strconv.Itoa(wb.Body.Len()))
		}

		wb.Flush()
	}
}

type jsonpResponse struct {
	Meta map[string]interface{}
	Data interface{}
}

func (j *jsonpResponse) MarshalJSON() ([]byte, error) {
	meta, err := json.Marshal(j.Meta)
	if err != nil {
		return nil, err
	}
	b := fmt.Sprintf("{\"meta\":%s,\"data\":%s}", meta, j.Data)
	return []byte(b), nil
}

type responseBuffer struct {
	Response gin.ResponseWriter // the actual ResponseWriter to flush to
	status   int                 // the HTTP response code from WriteHeader
	Body     *bytes.Buffer       // the response content body
	Flushed  bool
}

func NewResponseBuffer(w gin.ResponseWriter) *responseBuffer {
	return &responseBuffer{
		Response: w, status: defaultStatus, Body: &bytes.Buffer{},
	}
}

func (w *responseBuffer) Header() http.Header {
	return w.Response.Header() // use the actual response header
}

func (w *responseBuffer) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

func (w *responseBuffer) WriteString(s string) (n int, err error) {
	//w.WriteHeaderNow()
	//n, err = io.WriteString(w.ResponseWriter, s)
	//w.size += n
        n, err = w.Write([]byte(s))
	return
}

func (w *responseBuffer) Written() bool {
	return w.Body.Len() != noWritten
}

func (w *responseBuffer) WriteHeader(status int) {
	w.status = status
}

func (w *responseBuffer) WriteHeaderNow() {
	//if !w.Written() {
	//	w.size = 0
	//	w.ResponseWriter.WriteHeader(w.status)
	//}
}

func (w *responseBuffer) Status() int {
	return w.status
}

func (w *responseBuffer) Size() int {
	return w.Body.Len()
}

func (w *responseBuffer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	//if w.size < 0 {
	//	w.size = 0
	//}
	return w.Response.(http.Hijacker).Hijack()
}

func (w *responseBuffer) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

// Fake Flush
// TBD
func (w *responseBuffer) Flush() {
	w.realFlush()
}

func (w *responseBuffer) realFlush() {
	if w.Flushed {
		return
	}
	w.Response.WriteHeader(w.status)
	if w.Body.Len() > 0 {
		_, err := w.Response.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
		w.Body.Reset()
	}
	w.Flushed = true
}
