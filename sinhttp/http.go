package sinhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sin-z/sin-common/sinerrors"
	"github.com/sin-z/sin-common/sinlog"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"

	"net/url"
	"strings"
	"time"
)

const (
	_jumpInfoKey  = "_NVWA_HTTP_JUMP_INFO_"
	_appKeyHeader = "uberctx-_namespace_appkey_"
)

type requestImpl struct {
	ctx         context.Context
	serviceName string

	method string
	uri    string

	option RequestOption

	query      url.Values
	reqBodyRaw interface{}
	body       io.Reader

	req *http.Request
	rsp *http.Response

	err error
}
type requestDetailImpl struct {
	impl *requestImpl
}

// New return *Req, don't reuse Req
func Client(ctx context.Context, serviceName string) *requestImpl {
	return &requestImpl{
		ctx:         ctx,
		serviceName: serviceName,
		query:       url.Values{},
		req:         nil,
	}
}

// Get execute a HTTP GET request
func (r *requestImpl) Get(uri string) *requestDetailImpl {
	if r.err != nil {
		return &requestDetailImpl{
			impl: r,
		}
	}
	// todo get service name ip
	r.req, r.err = http.NewRequest("GET", uri, nil)
	if r.err != nil {
		return &requestDetailImpl{
			impl: r,
		}
	}
	r.method = http.MethodGet
	r.uri = uri
	return &requestDetailImpl{
		impl: r,
	}
}

// Post execute a HTTP POST request
func (r *requestImpl) Post(uri string, body interface{}) *requestDetailImpl {
	if r.err != nil {
		return &requestDetailImpl{
			impl: r,
		}
	}
	if body != nil {
		r.withBody(body)
		if r.err != nil {
			return &requestDetailImpl{
				impl: r,
			}
		}
	}
	// todo get service name ip
	r.req, r.err = http.NewRequest("POST", uri, r.body)
	if r.err != nil {
		return &requestDetailImpl{
			impl: r,
		}
	}
	r.method = http.MethodPost
	r.uri = uri
	return &requestDetailImpl{
		impl: r,
	}
}

// set request body, support ioReader, string, []byte, otherwise use json.Marshal
// body：如为struct，则将自动解析成json，否则直接流式放入body
func (r *requestImpl) withBody(body interface{}) *requestImpl {
	r.reqBodyRaw = body
	switch v := body.(type) {
	case io.Reader:
		r.body = v
	case []byte:
		r.body = bytes.NewReader(v)
	case string:
		r.body = strings.NewReader(v)
	default:
		buf, err := json.Marshal(body)
		if err != nil {
			r.err = multierr.Append(r.err, err)
			return r
		}
		r.body = bytes.NewReader(buf)
	}
	return r
}

func (r *requestDetailImpl) WithHeader(k string, v interface{}) *requestDetailImpl {
	r.impl.req.Header.Add(k, fmt.Sprint(v))
	return r
}

func (r *requestDetailImpl) WithHeaderMap(header map[string]interface{}) *requestDetailImpl {
	//r.initOption()
	for k, v := range header {
		//r.option.SetHeader(k, fmt.Sprint(v))
		r.impl.req.Header.Add(k, fmt.Sprint(v))
	}
	return r
}
func (r *requestDetailImpl) WithHeaders(keyAndValues ...interface{}) *requestDetailImpl {
	//r.initOption()
	l := len(keyAndValues) - 1
	//for k, v := range header {
	for i := 0; i < l; i += 2 {
		//r.option.SetHeader(k, fmt.Sprint(v))
		k := fmt.Sprint(keyAndValues[i])
		r.impl.req.Header.Add(k, fmt.Sprint(keyAndValues[i+1]))
	}
	if (l+1)%2 == 1 {
		sinlog.For(r.impl.ctx, zap.String("func", "sinhttp.Client().XXX().WithHeaders")).Warnw("the keys are not aligned")
		k := fmt.Sprint(keyAndValues[l])
		r.impl.req.Header.Add(k, "")
	}
	return r
}

// Param add query param
func (r *requestDetailImpl) WithQueryParam(k string, v interface{}) *requestDetailImpl {
	if r.impl.err != nil {
		return r
	}
	r.impl.query.Add(k, fmt.Sprint(v))
	return r
}
func (r *requestDetailImpl) WithQueryParams(keyAndValues ...interface{}) *requestDetailImpl {
	if r.impl.err != nil {
		return r
	}
	l := len(keyAndValues) - 1
	for i := 0; i < l; i += 2 {
		r.impl.query.Add(fmt.Sprint(keyAndValues[i]), fmt.Sprint(keyAndValues[i+1]))
	}
	if (l+1)%2 == 1 {
		sinlog.For(r.impl.ctx, zap.String("func", "sinhttp.Client().XXX().WithQueryParams")).Warnw("the keys are not aligned")
		r.impl.query.Add(fmt.Sprint(keyAndValues[l]), "")
	}
	return r
}

// Query add query param, support:
// Query("k1=v1&k2=v2")
// Query(url.Values{})
// Query(map[string] string{})
// Query(map[string] interface{})
// Query(struct{}) using url tag, reference: https://github.com/google/go-querystring

// WithQuery add query param, support:
// WithQuery("k1=v1&k2=v2")
// WithQuery(url.Values{})
// WithQuery(map[string] string{})
// WithQuery(map[string] interface{})
// WithQuery(struct{}) using url tag, reference: https://github.com/google/go-querystring
func (r *requestDetailImpl) WithQuery(query interface{}) *requestDetailImpl {
	switch query := query.(type) {
	case string:
		r.impl.queryString(query)
	case []byte:
		r.impl.queryString(string(query))
	case url.Values:
		r.impl.queryUrlValues(query)
	case map[string]string:
		for k, v := range query {
			r.impl.query.Add(k, v)
		}
	case map[string]interface{}:
		for k, v := range query {
			r.impl.query.Add(k, fmt.Sprint(v))
		}
		//slow path
	default:
		r.impl.queryReflect(query)
	}
	return r
}
func (r *requestDetailImpl) initOptions() RequestOption {
	return r.impl.option
}

// 设置超时时间。
// timeout：超时时间，单位毫秒，如不设置，则默认使用配置文件配置，如配置文件亦没有，则无限等待，直至网络异常断开。
// retryTimes：重试次数，如不设置，则默认使用配置文件配置，如配置文件亦没有，则默认0次重试
// slowTime：慢调用时间，单位毫秒，如不设置，则默认使用配置文件配置，如配置文件亦没有，则默认为0，不启用慢时间。
func (r *requestDetailImpl) WithTimeout(timeout int, retryTimes int, slowTime int) *requestDetailImpl {
	opt := r.initOptions()
	opt.Timeout = timeout
	opt.RetryTimes = retryTimes
	opt.SlowTime = slowTime
	return r
}

func (r *requestImpl) toBytes() (*http.Response, error) {
	if len(r.query) != 0 {
		if strings.IndexByte(r.uri, '?') == -1 {
			r.uri += "?" + r.query.Encode()
		} else {
			r.uri += "&" + r.query.Encode()
		}
	}

	client := &http.Client{}
	if r.option.Timeout != 0 {
		client.Timeout = time.Duration(r.option.Timeout) * time.Millisecond
	}
	// TODO url
	rsp, err := client.Do(r.req)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

//
func (r *requestDetailImpl) Do() *responseImpl {
	if r.impl.err != nil {
		return &responseImpl{
			req: r.impl,
		}
	}
	rsp, err := r.impl.toBytes()
	if err != nil {
		r.impl.err = multierr.Append(r.impl.err, err)
		return &responseImpl{
			req: r.impl,
		}
	}
	r.impl.rsp = rsp
	return &responseImpl{
		req: r.impl,
	}
}

type responseImpl struct {
	req *requestImpl
}

func (r *responseImpl) ParseEmpty() error {
	return r.ParseDataJson(nil)
}

// 解析数据
// json结构，解析整个response
func (r *responseImpl) Json(resp interface{}) error {
	if r.req.err != nil {
		return r.req.err
	}
	bodyData, err := ioutil.ReadAll(r.req.rsp.Body)
	if err != nil {
		r.req.err = multierr.Append(r.req.err, err)
		return r.req.err
	}
	err = json.Unmarshal(bodyData, resp)
	if err != nil {
		r.req.err = multierr.Append(r.req.err, err)
		return r.req.err
	}
	return nil
}

// 解析response的data区数据，并根据错误自动返回对应error。反解error，可以使用sinerrors.DMError函数
func (r *responseImpl) ParseDataJson(dataModel interface{}) error {
	if r.req.err != nil {
		return r.req.err
	}

	resp := sinWrapResp{
		WrapResp: WrapResp{
			Data: dataModel,
		},
	}

	err := r.Json(&resp)
	if err != nil {
		r.req.err = multierr.Append(r.req.err, err)
		return r.req.err
	}

	if resp.ErrCode != 0 {
		// 重定向错误码，则存储，以便透传到客户端
		if resp.ErrCode == sinerrors.Codes.Redirect.Code() {
			//context.WithValue(r.req.ctx, _jumpInfoKey, resp.Jump)
			timeout := time.Second * 3
			if r.req.option.Timeout > 0 {
				timeout = time.Millisecond * time.Duration(r.req.option.Timeout)
			}
			if timeout <= 0 { // 没有设置http timeout，则默认10秒
				timeout = time.Second * 3
			} else if timeout > time.Second*10 { // http timeout 超过10秒，则默认保留10秒
				timeout = time.Second * 10
			}
			Session(r.req.ctx).WithTimeout(timeout).Add(_jumpInfoKey, resp.Jump)
		}
		//return errors.Errorf("server response business error code")
		r.req.err = multierr.Append(r.req.err, sinerrors.Get(resp.ErrCode))
		// 注释，通过NewTmpError实现，避免下游服务的错误码不在错误码仓库内
		//return sinerrors.Get(resp.Code)
		return sinerrors.NewTmpError(resp.Code, resp.Msg)
	}
	return nil
}
