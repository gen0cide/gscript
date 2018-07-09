package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

//GetURLAsString will fetch a url and return an http response object, the body as a string, and an error
func GetURLAsString(url string, headers map[string]string, ignoresslerrors bool) (*http.Response, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	if ignoresslerrors == true {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	pageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return resp, string(pageData), nil
}

//GetURLAsBytes will fetch a url, headers, and a bool for ignoring ssl errors. this returns an http response object, the body as a string, and an error
func GetURLAsBytes(url string, headers map[string]string, ignoresslerrors bool) (*http.Response, []byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	if ignoresslerrors == true {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	pageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, pageData, nil
}

//PostJSON takes a url, json data, a map of headers, and a bool to ignore ssl errors, posts json data to url
func PostJSON(url string, jsondata string, headers map[string]string, ignoresslerrors bool) (*http.Response, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	if ignoresslerrors == true {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}
	// encode json to sanity check, then decode to ensure the transmition syntax is clean
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsondata), &jsonObj); err != nil {
		return nil, "", err
	}
	jsonStringCleaned, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, "", err
	}
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonStringCleaned))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	pageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	resp.Body.Close()
	return resp, string(pageData), nil
}

//PostURL posts the specified data to a url endpoint as text/plain data
func PostURL(url string, data string, headers map[string]string, ignoresslerrors bool) (*http.Response, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	if ignoresslerrors == true {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "text/plain", bytes.NewBufferString(data))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	pageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	resp.Body.Close()
	return resp, string(pageData), nil
}

//PostBinary posts the specified data to a url endpoint as application/octet-stream data
func PostBinary(url string, readPath string, headers map[string]string, ignoresslerrors bool) (*http.Response, string, error) {
	absPath, err := filepath.Abs(readPath)
	if err != nil {
		return nil, "", err
	}
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, "", err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	if ignoresslerrors == true {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "application/octet-stream", bytes.NewReader([]byte(contents)))
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	pageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	resp.Body.Close()
	return resp, string(pageData), nil
}
