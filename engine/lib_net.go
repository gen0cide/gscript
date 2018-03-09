package engine

// func DNSQuestion(target, request string) (string, error) {
// 	if request == "A" {
// 		var stringAnswerArray []string
// 		answerPTR, err := net.LookupIP(target)
// 		if err != nil {
// 			return "failed", err
// 		}
// 		for _, answrPTR := range answerPTR {
// 			stringAnswerArray = append(stringAnswerArray, answrPTR.String())
// 		}
// 		stringAnswer := strings.Join(stringAnswerArray, "/n")
// 		return stringAnswer, nil
// 	} else if request == "TXT" {
// 		answerTXT, err := net.LookupTXT(target)
// 		if err != nil {
// 			return "failed", err
// 		}
// 		stringAnswer := strings.Join(answerTXT, "/n")
// 		return stringAnswer, nil
// 	} else if request == "PTR" {
// 		answerA, err := net.LookupAddr(target)
// 		if err != nil {
// 			return "failed", err
// 		}
// 		stringAnswer := strings.Join(answerA, "/n")
// 		return stringAnswer, nil
// 	} else if request == "MX" {
// 		var stringAnswerArray []string
// 		answerMX, err := net.LookupMX(target)
// 		if err != nil {
// 			return "failed", err
// 		}
// 		for _, answrMX := range answerMX {
// 			stringAnswerArray = append(stringAnswerArray, answrMX.Host)
// 		}
// 		stringAnswer := strings.Join(stringAnswerArray, "/n")
// 		return stringAnswer, nil
// 	} else if request == "NS" {
// 		var stringAnswerArray []string
// 		answerNS, err := net.LookupNS(target)
// 		if err != nil {
// 			return "failed", err
// 		}
// 		for _, answrNS := range answerNS {
// 			stringAnswerArray = append(stringAnswerArray, answrNS.Host)
// 		}
// 		stringAnswer := strings.Join(stringAnswerArray, "/n")
// 		return stringAnswer, nil
// 	} else if request == "CNAME" {
// 		answerCNAME, err := net.LookupCNAME(target)
// 		if err != nil {
// 			return "failed", err
// 		}
// 		return string(answerCNAME), nil
// 	} else {
// 		answerA, err := net.LookupHost(target)
// 		if err != nil {
// 			return "failed", err
// 		}
// 		stringAnswer := strings.Join(answerA, "/n")
// 		return stringAnswer, nil
// 	}
// }

// func HTTPGetFile(url string) (int, []byte, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	respCode := resp.StatusCode
// 	pageData, err := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()
// 	return respCode, pageData, nil
// }

// func PostJSON(url string, jsonString []byte) (int, []byte, error) {
// 	// encode json to sanity check, then decode to ensure the transmition syntax is clean
// 	var jsonObj interface{}
// 	if err := json.Unmarshal(jsonString, &jsonObj); err != nil {
// 		return 0, nil, err
// 	}
// 	jsonStringCleaned, err := json.Marshal(jsonObj)
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	resp, err := http.Post(url, " application/json", bytes.NewReader(jsonStringCleaned))
// 	if err != nil {
// 		return 0, nil, err
// 	}
// 	respCode := resp.StatusCode
// 	pageData, err := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()
// 	return respCode, pageData, nil
// }

// func TCPRead(ip, port string) ([]byte, error) {
// 	host := ip + ":" + port
// 	conn, err := net.Dial("tcp", host)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()
// 	buffer := make([]byte, 1024)
// 	conn.Read(buffer)
// 	return buffer, nil
// }

// func TCPWrite(writeData []byte, ip, port string) ([]byte, error) {
// 	host := ip + ":" + port
// 	conn, err := net.Dial("tcp", host)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()
// 	buffer := make([]byte, 1024)
// 	conn.Read(buffer)
// 	conn.Write(writeData)
// 	buffer2 := make([]byte, 1024)
// 	conn.Read(buffer2)
// 	return buffer2, nil
// }

// func UDPWrite(writeData []byte, ip, port string) error {
// 	host := ip + ":" + port
// 	conn, err := net.Dial("tcp", host)
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()
// 	conn.Write(writeData)
// 	return nil
// }

// func (e *Engine) VMHTTPRequest(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMSSHCmd(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMGetTweet(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMServeFileOverHTTP(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMDetectSSLMITM(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMCanPing(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMCanSSHLogin(call otto.FunctionCall) otto.Value {
// 	e.Logger.WithField("trace", "true").Errorf("Function Not Implemented")
// 	return otto.FalseValue()
// }

// func (e *Engine) VMRetrieveFileFromURL(call otto.FunctionCall) otto.Value {
// 	readURL, err := call.ArgumentList[0].ToString()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	_, bytes, err := HTTPGetFile(readURL)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("HTTP Error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	vmResponse, err := e.VM.ToValue(bytes)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMPostJSON(call otto.FunctionCall) otto.Value {
// 	readURL, err := call.ArgumentList[0].ToString()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	readJSONString, err := call.ArgumentList[1].ToString()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	_, bytes, err := PostJSON(readURL, []byte(readJSONString))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("HTTP Error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	vmResponse, err := e.VM.ToValue(bytes)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMDNSQuery(call otto.FunctionCall) otto.Value {
// 	targetDomain := call.Argument(0)
// 	queryType := call.Argument(1)
// 	targetDomainAsString, err := targetDomain.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	queryTypeAsString, err := queryType.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	result, err := DNSQuestion(targetDomainAsString.(string), queryTypeAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("DNS Error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	vmResult, err := e.VM.ToValue(result)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResult
// }

// func (e *Engine) VMGetHostname(call otto.FunctionCall) otto.Value {
// 	hostString, err := GetHostname()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("OS error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	vmResponse, err := e.VM.ToValue(hostString)
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Return value casting error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	return vmResponse
// }

// func (e *Engine) VMHasPublicIP(call otto.FunctionCall) otto.Value {
// 	respCode, _, err := HTTPGetFile("http://icanhazip.com")
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
// 		return otto.FalseValue()
// 	} else if respCode == 200 {
// 		return otto.TrueValue()
// 	} else {
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMCanMakeTCPConn(call otto.FunctionCall) otto.Value {
// 	ip := call.Argument(0)
// 	ipString, err := ip.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	port := call.Argument(1)
// 	portString, err := port.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	tcpResponse, err := TCPRead(ipString.(string), portString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if tcpResponse != nil {
// 		return otto.TrueValue()
// 	}
// 	return otto.FalseValue()
// }

// func (e *Engine) VMExpectedDNS(call otto.FunctionCall) otto.Value {
// 	targetDomain := call.Argument(0)
// 	queryType := call.Argument(1)
// 	expectedResult := call.Argument(2)
// 	targetDomainAsString, err := targetDomain.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	queryTypeAsString, err := queryType.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	expectedResultAsString, err := expectedResult.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	result, err := DNSQuestion(targetDomainAsString.(string), queryTypeAsString.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if expectedResultAsString.(string) == string(result) {
// 		return otto.TrueValue()
// 	}
// 	return otto.FalseValue()
// }

// func (e *Engine) VMCanMakeHTTPConn(call otto.FunctionCall) otto.Value {
// 	url1 := call.Argument(0)
// 	url1String, err := url1.Export()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	respCode, _, err := HTTPGetFile(url1String.(string))
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Net error: %s", err.Error())
// 		return otto.FalseValue()
// 	} else if respCode >= 200 && respCode < 300 {
// 		return otto.TrueValue()
// 	} else {
// 		return otto.FalseValue()
// 	}
// }

// func (e *Engine) VMTCPPortInUse(call otto.FunctionCall) otto.Value {
// 	port, err := call.Argument(0).ToInteger()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if port < 0 || port > 65535 {
// 		return otto.FalseValue()
// 	}
// 	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
// 	if err != nil {
// 		return otto.TrueValue()
// 	}
// 	conn.Close()
// 	return otto.FalseValue()
// }

// func (e *Engine) VMUDPPortInUse(call otto.FunctionCall) otto.Value {
// 	port, err := call.Argument(0).ToInteger()
// 	if err != nil {
// 		e.Logger.WithField("trace", "true").Errorf("Parameter parsing error: %s", err.Error())
// 		return otto.FalseValue()
// 	}
// 	if port < 0 || port > 65535 {
// 		return otto.FalseValue()
// 	}
// 	conn, err := net.Listen("udp", fmt.Sprintf("127.0.0.1:%d", port))
// 	if err != nil {
// 		return otto.TrueValue()
// 	}
// 	conn.Close()
// 	return otto.FalseValue()
// }
