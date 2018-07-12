package time

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnix(t *testing.T) {
	unixTime := GetUnix()
	assert.NotNil(t, unixTime)
	assert.Equal(t, unixTime, int(1531208097), "will likyl fail")
}
