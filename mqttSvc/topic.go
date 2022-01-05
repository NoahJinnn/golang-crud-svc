package mqttSvc

const (
	TOPIC_GW_LOG_C           string = "gateway/log/create"
	TOPIC_GW_DOORLOCK_STATUS string = "gateway/door/status"
	TOPIC_GW_DOORLOCK_C      string = "gateway/door/create"
	TOPIC_GW_DOORLOCK_U      string = "gateway/door/update"
	TOPIC_GW_DOORLOCK_D      string = "gateway/door/delete"

	TOPIC_GW_BOOTUP   string = "gateway/bootup"
	TOPIC_GW_SHUTDOWN string = "gateway/shutdown"

	TOPIC_SV_DOORLOCK_CMD string = "server/door/cmd"
	TOPIC_SV_DOORLOCK_U   string = "server/door/update"
	TOPIC_SV_DOORLOCK_D   string = "server/door/delete"
	TOPIC_SV_GATEWAY_U    string = "server/gateway/update"

	TOPIC_SV_PASSWORD_C string = "server/password/create"
	TOPIC_SV_PASSWORD_U string = "server/password/update"
	TOPIC_SV_PASSWORD_D string = "server/password/delete"

	TOPIC_SV_SCHEDULER_C string = "server/register/create"
	TOPIC_SV_SCHEDULER_U string = "server/register/update"
	TOPIC_SV_SCHEDULER_D string = "server/register/delete"
)
