package requests

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// IMPORT EXCEPTION OVERRIDE (gen0cide):
	// This library is used for JSON serialization
	// and follows our rules of no dependencies outside
	// of the standard library. It's stable, widely used,
	// and is literally the exact thing I was going to
	// implement for serializing arbitrary JSON but
	// has the advantage of being more widely tested
	// and stable.
	"github.com/Jeffail/gabs"

	// IMPORT EXCEPTION OVERRIDE (gen0cide):
	// This is a small wrapper by the *AUTHORS OF GO ITSELF*
	// to provide the ability to stack and wrap
	// multiple error messages for better context.
	// It's stable, widely used, incredibly small,
	// and has zero dependencies except for a few
	// standard library packages.
	"github.com/pkg/errors"
)

//GetURLAsString will fetch a url and return an http response object, the body as a string, and an error
func GetURLAsString(url string, headers map[string]interface{}, ignoresslerrors bool) (*http.Response, string, error) {
	c := createClient(ignoresslerrors)
	req, err := createReq("GET", url, headers, nil)
	if err != nil {
		return nil, "", err
	}
	return do(c, req)
}

//GetURLAsBytes will fetch a url, headers, and a bool for ignoring ssl errors. this returns an http response object, the body as a string, and an error
func GetURLAsBytes(url string, headers map[string]interface{}, ignoresslerrors bool) (*http.Response, []byte, error) {
	c := createClient(ignoresslerrors)
	req, err := createReq("GET", url, headers, nil)
	if err != nil {
		return nil, []byte{}, err
	}
	resp, data, err := do(c, req)
	return resp, []byte(data), err
}

//PostJSON takes a url, json data, a map of headers, and a bool to ignore ssl errors, posts json data to url
func PostJSON(url string, data map[string]interface{}, headers map[string]interface{}, ignoresslerrors bool) (*http.Response, string, error) {
	jsonData, err := genJSON(data)
	if err != nil {
		return nil, "", err
	}
	c := createClient(ignoresslerrors)
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/json"
	}
	req, err := createReq("POST", url, headers, jsonData)
	if err != nil {
		return nil, "", err
	}
	return do(c, req)
}

//PostURL posts the specified data to a url endpoint as text/plain data
func PostURL(url string, data string, headers map[string]interface{}, ignoresslerrors bool) (*http.Response, string, error) {
	c := createClient(ignoresslerrors)
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "text/plain"
	}
	req, err := createReq("POST", url, headers, strings.NewReader(data))
	if err != nil {
		return nil, "", errors.Wrap(err, "error creating HTTP request")
	}
	return do(c, req)
}

//PostBinary posts the specified data to a url endpoint as application/octet-stream data
func PostBinary(url string, readPath string, headers map[string]interface{}, ignoresslerrors bool) (*http.Response, string, error) {
	absPath, err := filepath.Abs(readPath)
	if err != nil {
		return nil, "", errors.Wrap(err, "error reading body file")
	}
	fileReader, err := os.Open(absPath)
	if err != nil {
		return nil, "", errors.Wrap(err, "error reading body file")
	}
	c := createClient(ignoresslerrors)
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/octet-stream"
	}
	req, err := createReq("POST", url, headers, fileReader)
	if err != nil {
		return nil, "", errors.Wrap(err, "error creating HTTP request")
	}
	return do(c, req)
}

func do(c *http.Client, r *http.Request) (*http.Response, string, error) {
	resp, err := c.Do(r)
	if err != nil {
		return resp, "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, "", errors.Wrap(err, "could not read response body")
	}
	return resp, string(data), nil
}

func createClient(ignoreSSL bool) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreSSL,
			},
		},
	}
}

func genJSON(data map[string]interface{}) (*strings.Reader, error) {
	jsonObj := gabs.Wrap(data)
	if jsonObj == nil {
		return strings.NewReader("{}"), errors.New("could not generate JSON body")
	}
	return strings.NewReader(jsonObj.String()), nil
}

func createReq(method, url string, headers map[string]interface{}, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	return req, nil
}
