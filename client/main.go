package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var c http.Client

func main() {

	c = http.Client{Timeout: time.Duration(10) * time.Second}

	Getall()
	//Post()
	//Getpost()
	//Put()
	//Del()

}

func Getall() {

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body : %s", body)
}
func Post() {

	myJson := bytes.NewBuffer([]byte(`{"id":"13","name":"anjali","bal":"232"}`))
	resp, err := c.Post("http://localhost:8000/insert", "application/json", myJson)
	if err != nil {
		fmt.Errorf("Error %s", err)
		return
	}

	defer resp.Body.Close()
	body, err :=ioutil.ReadAll(resp.Body)
	fmt.Printf("Body : %s\n",body)
	
	
}
 
func Getpost(){

	myJson1 := bytes.NewBuffer([]byte(`{"id":"13"}`))

	req, err := http.NewRequest("GET", "http://localhost:8000/show_selc", myJson1)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body : %s", body)
}

func Put(){
	myJson2 := bytes.NewBuffer([]byte(`{"id":"13","bal":"100"}`))

	req, err := http.NewRequest("POST", "http://localhost:8000/update", myJson2)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body : %s", body)


}

func Del(){
	myJson3 := bytes.NewBuffer([]byte(`{"id":"10"}`))

	req, err := http.NewRequest("POST", "http://localhost:8000/delete", myJson3)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Body : %s", body)



}
