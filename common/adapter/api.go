package adapter

import (
	"app/utils"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/fatih/structs"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

var (
	gCache                 = cache.New(60*time.Minute, 1*time.Minute)
	logFile *utils.LogFile = utils.NewLogFile()
)

// APIs ..
type APIs map[string]API

// Get API
func (apis APIs) Get(name string) (result API) {
	if adapter, ok := apis[name]; ok {
		result = adapter
		result.Name = name
	} else {
		panic("Không tìm thấy config API " + name)
	}
	return
}

// API ..
type API struct {
	Name   string
	URL    string `mapstructure:"url"`
	Method string `mapstructure:"method"`
	Type   string `mapstructure:"content-type"`
	Auth   struct {
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"auth"`
	Headers    []header `mapstructure:"header"`
	Params     []param  `mapstructure:"params"`
	Cache      int      `mapstructure:"cache"`
	Timeout    int      `mapstructure:"timeout"`
	KeepLog    bool     `mapstructure:"keep_log"`
	SecretKey  string   `mapstructure:"secret_key"`
	ClientKey  string   `mapstructure:"client_key"`
	Agent      string   `mapstructure:"agent"`
	StatusCode int
}

type header struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

type param struct {
	Key   string      `mapstructure:"key"`
	Value interface{} `mapstructure:"value"`
}

// SetParams ..
func (api *API) SetParams(params map[string]interface{}) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		api.Params = append(api.Params, param{Key: k, Value: params[k]})
	}
}

// SetHeader ..
func (api *API) SetHeader(headers map[string]string) {
	for k, v := range headers {
		api.Headers = append(api.Headers, header{Key: k, Value: v})
	}
}

// SetRequestParamsWithStructsTag sử dụng để gắn một struct bất kỳ vào, chuyển thành map[string]interface với key dựa theo tag structs
func (api *API) SetRequestParamsWithStructsTag(params interface{}) {
	apiParams := structs.Map(params)
	api.SetParams(apiParams)
}

// SetRequestParamsWithJSONTag sử dụng để gắn một struct bất kỳ vào, chuyển thành map[string]interface với key dựa theo tag json
func (api *API) SetRequestParamsWithJSONTag(data interface{}) {
	res2B, _ := json.Marshal(data)
	type Params map[string]interface{}
	var params Params
	json.Unmarshal(res2B, &params)
	api.SetParams(params)
}

// SetParams ..
func (api *API) getParams() map[string]interface{} {
	params := make(map[string]interface{})
	for _, v := range api.Params {
		params[v.Key] = v.Value
	}
	return params
}

var (
	onceAPI      map[int]*sync.Once
	httpClient   map[int]*http.Client
	onceAPIMutex = sync.RWMutex{}
	transport    *http.Transport
)

func init() {
	transport = http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	onceAPI = make(map[int]*sync.Once)
	httpClient = make(map[int]*http.Client)
}

func (api *API) newHTTPClient() *http.Client {
	onceAPIMutex.Lock()

	if api.Timeout == 0 {
		api.Timeout = 10
	}

	if onceAPI[api.Timeout] == nil {
		onceAPI[api.Timeout] = &sync.Once{}
	}

	onceAPI[api.Timeout].Do(func() {
		httpClient[api.Timeout] = &http.Client{
			Timeout:   time.Duration(api.Timeout) * time.Second,
			Transport: transport,
		}
	})

	onceAPIMutex.Unlock()

	return httpClient[api.Timeout]
}

// MakeRequest ..
func (api *API) MakeRequest(requestID string) (body []byte, err error) {

	//log
	bodyRequest, _ := json.Marshal(api.getParams())

	bodyRequestLimit := string(bodyRequest)
	if len(bodyRequestLimit) > 5000 {
		bodyRequestLimit = bodyRequestLimit[:5000]
	}

	data := map[string]interface{}{
		"method":       api.Method,
		"headers":      api.Headers,
		"request-id":   requestID,
		"uri":          api.URL,
		"body-request": bodyRequestLimit,
	}

	defer api.send(data, true)

	//end log

	start := time.Now()

	key := hash(*api)

	cacheData, ok := gCache.Get(key)
	if ok {
		body = cacheData.([]byte)

		bodyResponseLimit := string(body)
		if len(bodyResponseLimit) > 5000 {
			bodyResponseLimit = bodyResponseLimit[:5000]
		}

		data["body-response"] = bodyResponseLimit
		data["status"] = http.StatusOK
		data["latency-human"] = time.Now().Sub(start).String()
		data["latency-micro"] = time.Now().Sub(start).Microseconds()
		data["cache"] = true

		return
	}

	req := new(http.Request)

	if api.Method == "" {
		api.Method = "GET"
	}

	switch api.Method {
	case http.MethodGet:
		req, err = http.NewRequest(api.Method, api.URL, nil)
		if err != nil {
			return
		}
		q := req.URL.Query()
		for _, param := range api.Params {
			q.Add(param.Key, fmt.Sprintf("%v", param.Value))
		}
		req.URL.RawQuery = q.Encode()
	case http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete:
		bodyRequest := &bytes.Buffer{}
		json.NewEncoder(bodyRequest).Encode(api.getParams())

		req, err = http.NewRequest(api.Method, api.URL, bodyRequest)
		fmt.Println(req)
		if err != nil {
			return
		}
	}

	body, err = api.doRequest(requestID, req)

	if err != nil {
		data["latency-human"] = time.Now().Sub(start).String()
		data["latency-micro"] = time.Now().Sub(start).Microseconds()
		data["error"] = err
	} else {

		bodyResponseLimit := string(body)
		if len(bodyResponseLimit) > 5000 {
			bodyResponseLimit = bodyResponseLimit[:5000]
		}

		data["body-response"] = bodyResponseLimit
		data["status"] = api.StatusCode
		data["latency-human"] = time.Now().Sub(start).String()
		data["latency-micro"] = time.Now().Sub(start).Microseconds()

		if api.Cache > 0 {
			if api.StatusCode == http.StatusOK {
				gCache.Set(key, body, time.Second*time.Duration(api.Cache))
			} else {
				log.Printf("Can't cache request: %v\n", api.URL)
			}
		}
	}

	return
}

// MakeRequestWithAttachments ..
func (api API) MakeRequestWithAttachments(requestID string) (body []byte, err error) {

	key := hash(api)

	cacheData, ok := gCache.Get(key)
	if ok {
		body = cacheData.([]byte)
		return
	}
	req := new(http.Request)
	bodyRequest := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyRequest)
	for key, val := range api.getParams() {
		valueType := reflect.TypeOf(val)
		if valueType.String() == "[]*multipart.FileHeader" {
			values := val.([]*multipart.FileHeader)
			for _, file := range values {
				rawFile, err := file.Open()
				if err != nil {
					return nil, err
				}
				defer rawFile.Close()

				fileWriter, err := bodyWriter.CreateFormFile(key, file.Filename)
				if err != nil {
					return nil, err
				}
				if _, err = io.Copy(fileWriter, rawFile); err != nil {
					return nil, err
				}
			}
		} else {
			if valueType.String() == "string" {
				newVal := val.(string)
				err = bodyWriter.WriteField(key, newVal)
				if err != nil {
					return
				}
			} else {
				var wt io.Writer
				value, _ := json.Marshal(val)
				wt, err = bodyWriter.CreateFormField(key)
				if err != nil {
					return
				}
				if _, err = io.Copy(wt, bytes.NewReader(value)); err != nil {
					return
				}
			}
		}
	}
	bodyWriter.Close()
	req, err = http.NewRequest(api.Method, api.URL, bodyRequest)
	req.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	return api.doRequest(requestID, req)
}

// MakeRequestURLEncoded ..
func (api *API) MakeRequestURLEncoded(requestID string) (body []byte, err error) {

	req := new(http.Request)

	if api.Method == "" {
		api.Method = "POST"
	}

	bodyRequest := &bytes.Buffer{}

	data := url.Values{}
	for k, v := range api.getParams() {
		data.Set(k, v.(string))
	}
	bodyRequest.WriteString(data.Encode())

	req, err = http.NewRequest(api.Method, api.URL, bodyRequest)
	if err != nil {
		return
	}

	return api.doRequest(requestID, req)
}

func (api *API) doRequest(requestID string, req *http.Request) (body []byte, err error) {

	// check api URL is null
	if api.URL == "" {
		err = fmt.Errorf("Empty URL is invalid")
		return
	}

	if api.Auth.UserName != "" {
		req.SetBasicAuth(api.Auth.UserName, api.Auth.Password)
	}

	for _, header := range api.Headers {
		if header.Key != "" {
			req.Header.Set(header.Key, header.Value)
		}
	}

	if api.Agent != "" {
		req.Header.Set("User-Agent", api.Agent)
	}

	//reset requestID
	req.Header.Set(echo.HeaderXRequestID, requestID)

	req.Header.Set(echo.HeaderAcceptEncoding, "gzip")

	client := api.newHTTPClient()

	response, errDoRequest := client.Do(req)
	if errDoRequest != nil {
		err = errDoRequest
		return
	}

	api.StatusCode = response.StatusCode

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	body, err = ioutil.ReadAll(reader)

	response.Body.Close()

	return
}

func hash(data API) string {
	arrBytes := []byte{}

	jsonBytes, _ := json.Marshal(data)
	arrBytes = append(arrBytes, jsonBytes...)

	return fmt.Sprintf("%x", md5.Sum(arrBytes))
}

// Send, Use this function to send data to elastic
func (api API) send(value interface{}, keep bool) {
	data := map[string]interface{}{
		"log":  value,
		"keep": strconv.FormatBool(keep),
	}

	jsonData, _ := json.Marshal(data)
	logFile.Write(string(jsonData))

	return
}
