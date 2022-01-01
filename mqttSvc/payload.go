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
