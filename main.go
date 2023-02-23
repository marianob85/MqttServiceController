package main

import (
	"fmt"
	"time"
)

//topic := "/Mariano-Oscam/oscam/online"
//topic := "/Mariano-Oscam/os/online"

func main() {
	go terminationSetup()
	var configuration = readConfiguration("config.json")
	client := NewMqttClient(configuration)
	for {
		time.Sleep(1 * time.Second)
	}

	//client.client.Disconnect(100)
}

func sub(client MqttClient) {
	topic := "/Mariano-Oscam/oscam/online"
	token := client.client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
