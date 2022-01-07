package mqttSvc

import (
	"encoding/json"
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

type BootupHPEmployee struct {
	userId     string `json:"user_id"`
	rfidPass   string `json:"rfid_pw"`
	keypadPass string `json:"keypad_pw"`
}

func ServerDeleteDoorlockPayload(doorlock *models.DoorlockDelete) string {
	msg := fmt.Sprintf(`{"door_id":%s}`, doorlock.DoorSerialID)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerUpdateDoorlockPayload(doorlock *models.Doorlock) string {
	msg := fmt.Sprintf(`{"door_id":%s,"description":%s,"location":%s, "state":%s}`,
		doorlock.DoorSerialID, doorlock.Description, doorlock.Location, doorlock.State)
	return PayloadWithGatewayId(doorlock.GatewayID, msg)
}

func ServerCmdDoorlockPayload(gwId string, doorSerialId string, cmd string) string {
	msg := fmt.Sprintf(`{"door_id":%s,"action":%s}`, doorSerialId, cmd)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerUpdateGatewayPayload(gw *models.Gateway) string {
	return fmt.Sprintf(`{"gateway_id":%d,"area_id":%d,"name":%s, "state":%s}`,
		gw.ID, gw.AreaID, gw.Name, gw.State)
}

func ServerUpsertRegisterPayload(usu models.UserSchedulerUpsert, rfidPass string, keypadPass string, userId string) string {
	sche := usu.Scheduler
	msg := fmt.Sprintf(`{"register_id":%d,"doorlock_id":%s, "user_id":%s, "rfid_pw":%s,"keypad_pw":%s,
	"start_date":%s,"stop_date":%s,"week_day":%d,"start_class":%d,"end_class":%d}`,
		sche.ID, usu.DoorlockID, userId, rfidPass, keypadPass,
		sche.StartDate, sche.EndDate, sche.WeekDay, sche.StartClassTime, sche.EndClassTime)
	return PayloadWithGatewayId(usu.GatewayID, msg)
}

func ServerDeleteRegisterPayload(usu models.UserSchedulerUpsert, userId string) string {
	sche := usu.Scheduler
	msg := fmt.Sprintf(`{"register_id":%d,"doorlock_id":%s, "user_id":%s}`,
		sche.ID, usu.DoorlockID, userId)
	return PayloadWithGatewayId(usu.GatewayID, msg)
}

func ServerBootuptHPEmployeePayload(gwId string, emps []models.Employee) string {
	bootupEmps := []BootupHPEmployee{}
	for _, emp := range emps {
		buEmp := BootupHPEmployee{
			userId:     emp.MSNV,
			rfidPass:   emp.RfidPass,
			keypadPass: emp.KeypadPass,
		}
		bootupEmps = append(bootupEmps, buEmp)
	}
	jsonInfo, _ := json.Marshal(bootupEmps)
	return PayloadWithGatewayId(gwId, string(jsonInfo))
}

func ServerUpsertHPEmployeePayload(gwId string, emp models.Employee) string {
	msg := fmt.Sprintf(`{"user_id":%s,"rfid_pw":%s, "keypad_pw":%s}`,
		emp.MSNV, emp.RfidPass, emp.KeypadPass)
	return PayloadWithGatewayId(gwId, msg)
}

func ServerDeleteHPEmployeePayload(gwId string, emp models.Employee) string {
	msg := fmt.Sprintf(`{"user_id":%s}`, emp.MSNV)
	return PayloadWithGatewayId(gwId, msg)
}

func PayloadWithGatewayId(gwId string, msg string) string {
	return fmt.Sprintf(`{"gateway_id":%s,message:%s}`, gwId, msg)
}
