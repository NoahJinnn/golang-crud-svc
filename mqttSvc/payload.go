package mqttSvc

import (
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

func ServerDeleteDoorlockPayload(doorlock *models.DoorlockDelete) string {
	return PayloadWithGatewayId(doorlock.GatewayID, fmt.Sprintf(`{"door_id":%d}`, doorlock.ID))
}

func ServerUpdateDoorlockPayload(doorlock *models.Doorlock) string {
	return PayloadWithGatewayId(doorlock.GatewayID,
		fmt.Sprintf(`{"door_id":%d,"description":%s,"location":%s, "state":%s}`,
			doorlock.ID, doorlock.Description, doorlock.Location, doorlock.State))
}

func ServerCmdDoorlockPayload(gwId string, doorlockId uint, cmd string) string {
	return PayloadWithGatewayId(gwId, fmt.Sprintf(`{"door_id":%d,"action":%s}`, doorlockId, cmd))
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

func PayloadWithGatewayId(gwId string, content string) string {
	return fmt.Sprintf(`{"gateway_id":%s,message:%s}`, gwId, content)
}
