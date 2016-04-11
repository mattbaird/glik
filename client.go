// Copyright 2016 Matthew Baird
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package glik

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AtScaleInc/apps-shared/httputil"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const content_type_header = "Content-Type"
const content_length_header = "Content-Length"
const xrf_header = "x-qlik-xrfkey"
const qlik_user_header = "X-Qlik-User"
const application_json_content_type = "application/xml"
const user_header_value = "UserDirectory=Internal; UserId=%s"
const POST = "POST"
const GET = "GET"
const DELETE = "DELETE"

var ErrDoesNotExist = errors.New("Does Not Exist")

var API_VERSION = "2.2"

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-About-API-Get-Description.htm
func (api *API) About() (About, error) {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s/qrs/about?xrfkey=%s", api.Server, xrfKey)
	var retval About
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = fmt.Sprintf(user_header_value, api.QlikUser)
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, GET, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval, err
}

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-App-Publish.htm
func (api *API) Publish(appId, streamId, name string) (PublishResults, error) {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s/qrs/app/%s/publish?stream=%s&name=%s&xrfkey=%s", api.Server, appId, streamId, name, xrfKey)
	var retval PublishResults
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = fmt.Sprintf(user_header_value, api.QlikUser)
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, GET, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval, err
}

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-App-Publish.htm
func (api *API) List() ([]AppListResults, error) {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s/qrs/app?xrfkey=%s", api.Server, xrfKey)
	var retval []AppListResults
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = fmt.Sprintf(user_header_value, api.QlikUser)
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, GET, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval, err
}

func makeXrfKey() string {
	id := uuid.NewV4().String()
	idNoDashes := strings.Replace(id, "-", "", -1)
	return idNoDashes[0:16]
}

func (api *API) makeRequest(requestUrl string, method string, payload []byte, result interface{}, headers map[string]string,
	cTimeout time.Duration, rwTimeout time.Duration) error {
	if false {
		fmt.Printf("%s:%v\n", method, requestUrl)
		if payload != nil {
			fmt.Printf("%v\n", string(payload))
		}
	}
	client := httputil.NewTimeoutClient(cTimeout, rwTimeout)
	var req *http.Request
	if len(payload) > 0 {
		var httpErr error
		req, httpErr = http.NewRequest(strings.TrimSpace(method), strings.TrimSpace(requestUrl), bytes.NewBuffer(payload))
		if httpErr != nil {
			return httpErr
		}
		req.Header.Add(content_length_header, strconv.Itoa(len(payload)))
	} else {
		var httpErr error
		req, httpErr = http.NewRequest(strings.TrimSpace(method), strings.TrimSpace(requestUrl), nil)
		if httpErr != nil {
			return httpErr
		}
	}
	if headers != nil {
		for header, headerValue := range headers {
			req.Header.Add(header, headerValue)
		}
	}
	var httpErr error
	resp, httpErr := client.Do(req)
	if httpErr != nil {
		return httpErr
	}
	defer resp.Body.Close()
	body, readBodyError := ioutil.ReadAll(resp.Body)
	if readBodyError != nil {
		return readBodyError
	}
	if false {
		fmt.Printf("body:%s", body)
	}
	if resp.StatusCode == 404 {
		return ErrDoesNotExist
	}
	if resp.StatusCode >= 300 {
		//		tErrorResponse := ErrorResponse{}
		//		err := xml.Unmarshal(body, &tErrorResponse)
		//		if err != nil {
		//			return err
		//		}
		//		return tErrorResponse.Error
		return fmt.Errorf("Error during request [%v]:%v", resp.StatusCode, resp.Status)
	}
	if result != nil {
		// else unmarshall to the result type specified by caller
		err := json.Unmarshal(body, &result)
		if err != nil {
			return err
		}
	}
	return nil
}
