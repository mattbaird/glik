package main

import (
	"fmt"
	"github.com/mattbaird/glik"
)

func main() {
	api := glik.DefaultApi()
	out, err := api.About()
	fmt.Printf("out:%v\n", out)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	listOfApps, err := api.List()
	fmt.Printf("list:%v\n", listOfApps)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}
	err = api.OpenWebSocket()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	defer api.CloseWebSocket()
	response, err := api.Create("scriptsetting", "main")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("Create:%v\n", response.Json())
	response, err = api.GetActiveDoc()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("GetActiveDoc:%v\n", response.Json())
	response, err = api.SetScript("Load RecNo() as NewNumbers AutoGenerate 10;")
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("SetScript:%v\n", response.Json())
	response, err = api.GetScript()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	fmt.Printf("GetScript:%v\n", response.Json())

	if false {
		copyResults, err := api.Copy("acdb78ec-a0ee-49f1-8741-3580e6af7f63", "testing published")
		if err != nil {
			fmt.Printf("error copying:%v\n", err)
			return
		}
		fmt.Printf("copied:%v\n", copyResults.Id)
		publishResults, err := api.Publish(copyResults.Id, "aaec8d41-5201-43ab-809f-3063750dfafd", "testing")
		fmt.Printf("publish:%v\n", publishResults.Published)
		if err != nil {
			fmt.Printf("err:%v\n", err)
		}
	}
}
