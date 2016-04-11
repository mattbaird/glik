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
