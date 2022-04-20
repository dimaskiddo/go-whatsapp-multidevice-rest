package types

type ResponseLogin struct {
	QRCode  string `json:"qrcode"`
	Timeout int    `json:"timeout"`
}
