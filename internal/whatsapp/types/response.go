package types

type ResponseLogin struct {
	QRCode  string `json:"qrcode"`
	Timeout int    `json:"timeout"`
}

type ResponsePairing struct {
	PairCode string `json:"paircode"`
	Timeout  int    `json:"timeout"`
}

type ResponseSendMessage struct {
	MsgID string `json:"msgid"`
}
