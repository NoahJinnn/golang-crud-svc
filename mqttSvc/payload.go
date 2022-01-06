package mqttSvc

import (
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

func ServerDeleteDoorlockPayload(doorlock *models.DoorlockDelete) string {
	return PayloadWithGatewayId(doorlock.GatewayID, fmt.Sprintf(`{"door_id":%s}`, doorlock.DoorSerialID))
}

func ServerUpdateDoorlockPayload(doorlock *models.Doorlock) string {
	return PayloadWithGatewayId(doorlock.GatewayID,
		fmt.Sprintf(`{"door_id":%s,"description":%s,"location":%s, "state":%s}`,
			doorlock.DoorSerialID, doorlock.Description, doorlock.Location, doorlock.State))
}

func ServerCmdDoorlockPayload(gwId string, doorSerialId string, cmd string) string {
	return PayloadWithGatewayId(gwId, fmt.Sprintf(`{"door_id":%s,"action":%s}`, doorSerialId, cmd))
}

func ServerUpdateGatewayPayload(gw *models.Gateway) string {
	return fmt.Sprintf(`{"gateway_id":%d,"area_id":%d,"name":%s, "state":%s}`,
		gw.ID, gw.AreaID, gw.Name, gw.State)
}

func ServerDeletePasswordPayload(pwId uint) string {
	return PayloadWithGatewayId("0",
		fmt.Sprintf(`"password_id":%d}`, pwId))
}

func ServerUpsertRegisterPayload(usu models.UserSchedulerUpsert, rfidPass string, keypadPass string, userId string) string {
	sche := usu.Scheduler
	registerTime := fmt.Sprintf(`{"start_date":%s,"stop_date":%s,"week_day":%d,"start_class":%d,"end_class":%d}`,
		sche.StartDate, sche.EndDate, sche.WeekDay, sche.StartClassTime, sche.EndClassTime)
	return fmt.Sprintf(`{"register_id":%d,"gateway_id":%s,"doorlock_id":%s, "user_id":%s, "rfid_pw":%s,"keypad_pw":%s,"register_time":%s}`,
		sche.ID, usu.GatewayID, usu.DoorlockID, userId, rfidPass, keypadPass, registerTime)
}

func ServerDeleteRegisterPayload(usu models.UserSchedulerUpsert, userId string) string {
	sche := usu.Scheduler
	return fmt.Sprintf(`{"register_id":%d,"gateway_id":%s,"doorlock_id":%s, "user_id":%s}`,
		sche.ID, usu.GatewayID, usu.DoorlockID, userId)
}

func PayloadWithGatewayId(gwId string, content string) string {
	return fmt.Sprintf(`{"gateway_id":%s,message:%s}`, gwId, content)
}
