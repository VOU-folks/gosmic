package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"reflect"
)

type Response = http.Response

// CLIENT
type Client struct {
	UserAgent string
	Transport *http.Client
	Jar       *cookiejar.Jar
}

type ClientOptions struct {
	UserAgent           string
	SkipTlsVerification bool
	Jar                 *cookiejar.Jar
}

func NoTlsVerificationOption() ClientOptions {
	return ClientOptions{
		SkipTlsVerification: true,
	}
}

const defaultUserAgent = "NaviMax/gosmic/1.0"

func NewClient() *Client {
	defaultOptions := ClientOptions{
		UserAgent:           defaultUserAgent,
		SkipTlsVerification: false,
		Jar:                 nil,
	}
	return NewClientWithOptions(defaultOptions)
}

func NewClientWithOptions(options ClientOptions) *Client {
	userAgent := defaultUserAgent
	if options.UserAgent != "" {
		userAgent = options.UserAgent
	}

	jar, _ := cookiejar.New(nil)
	if options.Jar != nil {
		jar = options.Jar
	}

	skipTlsCheck := options.SkipTlsVerification

	return &Client{
		UserAgent: userAgent,
		Transport: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: skipTlsCheck,
				},
			},
		},
		Jar: jar,
	}
}

func (c *Client) WithUserAgent(userAgent string) *Client {
	c.UserAgent = userAgent

	return c
}

func (c *Client) ResetCookies() {
	c.Jar, _ = cookiejar.New(nil)
}

// REQUEST
type Request struct {
	client *Client

	Method  RequestMethod
	Url     string
	Headers map[string]string
	Body    *bytes.Buffer
	IsJson  bool
}

type RequestMethod string

var (
	POST    RequestMethod = "POST"
	GET     RequestMethod = "GET"
	PUT     RequestMethod = "PUT"
	PATCH   RequestMethod = "PATCH"
	DELETE  RequestMethod = "DELETE"
	HEAD    RequestMethod = "HEAD"
	OPTIONS RequestMethod = "OPTIONS"
)

var defaultRequestMethod = GET

func NewRequest(client *Client) *Request {
	r := &Request{
		client: client,

		Method:  defaultRequestMethod,
		Url:     "",
		Headers: map[string]string{},
		IsJson:  false,
		Body:    nil,
	}

	r.WithUserAgent(client.UserAgent)

	return r
}

func (r *Request) WithMethod(method RequestMethod) *Request {
	r.Method = method

	return r
}

func (r *Request) WithUrl(url string) *Request {
	r.Url = url

	return r
}

func (r *Request) WithHeaders(headers map[string]string) *Request {
	r.Headers = headers

	return r
}

func (r *Request) SetHeader(key string, value string) *Request {
	r.Headers[key] = value

	return r
}

func (r *Request) WithBody(body *bytes.Buffer) *Request {
	r.Body = body

	return r
}

func (r *Request) WithJSONBody(body any) *Request {
	buff := JSONBodyToBuffer(body)
	r.IsJson = true
	r.Headers["Content-Type"] = "application/json"

	return r.WithBody(buff)
}

func JSONBodyToBuffer(body any) *bytes.Buffer {
	switch reflect.ValueOf(body).Kind() {
	case
		reflect.Struct,
		reflect.Map,
		reflect.Slice:
		return SerializableJSONToBuffer(body)

	case reflect.String:
		return StringToBuffer(body.(string))

	default:
		return StringToBuffer("{}")
	}
}

func SerializableJSONToBuffer(body any) *bytes.Buffer {
	serializedBody := new(bytes.Buffer)
	_ = json.NewEncoder(serializedBody).Encode(body)
	return serializedBody
}

func StringToBuffer(body string) *bytes.Buffer {
	return bytes.NewBufferString(body)
}

func (r *Request) WithUserAgent(userAgent string) *Request {
	r.Headers["User-Agent"] = userAgent

	return r
}

func (r *Request) Send() (*http.Response, error) {
	var bodyReader io.Reader = bytes.NewReader([]byte{})
	if r.Body != nil {
		bodyReader = bytes.NewReader(r.Body.Bytes())
	}

	req, err := http.NewRequest(string(r.Method), r.Url, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", r.client.UserAgent)
	}

	return SendRequest(req, r.client.Transport)
}

func SendRequest(req *http.Request, transport *http.Client) (*http.Response, error) {
	var err error
	var res *http.Response

	res, err = transport.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func PrintResponseHeaders(response *http.Response, printFn func(...any)) {
	for k, v := range response.Header {
		printFn(k, v)
	}
}

func PrintResponseBody(response *http.Response, printFn func(...any)) {
	body := new(bytes.Buffer)
	_, _ = body.ReadFrom(response.Body)
	printFn(body.String())
}

func GetResponseBodyAsString(response *http.Response) string {
	body := new(bytes.Buffer)
	_, _ = body.ReadFrom(response.Body)
	return body.String()
}

func GetResponseBodyAsBytes(response *http.Response) []byte {
	body := new(bytes.Buffer)
	_, _ = body.ReadFrom(response.Body)
	return body.Bytes()
}
