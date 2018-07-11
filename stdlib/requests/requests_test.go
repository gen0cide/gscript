package requests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURLAsString(t *testing.T) {
	obj, resp, err := GetURLAsString("http://icanhazip.com/", nil, true)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
	//assert.Equal(t, resp, "aok", "this should error")

	obj, resp, err = GetURLAsString("https://google.com/", nil, false)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
	//assert.Equal(t, resp, "aok", "this should error")

	obj, resp, err = GetURLAsString("https://wrong.host.badssl.com/", nil, true)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
	//assert.Equal(t, resp, "aok", "this should error")

	obj, resp, err = GetURLAsString("https://wrong.host.badssl.com/", nil, false)
	assert.NotNil(t, err)
	assert.Nil(t, obj)
	assert.Equal(t, "", resp, "blank means received no data, due to error")
	//assert.Equal(t, resp, "aok", "this should error")
}

func TestGetURLAsBytes(t *testing.T) {
	obj, resp, err := GetURLAsBytes("http://icanhazip.com/", nil, true)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
	//assert.Equal(t, string(resp), "aok", "this should error")
}

func TestPostJSON(t *testing.T) {
	json := "{\"menu\":\"item\"}"
	obj, resp, err := PostJSON("http://postb.in/5d83jRPR", json, nil, true)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
}

func TestPostURL(t *testing.T) {
	obj, resp, err := PostURL("http://postb.in/dEH0maRf", "first test data", nil, true)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
	//assert.Equal(t, resp, "blah blah", "should throw an error")
}

func TestPostBinary(t *testing.T) {
	obj, resp, err := PostBinary("http://postb.in/L7daTOz8", "./example_test.txt", nil, true)
	assert.Nil(t, err)
	assert.NotNil(t, obj)
	assert.NotNil(t, resp)
	//assert.Equal(t, resp, "blah blah", "should throw an error")
}
