package sinhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/gorilla/schema"
	"github.com/sin-z/sin-common/sinerrors"
	"github.com/sin-z/sin-common/sinlog"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"reflect"
)

var (
	_decoder *schema.Decoder
	_valid   = validator.New()
)

func init() {
	_decoder = schema.NewDecoder()
	_decoder.IgnoreUnknownKeys(true)
}

type query struct {
	ctx context.Context
	req *http.Request
	err error
}
type body struct {
	ctx       context.Context
	req       *http.Request
	bodyBytes []byte
	err       error
}
type request struct {
}

var Requests request

func (r *request) Query(ctx context.Context, request *http.Request) *query {
	return &query{
		ctx: ctx,
		req: request,
	}
}

func (r *request) Body(ctx context.Context, request *http.Request) *body {
	return &body{
		ctx: ctx,
		req: request,
	}
}

//
// query
//

func (q *query) Parse(model interface{}) *query {
	if q.err != nil {
		return q
	}
	if err := _decoder.Decode(model, q.req.URL.Query()); err != nil {
		sinlog.For(q.ctx, zap.String("func", "query::Parse"), zap.Any("request", q.req)).Errorw("decode error", "err", err)
		//q.err = errors.Wrapf(ecode.ParamErr, "decode failed,reqUrl(%s),err(%v)", q.req.URL.String(), err)
		q.err = sinerrors.Codes.ClientError.DetailF("error query params(maybe type error): %v", err)
		return q
	}
	if err := _valid.Struct(model); err != nil {
		sinlog.For(q.ctx, zap.String("func", "query::Parse"), zap.Any("query", q.req.URL.String())).Errorw("valid error", "err", err)
		//q.err = errors.Wrapf(ecode.ParamErr, "param check failed,reqUrl(%s),err(%v)", q.req.URL.String(), err)
		q.err = sinerrors.Codes.ClientError.DetailF("error query params: %v", err)
		return q
	}
	return q
}
func (q *query) Atom() (*Atom, error) {
	if q.err != nil {
		return nil, q.err
	}
	atom := &Atom{}
	q.Parse(atom)
	if q.err != nil {
		return nil, q.err
	}
	return atom, nil
}
func (q *query) AtomWeb() (*AtomWeb, error) {
	if q.err != nil {
		return nil, q.err
	}
	atom := &AtomWeb{}
	q.Parse(atom)
	if q.err != nil {
		return nil, q.err
	}
	return atom, nil
}
func (q *query) Error() error {
	return q.err
}

//
// body
//

func (b *body) ParseJson(model interface{}) *body {
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		//b.err = errors.Errorf("output model is not pointer")
		b.err = sinerrors.Codes.ServerError.DetailF("server code error, the pointer type is required, but the value type is passed in")
		return b
	}
	if b.bodyBytes == nil && b.req.Body != nil {
		b.bodyBytes, _ = ioutil.ReadAll(b.req.Body)
		defer func() {
			b.req.Body = ioutil.NopCloser(bytes.NewBuffer(b.bodyBytes))
		}()
	}
	if err := json.Unmarshal(b.bodyBytes, &model); err != nil {
		sinlog.For(b.ctx, zap.String("func", "query::Parse"), zap.Any("request", b.req)).Errorw("decode error", "err", err)
		//b.err = errors.Wrapf(ecode.ParamErr, "decode failed,reqUrl(%s), body(%s),err(%v)", b.req.URL.String(), string(b.bodyBytes), err)
		b.err = sinerrors.Codes.ClientError.DetailF("error body params(maybe type error): %v", err)
		return b
	}
	if err := _valid.Struct(model); err != nil {
		sinlog.For(b.ctx, zap.String("func", "query::Parse"), zap.Any("body", string(b.bodyBytes))).Errorw("valid error", "err", err)
		//b.err = errors.Wrapf(ecode.ParamErr, "param check failed,reqUrl(%s), body(%s),err(%v)", b.req.URL.String(), string(b.bodyBytes), err)
		b.err = sinerrors.Codes.ClientError.DetailF("error body params: %v", err)
		return b
	}
	return b
}

func (b *body) Error() error {
	return b.err
}
