package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVMDNSQuery(t *testing.T) {
	testScript := fmt.Sprintf(`
    var url = "google.com";
    var ip = "8.8.8.8";
    var type1 = "A";
    var return_value1 = DNSQuery(url, type1);
    var type2 = "CNAME";
    var return_value2 = DNSQuery(url, type2);
    var type3 = "TXT";
    var return_value3 = DNSQuery(url, type3);
    var type4 = "MX";
    var return_value4 = DNSQuery(url, type4);
    var type5 = "NS";
    var return_value5 = DNSQuery(url, type5);
    var type6 = "PTR";
    var return_value6 = DNSQuery(ip, type6);
  `)

	e := New("DNSQuery")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value1")
	assert.Nil(t, err)
	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
	retVal3, err := e.VM.Get("return_value3")
	assert.Nil(t, err)
	retValAsString3, err := retVal3.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString3)
	retVal4, err := e.VM.Get("return_value4")
	assert.Nil(t, err)
	retValAsString4, err := retVal4.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString4)
	retVal5, err := e.VM.Get("return_value5")
	assert.Nil(t, err)
	retValAsString5, err := retVal5.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString5)
	retVal6, err := e.VM.Get("return_value6")
	assert.Nil(t, err)
	retValAsString6, err := retVal6.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString6)
}

func TestVMCanMakeHTTPConn(t *testing.T) {
	testScript := `
    var url1 = "https://www.google.com";
    var return_value = CanMakeHTTPConn(url1);
  `

	e := New("HTTPTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMCanMakeTCPConn(t *testing.T) {
	testScript := `
    var ip = "towel.blinkenlights.nl";
    var port = "23";
    var return_value = CanMakeTCPConn(ip,port);
  `
	e := New("TCPTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMHasPublicIP(t *testing.T) {
	testScript := `
    var return_value = HasPublicIP();
  `
	e := New("PublicIPTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "true", retValAsString)
}

func TestVMTCPPortInUse(t *testing.T) {
	testScript := `
    var port = 8080;
    var return_value = TCPPortInUse(port);
  `
	e := New("TCPPortInUseTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal, err := e.VM.Get("return_value")
	assert.Nil(t, err)

	retValAsString, err := retVal.ToString()
	assert.Nil(t, err)

	assert.Equal(t, "false", retValAsString)
}

func TestVMExpectedDNS(t *testing.T) {
	testScript := fmt.Sprintf(`
    var url = "google.com";
    var type2 = "CNAME";
    var return_value2 = ExpectedDNS(url, type2, "google.com.");
  `)

	e := New("DNSQueryTest")
	e.CreateVM()

	e.VM.Run(testScript)
	retVal2, err := e.VM.Get("return_value2")
	assert.Nil(t, err)
	retValAsString2, err := retVal2.ToString()
	assert.Nil(t, err)
	assert.Equal(t, "true", retValAsString2)
}
