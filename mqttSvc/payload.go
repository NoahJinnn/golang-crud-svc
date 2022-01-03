package mqttSvc

import (
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

func ServerDeleteDoorlockPayload(doorlock *models.Doorlock) string {
	return fmt.Sprintf(`{"door_id":%s}`, doorlock.ID)
}

func ServerUpdateDoorlockPayload(doorlock *models.Doorlock) string {
	return fmt.Sprintf(`{"door_id":%s,"description":%s,"location":%s, "state":%s}`, doorlock.ID, doorlock.Description, doorlock.Location, doorlock.State)
}

func ServerCmdDoorlockPayload(doorlockId string, cmd string) string {
	return fmt.Sprintf(`{"door_id":%s,"action":%s}`, doorlockId, cmd)
}

func ServerUpdateGatewayPayload(gw *models.Gateway) string {
	return fmt.Sprintf(`{"gateway_id":%s,"area_id":%s,"name":%s, "state":%s}`, gw.ID, gw.AreaID, gw.Name, gw.State)
}

func ServerCreatePasswordPayload(pw *models.Password) string {
	return fmt.Sprintf(`{"user_id":%s,"password_id":%s,"password_type":%s,"password_hash":%s}`, pw.UserID, pw.ID, pw.PasswordType, pw.PasswordHash)
}

func ServerUpdatePasswordPayload(pw *models.Password) string {
	return fmt.Sprintf(`{"password_id":%s,"password_type":%s,"password_hash":%s}`, pw.ID, pw.PasswordType, pw.PasswordHash)
}

func ServerDeletePasswordPayload(pw *models.Password) string {
	return fmt.Sprintf(`{"password_id":%s}`, pw.ID)
}
