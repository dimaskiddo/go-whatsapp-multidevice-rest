package types

type ResponseLogin struct {
	Timeout string `json:"timeout"`
	QRCode  string `json:"qrcode"`
}
