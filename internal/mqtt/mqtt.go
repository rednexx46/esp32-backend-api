package mqtt

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

// InitMQTT initializes the MQTT client, connects to the broker, and subscribes to the specified topic.
func InitMQTT(handler mqtt.MessageHandler) {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://" + os.Getenv("MQTT_BROKER") + ":" + os.Getenv("MQTT_PORT")).
		SetUsername(os.Getenv("MQTT_USERNAME")).
		SetPassword(os.Getenv("MQTT_PASSWORD")).
		SetDefaultPublishHandler(handler)

	mqttClient = mqtt.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[MQTT] Failed to connect: %v", token.Error())
	}

	topic := os.Getenv("MQTT_TOPIC_SENSORS_DATA")
	if token := mqttClient.Subscribe(topic+"#", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("[MQTT] Failed to subscribe: %v", token.Error())
	}
	log.Printf("[MQTT] Subscribed to topic: %s", topic)
}
