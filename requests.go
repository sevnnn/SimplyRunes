package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func ApiGET(url string) ([]byte, int) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", Localhost+url, nil)
	Check(err)

	request.Header.Set("Authorization", Auth)
	response, err := client.Do(request)
	Check(err)

	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	Check(err)

	return read, response.StatusCode
}

func ApiDELETE(url string) int {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", Localhost+url, nil)
	Check(err)

	req.Header.Set("Authorization", Auth)
	resp, err := client.Do(req)
	Check(err)

	return resp.StatusCode
}

func ApiPOST(url string, payload []byte) int {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, Localhost+url, bytes.NewBuffer(payload))
	Check(err)

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", Auth)
	resp, err := client.Do(req)
	Check(err)

	return resp.StatusCode
}

func Get(url string) ([]byte, int) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	Check(err)

	response, err := client.Do(request)
	Check(err)

	defer response.Body.Close()
	read, err := ioutil.ReadAll(response.Body)
	Check(err)

	return read, response.StatusCode
}
