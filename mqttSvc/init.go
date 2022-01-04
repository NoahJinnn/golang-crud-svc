package mqttSvc

import (
	"context"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ecoprohcm/DMS_BackendServer/models"
	"github.com/tidwall/gjson"
)

var DoorlockStateCheck = make(chan bool)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func MqttClient(host string, port string, logSvc *models.LogSvc, doorlockSvc *models.DoorlockSvc, gwSvc *models.GatewaySvc) mqtt.Client {

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
	subGateway(client, logSvc, doorlockSvc, gwSvc)

	return client
}

func subGateway(client mqtt.Client, logSvc *models.LogSvc, doorlockSvc *models.DoorlockSvc, gwSvc *models.GatewaySvc) {

	t := client.Subscribe(string(TOPIC_GW_SHUTDOWN), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		macId := gjson.Get(payloadStr, "mac_id")
		gwMsg := gjson.Get(payloadStr, "message")
		fmt.Printf("[MQTT-INFO] Gateway %s is shutdown with message: %s", macId, gwMsg)
		gw := &models.Gateway{
			MacID: macId.String(),
		}
		gwSvc.DeleteGatewayByMacId(context.Background(), gw)
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_SHUTDOWN)
	}

	t = client.Subscribe(string(TOPIC_GW_BOOTUP), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		macId := gjson.Get(payloadStr, "mac_id")
		users := gjson.Get(payloadStr, "message.users")
		if users.Exists() {
			for _, v := range users.Array() {
				fmt.Printf("%s %s", v.Get("name"), v.Get("age"))
			}
		}

		doorlocks := gjson.Get(payloadStr, "message.doorlocks")
		if doorlocks.Exists() {
			for _, v := range doorlocks.Array() {
				doorID := v.Get("door_id")
				dl, _ := doorlockSvc.FindDoorlockByID(context.Background(), doorID.String())
				if dl == nil {
					location := v.Get("location")
					dl = &models.Doorlock{
						State:     "Close",
						Location:  location.String(),
						GatewayID: macId.String(),
					}
					dl.ID = uint(doorID.Uint())
					doorlockSvc.CreateDoorlock(context.Background(), dl)
				} else {
					doorlockSvc.UpdateDoorlockGateway(context.Background(), dl, macId.Str)
				}
			}
		}

	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_BOOTUP)
	}

	t = client.Subscribe(string(TOPIC_GW_LOG_C), 1, func(client mqtt.Client, msg mqtt.Message) {
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

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("Subscribed to topic %s", TOPIC_GW_LOG_C)
	}

	t = client.Subscribe(string(TOPIC_GW_DOORLOCK_STATUS), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		status := gjson.Get(payloadStr, "status")
		action := gjson.Get(payloadStr, "action")
		fmt.Printf("[MQTT-INFO] Action %s is %s", action, status)

		DoorlockStateCheck <- true
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_DOORLOCK_STATUS)
	}

	t = client.Subscribe(string(TOPIC_GW_DOORLOCK_U), 1, func(client mqtt.Client, msg mqtt.Message) {
		dl := parseDoorlockPayload(msg)
		doorlockSvc.UpdateDoorlock(context.Background(), dl)
		DoorlockStateCheck <- true
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_DOORLOCK_U)
	}

	t = client.Subscribe(string(TOPIC_GW_DOORLOCK_C), 1, func(client mqtt.Client, msg mqtt.Message) {
		dl := parseDoorlockPayload(msg)
		doorlockSvc.CreateDoorlock(context.Background(), dl)
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_DOORLOCK_C)
	}

	t = client.Subscribe(string(TOPIC_GW_DOORLOCK_D), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		doorId := gjson.Get(payloadStr, "door_id")
		doorlockSvc.DeleteDoorlock(context.Background(), uint(doorId.Uint()))
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_DOORLOCK_D)
	}
}

func parseDoorlockPayload(msg mqtt.Message) *models.Doorlock {
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
	dl.ID = uint(doorId.Uint())
	return dl
}
