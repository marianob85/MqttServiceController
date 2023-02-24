package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client mqtt.Client
	topic  string
}

func NewMqttClient(config Config) MqttClient {
	var mqttClient MqttClient

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s:%d", config.Server, config.Port))
	opts.SetClientID(config.ClientID)
	opts.SetUsername(config.UserName)
	opts.SetPassword(config.Password)
	opts.OnConnect = mqttClient.connectHandler
	opts.OnConnectionLost = mqttClient.connectLostHandler
	opts.OnReconnecting = mqttClient.reconnectingHandler
	opts.ConnectRetry = true
	opts.ConnectRetryInterval = 1 * time.Second
	mqttClient.client = mqtt.NewClient(opts)
	mqttClient.topic = config.Topic
	if token := mqttClient.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return mqttClient
}

func (o *MqttClient) subscribe(topic string, callback mqtt.MessageHandler) {
	token := o.client.Subscribe(fmt.Sprintf("%s/%s", o.topic, topic), 0, callback)
	token.Wait()
}

func (o *MqttClient) publish(topic string, message string, retained bool) {
	token := o.client.Publish(fmt.Sprintf("%s/%s", o.topic, topic), 0, retained, message)
	token.Wait()
}

func (o *MqttClient) reconnectingHandler(client mqtt.Client, options *mqtt.ClientOptions) {
	fmt.Printf("Connect attempt\n")
}

func (o *MqttClient) connectHandler(client mqtt.Client) {
	fmt.Println("Connected")
}

func (o *MqttClient) connectLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}
