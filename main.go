package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	var configuration = readConfiguration("config.json")
	client := NewMqttClient(configuration)
	closed := make(chan struct{})
	wait := &sync.WaitGroup{}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go oscamServiceHandler(client, wait, closed)
	go osServiceHandler(client, wait, closed)

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		close(closed)
	}
	wait.Wait()
	client.client.Disconnect(100)
}

func osServiceHandler(client MqttClient, wait *sync.WaitGroup, closed chan struct{}) {
	wait.Add(1)
	defer wait.Done()

	client.publish("os", "active", true)
	client.subscribe("os/command", func(client mqtt.Client, message mqtt.Message) { osControl(message.Payload()) })

	for {
		time.Sleep(1 * time.Second)

		select {
		case <-closed:
			client.publish("os", "inactive", true)
			return
		default:

		}
	}
}

func oscamServiceHandler(client MqttClient, wait *sync.WaitGroup, closed chan struct{}) {
	wait.Add(1)
	defer wait.Done()

	oscamService := Service{"oscam"}
	client.subscribe("oscam/command", func(client mqtt.Client, message mqtt.Message) { oscamService.setStatePayload(message.Payload()) })

	var status, _, statusText = oscamService.checkStatus()
	client.publish("oscam", statusText, true)

	for {
		time.Sleep(1 * time.Second)

		select {
		case <-closed:
			return
		default:
			var newStatus, _, newStatusText = oscamService.checkStatus()
			if newStatus != status {
				status = newStatus
				client.publish("oscam", newStatusText, true)
			}
		}
	}
}
