package doggy

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/hnlq715/doggy/utils"

	"go.uber.org/zap"
)

var (
	defaultClient = &http.Client{}
)

type Request struct {
	url string
	ctx context.Context
	req *http.Request
}

func newRequest(ctx context.Context, method, url string, body []byte) *Request {
	l := utils.LogFromContext(ctx).With(zap.String("query", url), zap.String("type", "http"), zap.String("direction", "out"))

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		l.Error("http.NewRequest failed", zap.Error(err))
		req = &http.Request{
			URL:        nil,
			Method:     method,
			Header:     make(http.Header),
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
	}

	return &Request{
		ctx: ctx,
		url: url,
		req: req,
	}
}

func Get(ctx context.Context, url string) *Request {
	return newRequest(ctx, "GET", url, nil)
}

func Post(ctx context.Context, url string, body []byte) *Request {
	return newRequest(ctx, "POST", url, body)
}

func (r *Request) Bytes() ([]byte, error) {
	resp, err := defaultClient.Do(r.req.WithContext(r.ctx))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return data, nil
}

func (r *Request) String() (string, error) {
	data, err := r.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *Request) ToJSON(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (r *Request) ToXML(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}

	return xml.Unmarshal(data, v)
}
