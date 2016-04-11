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
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

var (
	connectTimeOut   = time.Duration(30 * time.Second)
	readWriteTimeout = time.Duration(30 * time.Second)
)

const DEFAULT_SERVER = "https://192.168.99.5:4242"
const CRLF = "\r\n"

type API struct {
	Server     string
	Version    string
	XrfKey     string
	QlikUser   string
	ClientKey  string
	ClientCert string
	CertAuth   string
}

func DefaultApi() API {
	clientKeyBytes, err := ioutil.ReadFile("client_key.pem")
	if err != nil {
		fmt.Printf("error reading client key bytes:%v\n", err)
	}
	clientCertBytes, err := ioutil.ReadFile("client.pem")
	if err != nil {
		fmt.Printf("error reading client cert bytes:%v\n", err)
	}
	certAuthBytes, err := ioutil.ReadFile("client.pem")
	if err != nil {
		fmt.Printf("error reading ca bytes:%v\n", err)
	}
	api := NewAPI(DEFAULT_SERVER)
	api.ClientKey = string(clientKeyBytes)
	api.ClientCert = string(clientCertBytes)
	api.CertAuth = string(certAuthBytes)
	return api
}

func NewAPI(server string) API {
	fixedUpServer := server
	if strings.HasSuffix(server, "/") {
		fixedUpServer = server[0 : len(server)-1]
	}
	return API{Server: fixedUpServer}
}

type About struct {
	BuildVersion     string `json:"buildVersion,omitempty"`
	BuildDate        string `json:"buildDate,omitempty"`
	DatabaseProvider string `json:"databaseProvider,omitempty"`
	NodeType         int    `json:"nodeType,omitempty"`
	SchemaPath       string `json:"schemaPath,omitempty"`
}
