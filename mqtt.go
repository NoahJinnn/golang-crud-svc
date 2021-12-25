package main

import (
	"context"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/tidwall/gjson"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func initMqttClient(host string, port string, logSvc *models.LogSvc) {
	// var broker = "iot.hcmue.space"
	// var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", host, port))
	opts.SetClientID("go_mqtt_client")
	// opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subGateway(client, logSvc)
	publish(client)

	client.Disconnect(250)
}

func subGateway(client mqtt.Client, logSvc *models.LogSvc) {
	var gwlogTopic = "gateway/log/create"
	token := client.Subscribe(gwlogTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		macId := gjson.Get(payloadStr, "macId")
		logType := gjson.Get(payloadStr, "logType")
		content := gjson.Get(payloadStr, "content")
		logTime := gjson.Get(payloadStr, "logTime")
		fmt.Printf("Yo %s: %s \n", msg.Topic(), payloadStr)
		fmt.Printf("Lo %s: %s %s %s %s\n", msg.Topic(), macId,
			logType,
			content,
			logTime)
		logSvc.CreateGatewayLog(&models.GatewayLog{
			MacID:   macId.String(),
			LogType: logType.String(),
			Content: content.String(),
			LogTime: logTime.String(),
		}, context.Background())
	})
	token.Wait()
	fmt.Printf("Subscribed to topic %s", gwlogTopic)
}

func publish(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		json := fmt.Sprintf(`{"macId":%d,"logType":"SUCCESS","content": "Unlock door", "logTime": "198273198237"}`, i)
		token := client.Publish("gateway/log/create", 0, false, json)
		token.Wait()
		time.Sleep(time.Second)
	}
}
