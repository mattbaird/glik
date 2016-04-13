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
	"encoding/json"
)

type Request struct {
	JsonRPCVersion string   `json:"jsonrpc,omitempty"`
	Id             int      `json:"id,omitempty"`
	Method         string   `json:"method,omitempty"`
	Handle         int      `json:"handle,omitempty"`
	Delta          bool     `json:"delta,omitempty"`
	Params         []string `json:"params,omitempty"`
}

func (r *Request) Json() string {
	result, _ := json.Marshal(r)
	return string(result)
}

func CreateApp(name string) Request {
	return NewRequest(0, "CreateApp", -1, []string{name})
}

func SetScript(script string) Request {
	return NewRequest(3, "SetScript", 1, []string{script})
}

func GetScript() Request {
	return NewRequest(2, "GetScript", 1, []string{})
}

func GetActiveDoc() Request {
	return NewRequest(1, "GetActiveDoc", -1, []string{})
}

func NewRequest(id int, method string, handle int, params []string) Request {
	return Request{JsonRPCVersion: "2.0", Id: id, Method: method, Handle: handle, Params: params}
}

type Params struct {
	QID                        string `json:"qId,omitempty"`
	AppName                    string `json:"qAppName,omitempty"`
	LocalizedScriptMainSection string `json:"qLocalizedScriptMainSection,omitempty"`
}

type CreateAppParms struct {
	Params
}

type Response struct {
	JsonRPCVersion string  `json:"jsonrpc,omitempty"`
	Id             int     `json:"id,omitempty"`
	Result         *Result `json:"result,omitempty"`
}

func (r *Response) Json() string {
	result, _ := json.Marshal(r)
	return string(result)
}

type Result struct {
	Success bool   `json:"qSuccess,omitempty"`
	AppId   string `json:"qAppId,omitempty"`
}
