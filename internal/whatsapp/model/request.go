package model

type ReqLogin struct {
	Output string
}

type ReqSendMessage struct {
	MSISDN        string
	Message       string
	QuotedID      string
	QuotedMessage string
}

type ReqSendLocation struct {
	MSISDN        string
	Latitude      float64
	Longitude     float64
	QuotedID      string
	QuotedMessage string
}
