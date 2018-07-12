package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	rando := RandomInt(1, 2147483640)
	assert.NotZero(t, rando)
	//assert.EqualValues(t, rando, 0, "testing")
}

func TestGetAlphaNumericString(t *testing.T) {
	rando := GetAlphaNumericString(10)
	assert.NotEqual(t, "", rando, "should not be an empty string")
	assert.Len(t, rando, 10, "should be 10 chars long")
	//assert.Equal(t, rando, "", "testing")
}

func TestGetAlphaString(t *testing.T) {
	rando := GetAlphaString(5)
	assert.NotEqual(t, "", rando, "should not be an empty string")
	assert.Len(t, rando, 5, "should be 10 chars long")
	//assert.Equal(t, rando, "", "testing")
}

func TestGetAlphaNumericSpecialString(t *testing.T) {
	rando := GetAlphaNumericSpecialString(7)
	assert.NotEqual(t, "", rando, "should not be an empty string")
	assert.Len(t, rando, 7, "should be 10 chars long")
	//assert.Equal(t, rando, "", "testing")
}

func TestGetBool(t *testing.T) {
	rando := GetBool()
	assert.NotNil(t, rando)
	assert.Equal(t, false, rando, "passes half the time")
}
