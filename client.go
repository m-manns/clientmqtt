package main

import (
	"flag"
	"fmt"

	MqttLib "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Info("Client MQTT...")

	flagIP := flag.String("ip", "", "IP MQTT broker")
	flagPort := flag.String("port", "1883", "Port MQTT")
	flagUser := flag.String("username", "", "MQTT username")
	flagPass := flag.String("password", "", "MQTT password")
	flag.Parse()

	url := fmt.Sprintf("tcp://%v:%v", *flagIP, *flagPort)

	options := MqttLib.NewClientOptions().AddBroker(url)
	options.SetConnectionLostHandler(func(lost MqttLib.Client, err error) {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Connection to MQTT broker lost")
	})
	options.SetClientID(*flagUser)
	options.SetUsername(*flagUser)
	options.SetPassword(*flagPass)
	options.AutoReconnect = true
	options.OnConnect = func(c MqttLib.Client) {
		log.Info("Connected to MQTT broker")
	}

	client := MqttLib.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.WithFields(log.Fields{
			"error": token.Error(),
			"url":   url,
		}).Error("MQTT server connection error")
	} else {
		log.WithFields(log.Fields{
			"url": url,
		}).Debug("MQTT server connection")
	}

	log.WithFields(log.Fields{
		"is connected": client.IsConnected(),
	}).Info("Connection status")

	if client.IsConnected() {
		client.Disconnect(1)
	}
}
