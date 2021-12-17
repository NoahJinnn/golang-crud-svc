package models

type DeleteGateway struct {
	GatewayID string `json:"gatewayId" binding:"required"`
}

func convertToDeleteGw(dgw *DeleteGateway) *Gateway {
	return &Gateway{
		GatewayID: dgw.GatewayID,
	}
}
