package easyhttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

var client *http.Client

func init() {
	def := http.DefaultTransport
	defPot, ok := def.(*http.Transport)
	if !ok {
		panic("init transport出错")
	}
	defPot.MaxIdleConns = 100
	defPot.MaxIdleConnsPerHost = 100
	defPot.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client = &http.Client{
		Timeout:   time.Second * time.Duration(20),
		Transport: defPot,
	}
}

func Get(url string, header map[string]string, params map[string]interface{}) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	q := req.URL.Query()
	if params != nil {
		for Key, val := range params {
			v, _ := toString(val)
			q.Add(Key, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		return nil, err
	}
	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func Post(url string, header map[string]string, param map[string]interface{}) ([]byte, error) {
	dd, _ := json.Marshal(param)
	re := bytes.NewReader(dd)
	req, err := http.NewRequest("POST", url, re)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func PostMultipart(url string, header map[string]string, payload *bytes.Buffer) ([]byte, error) {
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	bb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

func toString(i interface{}) (string, error) {
	i = indirectToStringerOrError(i)

	switch s := i.(type) {
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(s), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case uint:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint64:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(s), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(s), 10), nil
	case []byte:
		return string(s), nil
	case template.HTML:
		return string(s), nil
	case template.URL:
		return string(s), nil
	case template.JS:
		return string(s), nil
	case template.CSS:
		return string(s), nil
	case template.HTMLAttr:
		return string(s), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return s.String(), nil
	case error:
		return s.Error(), nil
	default:
		return "", fmt.Errorf("unable to cast %#v of type %T to string", i, i)
	}
}

func indirectToStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}

	var errorType = reflect.TypeOf((*error)(nil)).Elem()
	var fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
