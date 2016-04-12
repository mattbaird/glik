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

const DEFAULT_SERVER = "192.168.99.5"
const CRLF = "\r\n"
const DEFAULT_USER = "sa_repository"
const DEFAULT_DIR = "Internal"
const DEFAULT_QRS_PORT = 4242
const DEFAULT_AUTH_PORT = 4243
const DEFAULT_WEBSOCKET_PORT = 4747

type API struct {
	Server        string
	QrsPort       int
	AuthPort      int
	WebsocketPort int
	Version       string
	Directory     string
	QlikUser      string
	ClientKey     string
	ClientCert    string
	XrfKey        string
	CertAuth      string
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
	api := NewAPI(DEFAULT_SERVER, DEFAULT_DIR, DEFAULT_USER, DEFAULT_QRS_PORT, DEFAULT_AUTH_PORT, DEFAULT_WEBSOCKET_PORT)
	api.ClientKey = string(clientKeyBytes)
	api.ClientCert = string(clientCertBytes)
	api.CertAuth = string(certAuthBytes)
	return api
}

func NewAPI(server string, directory, user string, qrsPort, authPort, websocketPort int) API {
	fixedUpServer := server
	if strings.HasSuffix(server, "/") {
		fixedUpServer = server[0 : len(server)-1]
	}
	api := API{Server: fixedUpServer}
	api.QlikUser = user
	api.Directory = directory
	api.QrsPort = qrsPort
	api.AuthPort = authPort
	api.WebsocketPort = websocketPort
	return api
}

type About struct {
	BuildVersion     string `json:"buildVersion,omitempty"`
	BuildDate        string `json:"buildDate,omitempty"`
	DatabaseProvider string `json:"databaseProvider,omitempty"`
	NodeType         int    `json:"nodeType,omitempty"`
	SchemaPath       string `json:"schemaPath,omitempty"`
}

type ApplicationResult struct {
	Id                    string           `json:"id,omitempty"`
	CreatedDate           string           `json:"createdDate,omitempty"`
	ModifiedDate          string           `json:"modifiedDate,omitempty"`
	ModifiedByUserName    string           `json:"modifiedByUserName,omitempty"`
	CustomProperties      []CustomProperty `json:"customProperites,omitempty"`
	Owner                 *Owner           `json:"owner,omitempty"`
	Name                  string           `json:"name,omitempty"`
	AppId                 string           `json:"appId,omitempty"`
	PublishTime           string           `json:"publishTime,omitempty"`
	Published             bool             `json:"published,omitempty"`
	Tags                  []string         `json:"tags,omitempty"`
	Description           string           `json:"description,omitempty"`
	Stream                *Stream          `json:"stream,omitempty"`
	FileSize              int              `json:"fileSize,omitempty"`
	LastReloadTime        string           `json:"lastReloadTime,omitempty"`
	Thumbnail             string           `json:"thumbnail,omitempty"`
	SavedInProductVersion string           `json:"savedInProductVersion,omitempty"`
	MigrationHash         string           `json:"migrationHash,omitempty"`
	Privileges            *Privileges      `json:"privileges,omitempty"`
	SchemaPath            string           `json:"schemaPath,omitempty"`
}

type CustomProperty struct {
}

type Owner struct {
	UserId        string `json:"userId,omitempty"`
	UserDirectory string `json:"userDirectory,omitempty"`
	Name          string `json:"name,omitempty"`
	Id            string `json:"id,omitempty"`
}

type Stream struct {
	Name       string      `json:"name,omitempty"`
	Id         string      `json:"id,omitempty"`
	Privileges *Privileges `json:"privileges,omitempty"`
}

type Privileges struct {
}
