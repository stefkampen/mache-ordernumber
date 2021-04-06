package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

func main() {
	validate := func(input string) error {
		_, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label: "VIN",
	}

	vin, _ := prompt.Run()

	prompt2 := promptui.Prompt{
		Label:    "Start at",
		Validate: validate,
	}

	numberString, _ := prompt2.Run()
	number, _ := strconv.ParseInt(numberString, 10, 64)

	for number >= 0 {
		numberString := fmt.Sprintf("%04d", number)

		success := retrieveData("FORD", numberString, "https://shop.ford.com/aemservices/shop/vot/api/customerorder/?orderNumber="+numberString+"&partAttributes=BP2.*&vin="+vin)

		if success {
			number -= 1
		} else {
			time.Sleep(time.Millisecond * 3000)
		}
	}
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	_, _ = buf.ReadBytes('\n')
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
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		_, _ = buf.ReadBytes('\n')
		os.Exit(0)
	}

	if resp.Status != "404 Not Found" {
		log.Println(number + " wrong status, retrying same code " + resp.Status)
		return false
	}

	log.Println(number + " not found")

	return true
}
