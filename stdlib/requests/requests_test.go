package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testHeaders = map[string]interface{}{
		"X-Test-ID": "gscript",
	}
)

func TestGetURLAsString(t *testing.T) {
	obj, resp, err := GetURLAsString("http://icanhazip.com/", testHeaders, false)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)

	obj, resp, err = GetURLAsString("https://google.com/", testHeaders, false)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)

	obj, resp, err = GetURLAsString("https://wrong.host.badssl.com/", testHeaders, true)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)

	obj, resp, err = GetURLAsString("https://wrong.host.badssl.com/", testHeaders, false)
	assert.NotNil(t, err)
	assert.Nil(t, obj)
	assert.Equal(t, "", resp, "blank means received no data, due to error")
}

func TestGetURLAsBytes(t *testing.T) {
	obj, resp, err := GetURLAsBytes("http://icanhazip.com/", testHeaders, false)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
}

func TestPostJSON(t *testing.T) {
	json := map[string]interface{}{
		"menu": "item",
	}
	obj, resp, err := PostJSON("http://postb.in/5d83jRPR", json, testHeaders, false)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
}

func TestPostURL(t *testing.T) {
	obj, resp, err := PostURL("http://postb.in/dEH0maRf", "first test data", testHeaders, false)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
}

//func TestPostBinary(t *testing.T) {
//	obj, resp, err := PostBinary("http://postb.in/L7daTOz8", "./example_test.txt", nil, true)
//	assert.Nil(t, err)
//	assert.NotNil(t, obj)
//	assert.NotNil(t, resp)
//	assert.Equal(t, resp, "blah blah", "should throw an error")
//}
