package engine

import (
	"fmt"
	"net"

	"github.com/robertkrimen/otto"
)

func (e *Engine) VMHTTPRequest(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMSSHCmd(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMGetTweet(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMServeFileOverHTTP(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMDetectSSLMITM(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMCanPing(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMCanSSHLogin(call otto.FunctionCall) otto.Value {
	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
	return otto.FalseValue()
}

func (e *Engine) VMRetrieveFileFromURL(call otto.FunctionCall) otto.Value {
	readURL, err := call.ArgumentList[0].ToString()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	_, bytes, err := HTTPGetFile(readURL)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("HTTP Error: %s", err.Error())
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(bytes)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMDNSQuery(call otto.FunctionCall) otto.Value {
	targetDomain := call.Argument(0)
	queryType := call.Argument(1)
	targetDomainAsString, err := targetDomain.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	queryTypeAsString, err := queryType.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	result, err := DNSQuestion(targetDomainAsString.(string), queryTypeAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("DNS Error: %s", err.Error())
		return otto.FalseValue()
	}
	vmResult, err := e.VM.ToValue(result)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResult
}

func (e *Engine) VMGetHostname(call otto.FunctionCall) otto.Value {
	hostString, err := GetHostname()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("OS error: %s", err.Error())
		return otto.FalseValue()
	}
	vmResponse, err := e.VM.ToValue(hostString)
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
		return otto.FalseValue()
	}
	return vmResponse
}

func (e *Engine) VMHasPublicIP(call otto.FunctionCall) otto.Value {
	respCode, _, err := HTTPGetFile("http://icanhazip.com")
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
		return otto.FalseValue()
	} else if respCode == 200 {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMCanMakeTCPConn(call otto.FunctionCall) otto.Value {
	ip := call.Argument(0)
	ipString, err := ip.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	port := call.Argument(1)
	portString, err := port.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	tcpResponse, err := TCPRead(ipString.(string), portString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
		return otto.FalseValue()
	}
	if tcpResponse != nil {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMExpectedDNS(call otto.FunctionCall) otto.Value {
	targetDomain := call.Argument(0)
	queryType := call.Argument(1)
	expectedResult := call.Argument(2)
	targetDomainAsString, err := targetDomain.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	queryTypeAsString, err := queryType.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	expectedResultAsString, err := expectedResult.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	result, err := DNSQuestion(targetDomainAsString.(string), queryTypeAsString.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
		return otto.FalseValue()
	}
	if expectedResultAsString.(string) == string(result) {
		return otto.TrueValue()
	}
	return otto.FalseValue()
}

func (e *Engine) VMCanMakeHTTPConn(call otto.FunctionCall) otto.Value {
	url1 := call.Argument(0)
	url1String, err := url1.Export()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	respCode, _, err := HTTPGetFile(url1String.(string))
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
		return otto.FalseValue()
	} else if respCode >= 200 && respCode < 300 {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}

func (e *Engine) VMTCPPortInUse(call otto.FunctionCall) otto.Value {
	port, err := call.Argument(0).ToInteger()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	if port < 0 || port > 65535 {
		return otto.FalseValue()
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return otto.TrueValue()
	}
	conn.Close()
	return otto.FalseValue()
}

func (e *Engine) VMUDPPortInUse(call otto.FunctionCall) otto.Value {
	port, err := call.Argument(0).ToInteger()
	if err != nil {
		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
		return otto.FalseValue()
	}
	if port < 0 || port > 65535 {
		return otto.FalseValue()
	}
	conn, err := net.Listen("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return otto.TrueValue()
	}
	conn.Close()
	return otto.FalseValue()
}
