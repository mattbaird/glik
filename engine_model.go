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
	"fmt"
)

type Request struct {
	JsonRPCVersion string      `json:"jsonrpc,omitempty"`
	Id             int         `json:"id,omitempty"`
	Method         string      `json:"method,omitempty"`
	Handle         int         `json:"handle,omitempty"`
	Delta          bool        `json:"delta,omitempty"`
	Params         interface{} `json:"params"`
}

func (r *Request) Json() string {
	result, _ := json.Marshal(r)
	return string(result)
}

func CreateApp(name string) Request {
	return NewRequest(0, "CreateApp", -1, []string{name})
}

func OpenDoc(name string, user, directory string) Request {
	return NewRequest(0, "OpenDoc", -1, []string{name, fmt.Sprintf("UserDirectory=%s; UserId=%s", directory, user)})
}

func CreateAppEx(name string) Request {
	return NewRequest(0, "CreateDocEx", -1, []string{name})
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

func NewRequest(id int, method string, handle int, params interface{}) Request {
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
	JsonRPCVersion string          `json:"jsonrpc,omitempty"`
	Id             int             `json:"id,omitempty"`
	Result         *Result         `json:"result,omitempty"`
	Error          *WebsocketError `json:"error,omitempty"`
}

func (r *Response) Json() string {
	result, _ := json.Marshal(r)
	return string(result)
}

type WebsocketError struct {
	Code      int    `json:"code,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Message   string `json:"message,omitempty"`
}

type Result struct {
	Success bool   `json:"qSuccess,omitempty"`
	AppId   string `json:"qAppId,omitempty"`
	Type    string `json:"qType,omitempty"`
	Handle  int    `json:"qHandle,omitempty"`
	Script  string `json:"qScript,omitempty"`
}

type SheetParams struct {
	Title        string        `json:"title,omitempty"`
	Description  string        `json:"description,omitempty"`
	Info         *Info         `json:"qInfo,omitempty"`
	ChildListDef *ChildListDef `json:"qChildListDef,omitempty"`
}

type Info struct {
	ID   string `json:"qId,omitempty"`
	Type string `json:"qType,omitempty"`
}

type ChildListDef struct {
	Data *Data `json:"qData,omitempty"`
}

type Data struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Meta        string `json:"meta,omitempty"`
	Order       string `json:"order,omitempty"`
	Type        string `json:"type,omitempty"`
	Id          string `json:"id,omitempty"`
	Lb          string `json:"lb,omitempty"`
	Hc          string `json:"hc,omitempty"`
}

func CreateSheetParams(title, description, id string) SheetParams {
	params := SheetParams{Title: title, Description: description}
	params.Info = &Info{ID: id, Type: "sheet"}
	params.ChildListDef = &ChildListDef{}
	params.ChildListDef.Data = &Data{Title: "/title",
		Description: "/description",
		Meta:        "/meta",
		Order:       "/order",
		Type:        "/qInfo/qType",
		Id:          "/qInfo/qId",
		Lb:          "/qListObjectDef",
		Hc:          "/qHyperCubeDef"}
	return params
}

func CreateSheet(params SheetParams) Request {
	return NewRequest(1, "CreateObject", 1, []SheetParams{params})
}
