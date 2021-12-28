package mqttSvc

import (
	"fmt"

	"github.com/ecoprohcm/DMS_BackendServer/models"
)

func ServerUpdateDoorlockPayload(doorlock *models.Doorlock) string {
	return fmt.Sprintf(`{"door_id":%s,"description":%s,"location":%s, "state":%s}`, doorlock.ID, doorlock.Description, doorlock.Location, doorlock.State)
}

func ServerCmdDoorlockPayload(doorlockId string, cmd string) string {
	return fmt.Sprintf(`{"door_id":%s,"action":%s}`, doorlockId, cmd)
}
