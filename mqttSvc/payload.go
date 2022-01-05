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

func ServerCreatePasswordPayload(pw *models.PasswordCreate) string {
	return PayloadWithGatewayId(pw.GatewayID,
		fmt.Sprintf(`{"user_id":%d,"password_id":%d,"password_type":%s,"password_hash":%s}`,
			pw.UserID, pw.ID, pw.PasswordType, pw.PasswordHash))
}

func ServerUpdatePasswordPayload(pw *models.Password) string {
	return PayloadWithGatewayId("0",
		fmt.Sprintf(`"password_id":%d,"password_type":%s,"password_hash":%s}`,
			pw.ID, pw.PasswordType, pw.PasswordHash))
}

func ServerDeletePasswordPayload(pwId uint) string {
	return PayloadWithGatewayId("0",
		fmt.Sprintf(`"password_id":%d}`, pwId))
}

func ServerUpsertRegisterPayload(ssu models.StudentSchedulerUpsert, userId string) string {
	sche := ssu.Scheduler
	registerTime := fmt.Sprintf(`{"start_date":%s,"stop_date":%s,"week_day":%d,"start_class":%d,"end_class":%d}`,
		sche.StartDate, sche.EndDate, sche.WeekDay, sche.StartClassTime, sche.EndClassTime)
	return fmt.Sprintf(`{"register_id":%d,"gateway_id":%s,"doorlock_id":%s, "user_id":%s, "register_time":%s}`,
		sche.ID, ssu.GatewayID, ssu.DoorlockID, userId, registerTime)
}

func ServerDeleteRegisterPayload(ssu models.StudentSchedulerUpsert, userId string) string {
	sche := ssu.Scheduler
	return fmt.Sprintf(`{"register_id":%d,"gateway_id":%s,"doorlock_id":%s, "user_id":%s}`,
		sche.ID, ssu.GatewayID, ssu.DoorlockID, userId)
}

func PayloadWithGatewayId(gwId string, content string) string {
	return fmt.Sprintf(`{"gateway_id":%s,message:%s}`, gwId, content)
}
