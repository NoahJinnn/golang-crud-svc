package mqttSvc

const (
	TOPIC_GW_LOG_C           string = "gateway/log/create"
	TOPIC_GW_DOORLOCK_STATUS string = "gateway/doorlock/status"
	TOPIC_GW_DOORLOCK_C      string = "gateway/doorlock/create"
	TOPIC_GW_DOORLOCK_U      string = "gateway/doorlock/update"
	TOPIC_GW_DOORLOCK_D      string = "gateway/doorlock/delete"

	TOPIC_GW_BOOTUP   string = "gateway/bootup"
	TOPIC_GW_SHUTDOWN string = "gateway/shutdown"

	TOPIC_SV_DOORLOCK_CMD string = "server/doorlock/cmd"
	TOPIC_SV_DOORLOCK_U   string = "server/doorlock/update"
	TOPIC_SV_DOORLOCK_D   string = "server/doorlock/delete"
	TOPIC_SV_GATEWAY_U    string = "server/gateway/update"
)
