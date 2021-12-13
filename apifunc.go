package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func ApiGET(url string) (string, int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", Localhost+url, nil)
	Check(err)

	req.Header.Set("Authorization", Auth)
	resp, err := client.Do(req)
	Check(err)

	defer resp.Body.Close()
	Check(err)

	returnstr, err := ioutil.ReadAll(resp.Body)
	Check(err)

	return string(returnstr), resp.StatusCode
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

func ApiDELETE(url string) int {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, Localhost+url, nil)
	Check(err)

	req.Header.Set("Authorization", Auth)
	resp, err := client.Do(req)
	Check(err)

	return resp.StatusCode
}
