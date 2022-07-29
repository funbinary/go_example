package srs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/funbinary/go_example/pkg/errors"
)

func RtcRequest(ctx context.Context, apiPath, rtcUrl, offerSdp string) (string, error) {
	param, err := url.Parse(rtcUrl)
	if err != nil {
		return "", errors.Wrapf(err, "Parse url %v", rtcUrl)
	}

	// Build api url.
	host := param.Host
	if !strings.Contains(host, ":") {
		host += ":1985"
	}

	api := fmt.Sprintf("http://%v", host)
	if !strings.HasPrefix(apiPath, "/") {
		api += "/"
	}
	api += apiPath

	if !strings.HasSuffix(apiPath, "/") {
		api += "/"
	}
	if param.RawQuery != "" {
		api += "?" + param.RawQuery
	}

	// Build JSON body.
	reqBody := struct {
		Api       string `json:"api"`
		ClientIP  string `json:"clientip"`
		SDP       string `json:"sdp"`
		StreamURL string `json:"streamurl"`
	}{
		api, "", offerSdp, rtcUrl,
	}

	resBody := struct {
		Code    int    `json:"code"`
		Session string `json:"sessionid"`
		SDP     string `json:"sdp"`
	}{}

	if err := apiRequest(ctx, api, reqBody, &resBody); err != nil {
		return "", errors.Wrapf(err, "request api=%v", api)
	}

	if resBody.Code != 0 {
		return "", errors.Errorf("Server fail code=%v", resBody.Code)
	}

	return resBody.SDP, nil
}

func apiRequest(ctx context.Context, r string, req interface{}, res interface{}) error {
	var b []byte
	if req != nil {
		if b0, err := json.Marshal(req); err != nil {
			return errors.Wrapf(err, "Marshal body %v", req)
		} else {
			b = b0
		}
	}

	method := "POST"
	if req == nil {
		method = "GET"
	}
	reqObj, err := http.NewRequest(method, r, strings.NewReader(string(b)))
	if err != nil {
		return errors.Wrapf(err, "HTTP request %v", string(b))
	}

	resObj, err := http.DefaultClient.Do(reqObj.WithContext(ctx))
	if err != nil {
		return errors.Wrapf(err, "Do HTTP request %v", string(b))
	}

	b2, err := ioutil.ReadAll(resObj.Body)
	if err != nil {
		return errors.Wrapf(err, "Read response for %v", string(b))
	}

	errorCode := struct {
		Code int `json:"code"`
	}{}
	if err := json.Unmarshal(b2, &errorCode); err != nil {
		return errors.Wrapf(err, "Unmarshal %v", string(b2))
	}
	if errorCode.Code != 0 {
		return errors.Errorf("Server fail code=%v %v", errorCode.Code, string(b2))
	}

	if err := json.Unmarshal(b2, res); err != nil {
		return errors.Wrapf(err, "Unmarshal %v", string(b2))
	}

	return nil
}
