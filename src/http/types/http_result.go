package types

import (
	"CASBanSyncer/src/utils"
	"io"
)

type HTTPResult struct {
	Body       io.ReadCloser
	Error      error
	StatusCode int
	cacheRead  []byte
}

func (r *HTTPResult) SetFallback(body []byte) {
	r.cacheRead = body
	r.Body = nil
}

func (r *HTTPResult) Read() []byte {
	if r.cacheRead != nil {
		return r.cacheRead
	}
	if r.Body == nil {
		return nil
	}
	buf, err := utils.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	defer func() {
		r.close()
		r.Body = nil
	}()
	r.cacheRead = buf
	return r.cacheRead
}

func (r *HTTPResult) ReadString() string {
	return string(r.Read())
}

func (r *HTTPResult) close() {
	_ = r.Body.Close()
}
