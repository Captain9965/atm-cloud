package mqttApi

import (
	"fmt"
	"strings"

	"log"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"encoding/json"
)

var mqttClient *mqtt.Client

type event_t struct {
	Ev string `json:"ev"`
}

type hello_t struct {
	Ms  string `json:"ms"`
	Ver string `json:"ver"`
}

func RunMqtt() {
	NewClient := Connect("master", "broker.emqx.io:1883", "lenny", "Lenny123")
	mqttClient = &NewClient
	Listen(mqttClient, "w/s/#")
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

func Listen(client *mqtt.Client, topic string) {
	(*client).Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// log mqtt message:
		log.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
		go handleMessage(msg.Payload(), msg.Topic())
	})

	fmt.Println("Subscription successful")
}

func handleMessage(message []byte, topic string) {
	var event event_t
	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("%s", err.Error())
	}

	// Handle different Event types:
	switch event.Ev {
	case "hello":
		var helloData hello_t
		if err := json.Unmarshal(message, &helloData); err != nil {
			log.Printf("%s", err.Error())
			return
		}
		handleHelloEvent(&helloData, topic)
	default:
		log.Printf("%s Event not supported", event.Ev)

	}
}

func handleHelloEvent(data *hello_t, topic string) {
	//get settings for different machine versions:
	switch data.Ver {
		case "1.0.0":
			deviceUid, ok := getUidFromTopic(topic)
			if !ok {
				deviceUid = "None"
				log.Printf("Unable to get uid from %s \n", topic)
			}
			(*mqttClient).Publish("w/s/"+deviceUid, 0, false, `{"ev":"config","admin id":"5678","vending id":"5783","service id":"8990","cash":"6789","taps":"1111"}`)
		default:
			log.Printf("%s not supported", data.Ver)

	}
}

func getUidFromTopic(topic string) (string, bool) {
	parts := strings.Split(topic, "/")
	if len(parts) > 1 {
		return parts[len(parts)-1], true
	} else {
		return "", false
	}
}

func Publish(topic string, payload interface{}){
	fmt.Printf("Sending --> %s", payload)
	(*mqttClient).Publish(topic, 0, false, payload)
}
