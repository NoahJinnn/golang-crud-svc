package mqttSvc

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func HandleMqttErr(t *mqtt.Token) error {
	<-(*t).Done() // Can also use '<-t.Done()' in releases > 1.2.0
	if (*t).Error() != nil {
		fmt.Println((*t).Error()) // Use your preferred logging technique (or just fmt.Printf)
		return (*t).Error()
	}
	return nil
}
