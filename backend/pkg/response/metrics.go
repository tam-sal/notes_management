package response

import "net/http"

type MetricsWriter struct {
	StatusCode      int
	BytesCount      int
	isHeaderWritten bool
	wrapped         http.ResponseWriter
}

func NewMetricsWriter(w http.ResponseWriter) *MetricsWriter {
	return &MetricsWriter{
		StatusCode: http.StatusOK,
		wrapped:    w,
	}
}

func (mxWriter *MetricsWriter) Header() http.Header {
	return mxWriter.wrapped.Header()
}

func (mxWriter *MetricsWriter) WriteHeader(StatusCode int) {
	mxWriter.wrapped.WriteHeader(StatusCode)
	if !mxWriter.isHeaderWritten {
		mxWriter.isHeaderWritten = true
		mxWriter.StatusCode = StatusCode
	}
}

func (mxWriter *MetricsWriter) Write(b []byte) (int, error) {
	mxWriter.isHeaderWritten = true
	n, err := mxWriter.wrapped.Write(b)
	mxWriter.BytesCount += n
	return n, err
}

func (mxWriter *MetricsWriter) Unwrap() http.ResponseWriter {
	return mxWriter.wrapped
}
