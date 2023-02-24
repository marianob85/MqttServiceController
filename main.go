package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kardianos/service"
)

var logger service.Logger

type program struct {
	closed chan struct{}
	client MqttClient
	wg     sync.WaitGroup
}

func main() {
	svcConfig := &service.Config{
		Name:        "OscamControl",
		DisplayName: "OscamControl_service",
		Description: "OscamControl service.",
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

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	var configuration = readConfiguration("config.json")
	p.client = NewMqttClient(configuration)
	p.closed = make(chan struct{})

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go p.oscamServiceHandler()
	go p.osServiceHandler()

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		close(p.closed)
	case <-p.closed:
		fmt.Printf("Got closed signal. Aborting...\n")
	}
	p.wg.Wait()
	p.client.client.Disconnect(100)
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	close(p.closed)
	p.wg.Wait()
	return nil
}

func (p *program) osServiceHandler() {
	p.wg.Add(1)
	defer p.wg.Done()

	p.client.publish("os", "active", true)
	p.client.subscribe("os/command", func(client mqtt.Client, message mqtt.Message) { osControl(message.Payload()) })

	for {
		time.Sleep(1 * time.Second)

		select {
		case <-p.closed:
			p.client.publish("os", "inactive", true)
			return
		default:

		}
	}
}

func (p *program) oscamServiceHandler() {
	p.wg.Add(1)
	defer p.wg.Done()

	oscamService := Service{"oscam"}
	p.client.subscribe("oscam/command", func(client mqtt.Client, message mqtt.Message) { oscamService.setStatePayload(message.Payload()) })

	var status, _, statusText = oscamService.checkStatus()
	p.client.publish("oscam", statusText, true)

	for {
		time.Sleep(1 * time.Second)

		select {
		case <-p.closed:
			return
		default:
			var newStatus, _, newStatusText = oscamService.checkStatus()
			if newStatus != status {
				status = newStatus
				p.client.publish("oscam", newStatusText, true)
			}
		}
	}
}
