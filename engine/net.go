package engine

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func DNSQuestion(target, request string) (string, error) {
	if request == "A" {
		var stringAnswerArray []string
		answerPTR, err := net.LookupIP(target)
		if err != nil {
			return "failed", err
		}
		for _, answrPTR := range answerPTR {
			stringAnswerArray = append(stringAnswerArray, answrPTR.String())
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		return stringAnswer, nil
	} else if request == "TXT" {
		answerTXT, err := net.LookupTXT(target)
		if err != nil {
			return "failed", err
		}
		stringAnswer := strings.Join(answerTXT, "/n")
		return stringAnswer, nil
	} else if request == "PTR" {
		answerA, err := net.LookupAddr(target)
		if err != nil {
			return "failed", err
		}
		stringAnswer := strings.Join(answerA, "/n")
		return stringAnswer, nil
	} else if request == "MX" {
		var stringAnswerArray []string
		answerMX, err := net.LookupMX(target)
		if err != nil {
			return "failed", err
		}
		for _, answrMX := range answerMX {
			stringAnswerArray = append(stringAnswerArray, answrMX.Host)
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		return stringAnswer, nil
	} else if request == "NS" {
		var stringAnswerArray []string
		answerNS, err := net.LookupNS(target)
		if err != nil {
			return "failed", err
		}
		for _, answrNS := range answerNS {
			stringAnswerArray = append(stringAnswerArray, answrNS.Host)
		}
		stringAnswer := strings.Join(stringAnswerArray, "/n")
		return stringAnswer, nil
	} else if request == "CNAME" {
		answerCNAME, err := net.LookupCNAME(target)
		if err != nil {
			return "failed", err
		}
		return string(answerCNAME), nil
	} else {
		answerA, err := net.LookupHost(target)
		if err != nil {
			return "failed", err
		}
		stringAnswer := strings.Join(answerA, "/n")
		return stringAnswer, nil
	}
}

func HTTPGetFile(url string) (int, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}
	respCode := resp.StatusCode
	pageData, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return respCode, pageData, nil
}

func TCPRead(ip, port string) ([]byte, error) {
	host := ip + ":" + port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	return buffer, nil
}

func TCPWrite(writeData []byte, ip, port string) ([]byte, error) {
	host := ip + ":" + port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	conn.Read(buffer)
	conn.Write(writeData)
	buffer2 := make([]byte, 1024)
	conn.Read(buffer2)
	return buffer2, nil
}

func UDPWrite(writeData []byte, ip, port string) error {
	host := ip + ":" + port
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.Write(writeData)
	return nil
}

func GetLocalIPs() []string {
	addresses := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return addresses
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addresses = append(addresses, ipnet.IP.String())
			}
		}
	}
	return addresses
}
