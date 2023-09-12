package mqtt

import (
	"fmt"
	"os"
	"vitalsign-publisher/config"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/fatih/color"
)

var (
	conf         = config.GetConfig()
	clientOrigin MQTT.Client
	clientWeb    MQTT.Client
)

func init() {
	brokerOrigin := fmt.Sprintf("tcp://%s:%d", conf.Wisepaas.Host, conf.Wisepaas.Port)
	opts := MQTT.NewClientOptions()
	opts.AddBroker(brokerOrigin)
	opts.SetUsername(os.Getenv("MQTT_USER"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetClientID("Origin")

	clientOrigin = MQTT.NewClient(opts)
	if token := clientOrigin.Connect(); token.Wait() && token.Error() != nil {
		color.Red("MQTT clientOrigin connect FAIL: %s", token.Error())
		panic(token.Error())
	} else {
		color.Green("MQTT clientOrigin connect %v: %s", clientOrigin.IsConnected(), brokerOrigin)
	}

	brokerWeb := fmt.Sprintf("tcp://%s:%d", conf.Wisepaas.Host, conf.Wisepaas.Websocket)
	optsA := MQTT.NewClientOptions()
	optsA.AddBroker(brokerWeb)
	optsA.SetUsername(os.Getenv("MQTT_USER"))
	optsA.SetPassword(os.Getenv("MQTT_PASSWORD"))
	optsA.SetClientID("Web")

	clientWeb = MQTT.NewClient(optsA)
	if token := clientWeb.Connect(); token.Wait() && token.Error() != nil {
		color.Red("MQTT clientWeb connect FAIL: %s", token.Error())
		panic(token.Error())
	} else {
		color.Green("MQTT clientWeb connect %v: %s", clientWeb.IsConnected(), brokerWeb)
	}
}

func GetClient(name string) MQTT.Client {
	if name == "Web" {
		return clientWeb
	}
	return clientOrigin
}

func MQTTPublish(client MQTT.Client, topic string, payload interface{}) {
	if client.IsConnected() {
		client.Publish(topic, byte(0), false, payload)
	} else {
		rd := client.OptionsReader()
		color.Red("client %s is unconnected - AutoReconnect setting: %v", rd.ClientID(), rd.AutoReconnect())
	}
}
