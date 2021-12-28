package mqttSvc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/tidwall/gjson"
)

const (
	TOPIC_GW_LOG_C        string = "gateway/log/create"
	TOPIC_SV_DOORLOCK_R   string = "server/doorlock/read"
	TOPIC_SV_DOORLOCK_U   string = "server/doorlock/update"
	TOPIC_SV_DOORLOCK_CMD string = "server/doorlock/cmd"
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

func MqttClient(host string, port string, logSvc *models.LogSvc, doorlockSvc *models.DoorlockSvc) mqtt.Client {

	mqtt.ERROR = log.New(os.Stdout, "[MQTT-ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[MQTT-CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[MQTT-WARN]  ", 0)

	//Enable when need DEBUG
	// mqtt.DEBUG = log.New(os.Stdout, "[MQTT-DEBUG] "+host+":"+port, 0)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", host, port))
	opts.SetClientID("go_mqtt_client_1")
	// opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subGateway(client, logSvc, doorlockSvc)
	publishTest(client)

	// client.Disconnect(250)
	return client
}

func subGateway(client mqtt.Client, logSvc *models.LogSvc, doorlockSvc *models.DoorlockSvc) {
	t := client.Subscribe(string(TOPIC_GW_LOG_C), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		macId := gjson.Get(payloadStr, "macId")
		logType := gjson.Get(payloadStr, "logType")
		content := gjson.Get(payloadStr, "content")
		logTime := gjson.Get(payloadStr, "logTime")
		fmt.Printf(" %s: %s \n", msg.Topic(), payloadStr)
		logSvc.CreateGatewayLog(context.Background(), &models.GatewayLog{
			MacID:   macId.String(),
			LogType: logType.String(),
			Content: content.String(),
			LogTime: logTime.String(),
		})
	})

	t = client.Subscribe(string(TOPIC_SV_DOORLOCK_R), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		doorId := gjson.Get(payloadStr, "door_id")
		description := gjson.Get(payloadStr, "description")
		location := gjson.Get(payloadStr, "location")
		state := gjson.Get(payloadStr, "state")
		fmt.Printf(" %s: %s \n", msg.Topic(), payloadStr)

		dl := &models.Doorlock{
			Description: description.String(),
			Location:    location.String(),
			State:       state.String(),
		}
		dl.ID = doorId.String()
		doorlockSvc.UpdateDoorlock(context.Background(), dl)
	})

	// Test only
	t = client.Subscribe(string(TOPIC_SV_DOORLOCK_CMD), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		fmt.Printf(" %s: %s \n", msg.Topic(), payloadStr)
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("Subscribed to topic %s", TOPIC_GW_LOG_C)
	}
}

func publishTest(client mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		json := fmt.Sprintf(`{"macId":%d,"logType":"SUCCESS","content": "Unlock door", "logTime": "198273198237"}`, i)
		t := client.Publish(string(TOPIC_GW_LOG_C), 1, false, json)
		go HandleMqttErr(&t)
		time.Sleep(time.Second)
	}
}
