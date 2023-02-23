package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

//topic := "/Mariano-Oscam/oscam/online"
//topic := "/Mariano-Oscam/os/online"

// https://guzalexander.com/2017/05/31/gracefully-exit-server-in-go.html

func main() {
	var configuration = readConfiguration("config.json")
	client := NewMqttClient(configuration)

	task := &Task{
		closed: make(chan struct{}),
		client: client,
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task.wg.Add(1)
	go func() { defer task.wg.Done(); task.Run() }()

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		task.Stop()
	}
}

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	client MqttClient
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			return
		}
	}
}

func (t *Task) Stop() {
	t.client.client.Disconnect(100)
	close(t.closed)
	t.wg.Wait()
}

// func sub(client MqttClient) {
// 	topic := "/Mariano-Oscam/oscam/online"
// 	token := client.client.Subscribe(topic, 1, nil)
// 	token.Wait()
// 	fmt.Printf("Subscribed to topic: %s", topic)
// }
