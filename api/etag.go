package api

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"net/http"

	"github.com/wtg/shuttletracker/log"
)

type etagResponseWriter struct {
	http.ResponseWriter
	buf  bytes.Buffer
	hash hash.Hash
	w    io.Writer
}

func (e *etagResponseWriter) Write(p []byte) (int, error) {
	return e.w.Write(p)
}

func etag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ew := &etagResponseWriter{
			ResponseWriter: w,
			buf:            bytes.Buffer{},
			hash:           sha1.New(),
		}
		ew.w = io.MultiWriter(&ew.buf, ew.hash)

		next.ServeHTTP(ew, r)

		sum := fmt.Sprintf("%x", ew.hash.Sum(nil))
		w.Header().Set("ETag", sum)

		if r.Header.Get("If-None-Match") == sum {
			w.WriteHeader(304)
		} else {
			_, err := ew.buf.WriteTo(w)
			if err != nil {
				log.WithError(err).Error("unable to write HTTP response")
			}
		}
	})
}
