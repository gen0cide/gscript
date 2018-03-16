package engine

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// DNSQuestion - Issues a DNS query and returns it's response
//
// Package
//
// net
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  DNSQuestion(target, request)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * target (string)
//  * request (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.answer (string)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = DNSQuestion(target, request);
//  // obj.answer
//  // obj.runtimeError
//
func (e *Engine) DNSQuestion(target, request string) (string, error) {
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

// HTTPGetFile - Retrives a file from an HTTP(s) endpoint
//
// Package
//
// net
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  HTTPGetFile(url)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * url (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.statusCode (int)
//  * obj.file ([]byte)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = HTTPGetFile(url);
//  // obj.statusCode
//  // obj.file
//  // obj.runtimeError
//
func (e *Engine) HTTPGetFile(url string) (int, []byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	pageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, pageData, nil
}

// PostJSON - Transmits a JSON object to a URL and retruns the HTTP status code and response
//
// Package
//
// net
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  PostJSON(url, json)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * url (string)
//  * json (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.statusCode (int)
//  * obj.response ([]byte)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = PostJSON(url, json);
//  // obj.statusCode
//  // obj.response
//  // obj.runtimeError
//
func (e *Engine) PostJSON(url string, jsonString string) (int, []byte, error) {
	// encode json to sanity check, then decode to ensure the transmition syntax is clean
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonString), &jsonObj); err != nil {
		return 0, nil, err
	}
	jsonStringCleaned, err := json.Marshal(jsonObj)
	if err != nil {
		return 0, nil, err
	}
	resp, err := http.Post(url, " application/json", bytes.NewReader(jsonStringCleaned))
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	pageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	resp.Body.Close()
	return resp.StatusCode, pageData, nil
}

// SSHCmd - Runs a command on a target host via SSH and returns the result (stdOut only). Uses both password and key authentication options
//
// Package
//
// net
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  SSHCmd(hostAndPort, cmd, username, password, key)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * hostAndPort (string)
//  * cmd (string)
//  * username (string)
//  * password (string)
//  * key ([]byte)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.response (string)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = SSHCmd(hostAndPort, cmd, username, password, key);
//  // obj.response
//  // obj.runtimeError
//
func (e *Engine) SSHCmd(hostAndPort, cmd, username, password string, key []byte) (string, error) {
	// create config
	sshConfig := &ssh.ClientConfig{
		User:            strings.TrimSpace(username),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if len(password) > 0 {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(password))
	}
	if len(key) > 0 {
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return "", err
		}
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(signer))
	}

	// connect
	client, err := ssh.Dial("tcp", hostAndPort, sshConfig)
	if err != nil {
		return "", err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	// run cmd
	responseBuffer := bytes.NewBuffer([]byte{})
	session.Stdout = responseBuffer
	if err = session.Run(cmd); err != nil {
		return "", err
	}
	return responseBuffer.String(), nil
}

// ServePathOverHTTPS - Starts an HTTPS webserver on a given port (default 443) for $X (default 30) number of seconds that acts as a file server rooted in a given path
//
// Package
//
// net
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  ServePathOverHTTPS(port, path, timeout)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * port (string)
//  * path (string)
//  * timeout (int64)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = ServePathOverHTTPS(port, path, timeout);
//  // obj.runtimeError
//
func (e *Engine) ServePathOverHTTPS(port, path string, timeout int64) error {
	// init
	if len(port) < 1 {
		port = "443"
	}
	if timeout == int64(0) {
		timeout = int64(30)
	}

	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// make cert
	template := x509.Certificate{
		IsCA: true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"Wonkaville"},
			Organization: []string{"Wonka Co"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(48 * time.Hour),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	template.DNSNames = append(template.DNSNames, hostname)
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			template.IPAddresses = append(template.IPAddresses, net.ParseIP(addr.String()))
		}
	}
	template.KeyUsage |= x509.KeyUsageCertSign
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return errors.New("Failed to create certificate: " + err.Error())
	}

	// make web server obj
	config := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{certData},
			PrivateKey:  privateKey,
		}},
	}
	srv := &http.Server{
		Addr:      "0.0.0.0:" + port,
		Handler:   http.FileServer(http.Dir(path)),
		TLSConfig: config,
	}

	// make new thread that waits for timeout to expire, then kills the server
	killSwitch := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		srv.Shutdown(context.Background())
		close(killSwitch)
	}()

	// start server
	return srv.ListenAndServe()
	<-killSwitch
	return nil
}

// IsTCPPortInUse - States whether or not a given TCP port is avalible for use
//
// Package
//
// net
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  IsTCPPortInUse(port)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * port (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.state (bool)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = IsTCPPortInUse(port);
//  // obj.state
//
func (e *Engine) IsTCPPortInUse(port string) bool {
	conn, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return true
	}
	defer conn.Close()
	return false
}

// IsUDPPortInUse - States whether or not a given UDP port is avalible for use
//
// Package
//
// net
//
// Author
//
// - Vyrus (https://github.com/vyrus001)
//
// Javascript
//
// Here is the Javascript method signature:
//  IsUDPPortInUse(port)
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//  * port (string)
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.state (bool)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = IsUDPPortInUse(port);
//  // obj.state
//
func (e *Engine) IsUDPPortInUse(port string) bool {
	conn, err := net.Listen("udp", "0.0.0.0:"+port)
	if err != nil {
		return true
	}
	defer conn.Close()
	return false
}

// GetLocalIPs - Gets an array of Ip addresses for the host
//
// Package
//
// net
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  GetLocalIPs()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.addresses ([]string)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = GetLocalIPs();
//  // obj.addresses
//
func (e *Engine) GetLocalIPs() []string {
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

// GetMACAddress - Gets the MAC address of the interface with an IPv4 address
//
// Package
//
// net
//
// Author
//
// - ahhh (https://github.com/ahhh)
//
// Javascript
//
// Here is the Javascript method signature:
//  GetMACAddress()
//
// Arguments
//
// Here is a list of the arguments for the Javascript function:
//
// Returns
//
// Here is a list of fields in the return object:
//  * obj.address (string)
//  * obj.runtimeError (error)
//
// Example
//
// Here is an example of how to use this function in gscript:
//  var obj = GetMACAddress();
//  // obj.address
//  // obj.runtimeError
//
func (e *Engine) GetMACAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var currentIP, currentNetworkHardwareName string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
			}
		}
	}
	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {
		if addrs, err := interf.Addrs(); err == nil {
			for _, addr := range addrs {
				if strings.Contains(addr.String(), currentIP) {
					currentNetworkHardwareName = interf.Name
				}
			}
		}
	}
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)
	if err != nil {
		return "", err
	}
	macAddress := netInterface.HardwareAddr
	return macAddress.String(), nil
}
