package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kardianos/service"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello astaxie!") // send data to client side
}

var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	http.HandleFunc("/", sayhelloName)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "GenesisEicarService",
		DisplayName: "Genesis Eicar Example Service",
		Description: "This is a cross platform test binary for multiple operating systems.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
