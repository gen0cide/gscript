package engine

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"math/big"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

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

func (e *Engine) PostJSON(url string, jsonString []byte) (int, []byte, error) {
	// encode json to sanity check, then decode to ensure the transmition syntax is clean
	var jsonObj interface{}
	if err := json.Unmarshal(jsonString, &jsonObj); err != nil {
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
	client, err := ssh.Dial("tcp", host, sshConfig)
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
	responseBuffer := bytes.NewBuffer([]byte)
	session.Stdout = &responseBuffer
	if err = session.Run(cmd); err != nil {
		return "", err
	}
	return responseBuffer.String(), nil
}

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
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return errors.New("failed to generate serial number: " + err.Error())
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Wonka Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(48 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	template.DNSNames = append(template.DNSNames, hostname)
	for _, iface := range net.Interfaces() {
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			template.IPAddresses = append(template.IPAddresses, addr.String())
		}
	}
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign
	certData, err := x509.CreateCertificate(rand.Reader, &template, &template, privateKey.PublicKey, privateKey)
	if err != nil {
		return errors.New("Failed to create certificate: " + err)
	}
	cert, err := tls.LoadX509KeyPair(certData, keyData)
	if err != nil {
		return err
	}

	// make web server obj
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	srv := &http.Server{
		Addr:      "0.0.0.0:" + port,
		Handler:   http.FileServer(http.Dir(path)),
		TLSConfig: config,
	}

	// make new thread that waits for timeout to expire, then kills the server
	killSwitch := make(chan struct{})
	go func() {
		time.Sleep(timeout * time.Second)
		srv.Shutdown(context.Background())
		close(killSwitch)
	}()

	// start server
	return srv.ListenAndServe()
	<-killSwitch
}

func (e *Engine) IsTCPPortInUse(port string) bool {
	conn, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (e *Engine) IsUDPPortInUse(port string) bool {
	conn, err := net.Listen("udp", "0.0.0.0:"+port)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
