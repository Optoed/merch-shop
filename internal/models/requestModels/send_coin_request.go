package requestModels

type SendCoinRequest struct {
	SenderName   string `json:"name"` //TODO проверь что просто name
	ReceiverName string `json:"receiver"`
	Amount       int    `json:"amount"`
}
