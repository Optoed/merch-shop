package requestModels

type SendCoinRequest struct {
	ReceiverName string `json:"toUser" binding:"required"`
	Amount       int    `json:"amount" binding:"required"`
}
