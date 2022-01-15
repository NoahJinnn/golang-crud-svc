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

// TODO: Guarantee mqtt req/res
// var DoorlockStateCheck = make(chan bool)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func MqttClient(host string, port string, logSvc *models.LogSvc, doorlockSvc *models.DoorlockSvc, gwSvc *models.GatewaySvc, empSvc *models.EmployeeSvc) mqtt.Client {

	mqtt.ERROR = log.New(os.Stdout, "[MQTT-ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[MQTT-CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[MQTT-WARN]  ", 0)

	//* Enable when need DEBUG
	// mqtt.DEBUG = log.New(os.Stdout, "[MQTT-DEBUG] "+host+":"+port, 0)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", host, port))
	opts.SetClientID("go_mqtt_client_10")
	// opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subGateway(client, logSvc, doorlockSvc, gwSvc, empSvc)

	return client
}

func subGateway(client mqtt.Client, logSvc *models.LogSvc, doorlockSvc *models.DoorlockSvc, gwSvc *models.GatewaySvc, empSvc *models.EmployeeSvc) {

	t := client.Subscribe(string(TOPIC_GW_SHUTDOWN), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		gwId := gjson.Get(payloadStr, "gateway_id")
		gwMsg := gjson.Get(payloadStr, "message")
		fmt.Printf("[MQTT-INFO] Gateway %s is shutdown with message: %s", gwId, gwMsg)
		gwSvc.DeleteGateway(context.Background(), gwId.String())
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_SHUTDOWN)
	}

	t = client.Subscribe(string(TOPIC_GW_BOOTUP), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())
		fmt.Println(payloadStr)
		newGw := &models.Gateway{}
		gwId := gjson.Get(payloadStr, "gateway_id")
		checkGw, _ := gwSvc.FindGatewayByMacID(context.Background(), gwId.String())
		if checkGw == nil {
			newGw.GatewayID = gwId.String()
			gwSvc.CreateGateway(context.Background(), newGw)
		}

		doorlocks := gjson.Get(payloadStr, "message.doorlocks")
		if doorlocks.Exists() {
			for _, v := range doorlocks.Array() {
				doorID := v.Get("doorlock_id")
				location := v.Get("location")
				description := v.Get("description")

				dl := &models.Doorlock{
					DoorSerialID: doorID.String(),
					Location:     location.String(),
					GatewayID:    gwId.String(),
					Description:  description.String(),
				}

				checkDl, _ := doorlockSvc.FindDoorlockBySerialID(context.Background(), doorID.String())
				if checkDl == nil {
					doorlockSvc.CreateDoorlock(context.Background(), dl)
				}
			}
		}

		hpEmployees, err := empSvc.FindAllHPEmployee(context.Background())
		if err != nil {
			fmt.Println(err.Error())
		}

		t := client.Publish(TOPIC_SV_HP_BOOTUP, 1, false, ServerBootuptHPEmployeePayload(gwId.String(), hpEmployees))
		HandleMqttErr(&t)

	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_BOOTUP)
	}

	t = client.Subscribe(string(TOPIC_GW_LOG_C), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())

		logMsg := gjson.Get(payloadStr, "message").String()
		gatewayId := gjson.Get(logMsg, "gateway_id")
		logType := gjson.Get(logMsg, "log_type")
		content := gjson.Get(logMsg, "log_data")
		logTime := gjson.Get(logMsg, "log_time")
		fmt.Printf(" %s: %s \n", msg.Topic(), payloadStr)
		logSvc.CreateGatewayLog(context.Background(), &models.GatewayLog{
			GatewayID: gatewayId.String(),
			LogType:   logType.String(),
			Content:   content.String(),
			LogTime:   logTime.String(),
		})
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("Subscribed to topic %s", TOPIC_GW_LOG_C)
	}

	// TODO: Guarantee mqtt req/res
	// t = client.Subscribe(string(TOPIC_GW_DOORLOCK_STATUS), 1, func(client mqtt.Client, msg mqtt.Message) {
	// 	var payloadStr = string(msg.Payload())
	// 	status := gjson.Get(payloadStr, "status")
	// 	action := gjson.Get(payloadStr, "action")
	// 	fmt.Printf("[MQTT-INFO] Action %s is %s", action, status)

	// 	DoorlockStateCheck <- true
	// })

	// if err := HandleMqttErr(&t); err == nil {
	// 	fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_DOORLOCK_STATUS)
	// }

	t = client.Subscribe(string(TOPIC_GW_DOORLOCK_U), 1, func(client mqtt.Client, msg mqtt.Message) {
		var payloadStr = string(msg.Payload())

		doorStateMsg := gjson.Get(payloadStr, "message").String()
		doorId := gjson.Get(doorStateMsg, "doorlock_id")
		state := gjson.Get(doorStateMsg, "state")
		doorlockSvc.UpdateDoorlockState(context.Background(), &models.DoorlockCmd{
			DoorSerialID: doorId.String(),
			State:        state.String(),
		})
		// TODO: Guarantee mqtt req/res
		// DoorlockStateCheck <- true
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
		doorId := gjson.Get(payloadStr, "doorlock_id")
		doorlockSvc.DeleteDoorlock(context.Background(), doorId.String())
	})

	if err := HandleMqttErr(&t); err == nil {
		fmt.Printf("[MQTT-INFO] Subscribed to topic %s", TOPIC_GW_DOORLOCK_D)
	}
}

func parseDoorlockPayload(msg mqtt.Message) *models.Doorlock {
	var payloadStr = string(msg.Payload())
	doorId := gjson.Get(payloadStr, "doorlock_id")
	description := gjson.Get(payloadStr, "description")
	location := gjson.Get(payloadStr, "location")
	fmt.Printf(" %s: %s \n", msg.Topic(), payloadStr)

	dl := &models.Doorlock{
		Description: description.String(),
		Location:    location.String(),
	}
	dl.ID = uint(doorId.Uint())
	return dl
}
