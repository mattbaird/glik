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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AtScaleInc/apps-shared/httputil"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const content_type_header = "Content-Type"
const content_length_header = "Content-Length"
const xrf_header = "X-Qlik-Xrfkey"
const qlik_user_header = "X-Qlik-User"
const application_json_content_type = "application/xml"
const user_header_value = "UserDirectory=%s; UserId=%s"
const POST = "POST"
const GET = "GET"
const PUT = "PUT"
const DELETE = "DELETE"

var ErrDoesNotExist = errors.New("Does Not Exist")

var API_VERSION = "2.2"

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-About-API-Get-Description.htm
func (api *API) About() (About, error) {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s://%s:%v/qrs/about?xrfkey=%s", "https", api.Server, api.QrsPort, xrfKey)
	var retval About
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = api.makeQlikUserHeader()
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, GET, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval, err
}

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-App-Publish.htm
func (api *API) Publish(appId, streamId, name string) (ApplicationResult, error) {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s://%s:%v/qrs/app/%s/publish?stream=%s&name=%s&xrfkey=%s", "https", api.Server, api.QrsPort, appId, streamId, url.QueryEscape(name), xrfKey)
	var retval ApplicationResult
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = api.makeQlikUserHeader()
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, PUT, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval, err
}

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-App-Make-Copy.htm
func (api *API) Copy(appId, name string) (ApplicationResult, error) {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s://%s:%v/qrs/app/%s/copy?name=%s&xrfkey=%s", "https", api.Server, api.QrsPort, appId, url.QueryEscape(name), xrfKey)
	fmt.Printf("url:%v\n", url)
	var retval ApplicationResult
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = api.makeQlikUserHeader()
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, POST, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval, err
}

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-App-Publish.htm
func (api *API) List() ([]ApplicationResult, error) {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s://%s:%v/qrs/app?xrfkey=%s", "https", api.Server, api.QrsPort, xrfKey)
	var retval []ApplicationResult
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = api.makeQlikUserHeader()
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, GET, nil, &retval, headers, connectTimeOut, readWriteTimeout)
	return retval, err
}

//http://help.qlik.com/en-US/sense-developer/2.2/Subsystems/RepositoryServiceAPI/Content/RepositoryServiceAPI/RepositoryServiceAPI-App-Reload.htm
func (api *API) Reload(appId string) error {
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s://%s:%v/qrs/app/%s/reload?xrfkey=%s", "https", api.Server, api.QrsPort, appId, xrfKey)
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = api.makeQlikUserHeader()
	headers[content_type_header] = application_json_content_type
	err := api.makeRequest(url, GET, nil, nil, headers, connectTimeOut, readWriteTimeout)
	return err
}

func (api *API) getTicket() error {
	//	https://localhost:4243/qps/ticket
	xrfKey := makeXrfKey()
	url := fmt.Sprintf("%s://%s:%v/qps/ticket?xrfkey=%s", "https", api.Server, api.AuthPort, xrfKey)
	fmt.Printf("url:%v\n", url)
	headers := make(map[string]string)
	headers[xrf_header] = xrfKey
	headers[qlik_user_header] = api.makeQlikUserHeader()
	headers[content_type_header] = application_json_content_type
	body := `{"UserDirectory":"` + "WIN8-VBOX" + `", "UserId":"` + "atscale" + `"}`
	err := api.makeRequest(url, POST, []byte(body), nil, headers, connectTimeOut, readWriteTimeout)
	return err
}

//https://help.qlik.com/en-US/sense-developer/2.1/Subsystems/EngineAPI/Content/CreatingAppLoadingData/CreateApps/create-app.htm
func (api *API) Create(name, localizedScriptMainSection string) (Response, error) {
	command := CreateApp(name)
	return api.executeWebsocketCommand(command)
}

//https://help.qlik.com/en-US/sense-developer/2.1/Subsystems/EngineAPI/Content/CreatingAppLoadingData/CreateApps/open-app.htm
func (api *API) Open(name, directory, user string) (Response, error) {
	command := OpenDoc(name, directory, user)
	return api.executeWebsocketCommand(command)
}

//https://help.qlik.com/en-US/sense-developer/2.1/Subsystems/EngineAPI/Content/CreatingAppLoadingData/CreateApps/create-and-open-app.htm
func (api *API) GetActiveDoc() (Response, error) {
	command := GetActiveDoc()
	return api.executeWebsocketCommand(command)
}

//https://help.qlik.com/en-US/sense-developer/2.1/Subsystems/EngineAPI/Content/CreatingAppLoadingData/EditDataLoadScript/set-get-script.htm
func (api *API) SetScript(handle int, script string) (Response, error) {
	command := SetScript(handle, script)
	return api.executeWebsocketCommand(command)
}

func (api *API) GetScript(handle int) (Response, error) {
	command := GetScript(handle)
	return api.executeWebsocketCommand(command)
}

func (api *API) CreateSheet(handle int, title, description, thumbnail, id string, rows, columns, rank int) (Response, error) {
	params := CreateSheetParams(title, description, thumbnail, id, rows, columns, rank)
	command := CreateSheet(handle, params)
	fmt.Printf("createsheet:%v\n", command.Json())
	return api.executeWebsocketCommand(command)
}

func (api *API) ListStreams() (Response, error) {
	command := GetStreamList()
	return api.executeWebsocketCommand(command)
}

func (api *API) DoReload(handle int) (Response, error) {
	command := DoReload(handle)
	return api.executeWebsocketCommand(command)
}

func (api *API) OpenWebSocket() error {
	ws := fmt.Sprintf("wss://%s:%v/app", api.Server, api.WebsocketPort)
	u, err := url.Parse(ws)
	if err != nil {
		return err
	}
	dialer := net.Dialer{}
	rawConn, err := tls.DialWithDialer(&dialer, "tcp", u.Host, api.getTlsConfig(api.ClientCert, api.ClientKey, api.CertAuth))

	if err != nil {
		return err
	}

	xrfKey := makeXrfKey()
	// your milage may differ
	wsHeaders := http.Header{
		"Origin":                   {"http://192.168.5.131"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
		xrf_header:                 {xrfKey},
		qlik_user_header:           {api.makeQlikUserHeader()},
		content_type_header:        {application_json_content_type},
	}
	websocketConnection, resp, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
	api.WebsocketConnection = websocketConnection
	if err != nil {
		return fmt.Errorf("websocket.NewClient Error: %s\nResp:%+v", err, resp)
	}
	return nil
}

func (api *API) CloseWebSocket() error {
	return api.WebsocketConnection.Close()
}

func (api *API) executeWebsocketCommand(command interface{}) (Response, error) {
	response := Response{}
	err := api.WebsocketConnection.WriteJSON(command)
	if err != nil {
		return response, err
	}
	if false {
		_, res, err := api.WebsocketConnection.ReadMessage()
		if err != nil {
			return response, err
		}
		fmt.Printf("res:%v\n", string(res))
		err = json.Unmarshal(res, &response)
		if err != nil {
			return response, err
		}
		if response.Error != nil {
			return response, response.Error.GetError()
		}
	} else {
		err = api.WebsocketConnection.ReadJSON(&response)
		if err == nil {
			return response, err
		}
		if response.Error != nil {
			return response, response.Error.GetError()
		}
	}
	return response, nil

}

func (api *API) makeQlikUserHeader() string {
	return fmt.Sprintf(user_header_value, api.Directory, api.QlikUser)
}

func makeXrfKey() string {
	id := uuid.NewV4().String()
	idNoDashes := strings.Replace(id, "-", "", -1)
	return idNoDashes[0:16]
}

func (api *API) makeRequest(requestUrl string, method string, payload []byte, result interface{}, headers map[string]string,
	cTimeout time.Duration, rwTimeout time.Duration) error {
	if true {
		fmt.Printf("%s:%v\n", method, requestUrl)
		if payload != nil {
			fmt.Printf("%v\n", string(payload))
		}
	}
	client := httputil.NewTimeoutClient(cTimeout, rwTimeout, true)
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
		fmt.Printf("body:%s\n", body)
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

func (api *API) getTlsConfig(certLocation, keyLocation, caFile string) *tls.Config {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	if len(certLocation) > 0 && len(keyLocation) > 0 {
		// Load client cert if available
		cert, err := tls.LoadX509KeyPair(certLocation, keyLocation)
		if err == nil {
			if len(caFile) > 0 {
				caCertPool := x509.NewCertPool()
				caCert, err := ioutil.ReadFile(caFile)
				if err != nil {
					fmt.Printf("Error setting up caFile [%s]:%v\n", caFile, err)
				}
				caCertPool.AppendCertsFromPEM(caCert)
				tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true, RootCAs: caCertPool}
				tlsConfig.BuildNameToCertificate()
			} else {
				tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
			}
		}
	}
	return tlsConfig
}
