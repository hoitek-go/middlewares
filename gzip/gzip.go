package gzip

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
)

const acceptEncoding string = "Accept-Encoding"

type (
	compressResponseWriter struct {
		compressor io.Writer
		w          http.ResponseWriter
	}
	flusher interface {
		Flush() error
	}
)

func (cw *compressResponseWriter) WriteHeader(c int) {
	cw.w.Header().Del("Content-Length")
	cw.w.WriteHeader(c)
}

func (cw *compressResponseWriter) Write(b []byte) (int, error) {
	h := cw.w.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}
	h.Del("Content-Length")

	return cw.compressor.Write(b)
}

func (cw *compressResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	return io.Copy(cw.compressor, r)
}

func (w *compressResponseWriter) Flush() {
	if f, ok := w.compressor.(flusher); ok {
		f.Flush()
	}
	if f, ok := w.w.(http.Flusher); ok {
		f.Flush()
	}
}

func Middleware(level int) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		if level < gzip.DefaultCompression || level > gzip.BestCompression {
			level = gzip.DefaultCompression
		}
		const (
			gzipEncoding  = "gzip"
			flateEncoding = "deflate"
		)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var encoding string

			for _, curEnc := range strings.Split(r.Header.Get(acceptEncoding), ",") {
				curEnc = strings.TrimSpace(curEnc)
				if curEnc == gzipEncoding || curEnc == flateEncoding {
					encoding = curEnc
					break
				}
			}

			w.Header().Add("Vary", acceptEncoding)

			if encoding == "" {
				h.ServeHTTP(w, r)
				return
			}

			if r.Header.Get("Upgrade") != "" {
				h.ServeHTTP(w, r)
				return
			}

			var encWriter io.WriteCloser
			if encoding == gzipEncoding {
				encWriter, _ = gzip.NewWriterLevel(w, level)
			} else if encoding == flateEncoding {
				encWriter, _ = flate.NewWriter(w, level)
			}
			defer encWriter.Close()

			w.Header().Set("Content-Encoding", encoding)
			r.Header.Del(acceptEncoding)

			cw := &compressResponseWriter{
				w:          w,
				compressor: encWriter,
			}

			w = httpsnoop.Wrap(w, httpsnoop.Hooks{
				Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
					return cw.Write
				},
				WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return cw.WriteHeader
				},
				Flush: func(httpsnoop.FlushFunc) httpsnoop.FlushFunc {
					return cw.Flush
				},
				ReadFrom: func(rff httpsnoop.ReadFromFunc) httpsnoop.ReadFromFunc {
					return cw.ReadFrom
				},
			})

			h.ServeHTTP(w, r)
		})
	}
}
