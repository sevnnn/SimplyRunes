package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Get(url string) ([]byte, int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	Check(err)

	resp, err := client.Do(req)
	Check(err)

	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	Check(err)

	return read, resp.StatusCode
}

func ApiGet(endpoint string) ([]byte, int) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", APILink+endpoint, nil)
	Check(err)

	request.Header.Set("Authorization", Auth)
	response, err := client.Do(request)
	Check(err)

	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	Check(err)

	return read, response.StatusCode
}

func ApiDELETE(endpoint string) int {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", APILink+endpoint, nil)
	Check(err)

	req.Header.Set("Authorization", Auth)
	resp, err := client.Do(req)
	Check(err)

	return resp.StatusCode
}

func ApiPOST(endpoint string, payload []byte) int {
	client := &http.Client{}
	req, err := http.NewRequest("POST", APILink+endpoint, bytes.NewBuffer(payload))
	Check(err)

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", Auth)
	resp, err := client.Do(req)
	Check(err)

	return resp.StatusCode
}
