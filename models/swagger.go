package models

import "time"

type SwagCreateGateway struct {
	AreaID          uint      `json:"areaId"`
	MacID           string    `json:"macId"`
	Name            string    `json:"name"`
	LastConnectTime time.Time `json:"lastConnectTime"`
	State           string    `json:"state"`
}

type SwagUpateGateway struct {
	GormModel
	SwagCreateGateway
}

type SwagCreateArea struct {
	Gateway Gateway `json:"gateway"`
	Name    string  `json:"name"`
	Manager string  `json:"manager"`
}

type SwagUpdateArea struct {
	GormModel
	SwagCreateArea
}

type SwagCreateDoorlock struct {
	AreaID          uint      `json:"areaId"`
	GatewayID       uint      `json:"gatewayId"`
	SchedulerID     uint      `json:"schedulerId"`
	Description     string    `json:"description"`
	Location        string    `json:"location"`
	LastConnectTime time.Time `json:"lastConnectTime"`
	State           string    `json:"state"`
}

type SwagUpdateDoorlock struct {
	GormModel
	SwagCreateDoorlock
}
