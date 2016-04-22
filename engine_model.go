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

func SetScript(handle int, script string) Request {
	return NewRequest(3, "SetScript", handle, []string{script})
}

func GetScript(handle int) Request {
	return NewRequest(2, "GetScript", handle, []string{})
}

func GetStreamList() Request {
	return NewRequest(0, "GetStreamList", -1, []string{})
}

func DoReload(handle int) Request {
	return NewRequest(2, "DoReload", handle, []string{})
}

func GetActiveDoc() Request {
	return NewRequest(1, "GetActiveDoc", -1, []string{})
}

func NewRequest(id int, method string, handle int, params interface{}) Request {
	if params != nil {
		fmt.Printf("Params:%v\n", params)
	}
	return Request{JsonRPCVersion: "2.0", Id: id, Method: method, Handle: handle, Params: params}
}

func CreateSheet(handle int, params SheetParams) Request {
	return NewRequest(12, "CreateObject", handle, []SheetParams{params})
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
	Id             int             `json:"id"`
	Result         *Result         `json:"result,omitempty"`
	Error          *WebsocketError `json:"error,omitempty"`
	Change         []int           `json:"change"`
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
	Success    bool           `json:"qSuccess,omitempty"`
	AppId      string         `json:"qAppId,omitempty"`
	Type       string         `json:"qType,omitempty"`
	Handle     int            `json:"qHandle,omitempty"`
	Script     string         `json:"qScript,omitempty"`
	StreamList []EngineStream `json:"qStreamList,omitempty"`
	Return     Return         `json:"qReturn,omitempty"`
}

type Return struct {
	Type   string `json:"qType,omitempty"`
	Handle int    `json:"qHandle,omitempty"`
}

func (e *WebsocketError) GetError() error {
	return fmt.Errorf("Error [%v]: %s - %s", e.Code, e.Message, e.Parameter)
}

type EngineStream struct {
	Id   string `json:"qId,omitempty"`
	Name string `json:"qName,omitempty"`
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

type SheetParamsEx struct {
	Title        string        `json:"title,omitempty"`
	Description  string        `json:"description,omitempty"`
	Info         *Info         `json:"qInfo,omitempty"`
	ChildListDef *ChildListDef `json:"qChildListDef,omitempty"`
}

func CreateSheetParamsEx(title, description, id string) SheetParamsEx {
	params := SheetParamsEx{Title: title, Description: description}
	params.Info = &Info{ID: "SH01", Type: "sheet"}
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

type SheetParams struct {
	MetaDef   *MetaDef      `json:"qMetaDef,omitempty"`
	Rank      int           `json:"rank,omitempty"`
	Thumbnail string        `json:"thumbnail,omitempty"`
	Columns   int           `json:"columns,omitempty"`
	Rows      int           `json:"rows,omitempty"`
	Cells     []interface{} `json:"cells,omitempty"`
	Info      *Info         `json:"qInfo,omitempty"`
}

type MetaDef struct {
	Title        string        `json:"title,omitempty"`
	Description  string        `json:"description,omitempty"`
	ChildListDef *ChildListDef `json:"qChildListDef,omitempty"`
}

type Info struct {
	ID   string `json:"qId,omitempty"`
	Type string `json:"qType,omitempty"`
}

func CreateSheetParams(title, description, thumbnail, id string, rows, columns, rank int) SheetParams {
	params := SheetParams{Columns: columns, Rows: rows, Rank: rank, Thumbnail: thumbnail}
	params.MetaDef = &MetaDef{Title: title, Description: description}
	params.Info = &Info{ID: id, Type: "sheet"}
	return params
}
