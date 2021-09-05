package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
	"github.com/valyala/fasttemplate"
)

// TestConfig sets up the app for a test environment.
func TestConfig(t *testing.T) Configurer {
	return func(app *App) {
		app.Debug = true
		app.Logger = _loggerWithConfig(_testLogWriter{t})
	}
}

type _testLogWriter struct{ *testing.T }

func (tl _testLogWriter) Write(p []byte) (int, error) {
	tl.Logf(string(p))
	return len(p), nil
}

// Extracted from github.com/labstack/echo/v4/middleware LoggerWithConfig to
// allow padding the method for easier to read output.
func _loggerWithConfig(output io.Writer) echo.MiddlewareFunc {
	skipper := func(c echo.Context) bool {
		return strings.HasPrefix(c.Path(), "/css")
	}
	format := "${status} ${method} ${uri}\n"
	customTimeFormat := "2006-01-02 15:04:05.00000"
	template := fasttemplate.New(format, "${", "}")
	colorer := color.New()
	colorer.SetOutput(output)
	pool := &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			buf := pool.Get().(*bytes.Buffer)
			buf.Reset()
			defer pool.Put(buf)

			if _, err = template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
				switch tag {
				case "time_unix":
					return buf.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
				case "time_unix_nano":
					return buf.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
				case "time_rfc3339":
					return buf.WriteString(time.Now().Format(time.RFC3339))
				case "time_rfc3339_nano":
					return buf.WriteString(time.Now().Format(time.RFC3339Nano))
				case "time_custom":
					return buf.WriteString(time.Now().Format(customTimeFormat))
				case "id":
					id := req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}
					return buf.WriteString(id)
				case "remote_ip":
					return buf.WriteString(c.RealIP())
				case "host":
					return buf.WriteString(req.Host)
				case "uri":
					return buf.WriteString(req.RequestURI)
				case "method":
					return buf.WriteString(fmt.Sprintf("%-8s", req.Method))
				case "path":
					p := req.URL.Path
					if p == "" {
						p = "/"
					}
					return buf.WriteString(p)
				case "protocol":
					return buf.WriteString(req.Proto)
				case "referer":
					return buf.WriteString(req.Referer())
				case "user_agent":
					return buf.WriteString(req.UserAgent())
				case "status":
					n := res.Status
					s := colorer.Green(n)
					switch {
					case n >= 500:
						s = colorer.Red(n)
					case n >= 400:
						s = colorer.Yellow(n)
					case n >= 300:
						s = colorer.Cyan(n)
					}
					return buf.WriteString(s)
				case "error":
					if err != nil {
						// Error may contain invalid JSON e.g. `"`
						b, _ := json.Marshal(err.Error())
						b = b[1 : len(b)-1]
						return buf.Write(b)
					}
				case "latency":
					l := stop.Sub(start)
					return buf.WriteString(strconv.FormatInt(int64(l), 10))
				case "latency_human":
					return buf.WriteString(stop.Sub(start).String())
				case "bytes_in":
					cl := req.Header.Get(echo.HeaderContentLength)
					if cl == "" {
						cl = "0"
					}
					return buf.WriteString(cl)
				case "bytes_out":
					return buf.WriteString(strconv.FormatInt(res.Size, 10))
				default:
					switch {
					case strings.HasPrefix(tag, "header:"):
						return buf.Write([]byte(c.Request().Header.Get(tag[7:])))
					case strings.HasPrefix(tag, "query:"):
						return buf.Write([]byte(c.QueryParam(tag[6:])))
					case strings.HasPrefix(tag, "form:"):
						return buf.Write([]byte(c.FormValue(tag[5:])))
					case strings.HasPrefix(tag, "cookie:"):
						cookie, err := c.Cookie(tag[7:])
						if err == nil {
							return buf.Write([]byte(cookie.Value))
						}
					}
				}
				return 0, nil
			}); err != nil {
				return
			}

			_, err = output.Write(buf.Bytes())
			return
		}
	}
}