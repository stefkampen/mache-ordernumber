package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	number := 9999
	vin := "WF0xxxxxx"

	for number >= 0 {
		numberString := fmt.Sprintf("%04d", number)

		success := retrieveData("FORD", numberString, "https://shop.ford.com/aemservices/shop/vot/api/customerorder/?orderNumber="+numberString+"&partAttributes=BP2.*&vin="+vin)

		if success {
			number -= 1
		}
		time.Sleep(time.Second * 15)
	}
}

func retrieveData(name string, number string, url string) bool {
	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(number + " error, retrying")
		return false
	}

	if resp.Status == "200 OK" {
		log.Println(number + " found, exiting")
		os.Exit(0)
	}

	if resp.Status != "404 Not Found" {
		log.Println(number + " wrong status, retrying same code " + resp.Status)
		return false
	}

	log.Println(number + " not found")

	return true
}
