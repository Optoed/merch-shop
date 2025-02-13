package requestModels

type SendCoinRequest struct {
	ReceiverName string `json:"receiver_name"`
	Amount       int    `json:"amount"`
}
