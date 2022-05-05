package types

type ResponseLogin struct {
	QRCode  string `json:"qrcode"`
	Timeout int    `json:"timeout"`
}

type ResponseSendMessage struct {
	MsgID string `json:"msgid"`
}
