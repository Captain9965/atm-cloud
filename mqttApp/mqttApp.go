package mqttApi

import (
	"fmt"

	"log"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func DisplayMqtt() {
	fmt.Println("Launching mqtt app")
}

func Connect(clientId string, uri string, username string, password string) mqtt.Client {
	opts := CreateClientOptions(clientId, uri, username, password)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client

}

func CreateClientOptions(clientId string, uri string, username string, password string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(uri)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func Listen(uri string, topic string, clientId string, username string, password string) {
	client := Connect(clientId, uri, username, password)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})

	fmt.Println("Subscription successful")
}
