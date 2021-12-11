package handlers

import "net/http"

type notFoundRescuingResponseWriter struct {
	http.ResponseWriter
	rescued bool
}

func newNotFoundRescuer(w http.ResponseWriter, toRescue ...int) *notFoundRescuingResponseWriter {
	return &notFoundRescuingResponseWriter{w, false}
}

func (r *notFoundRescuingResponseWriter) Write(data []byte) (int, error) {
	if r.rescued {
		return 0, nil
	}

	return r.ResponseWriter.Write(data)
}

func (r *notFoundRescuingResponseWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusNotFound {
		r.rescued = true
		return
	}

	r.ResponseWriter.WriteHeader(statusCode)
}
