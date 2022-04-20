package types

type RequestLogin struct {
	Output string
}

type RequestSendMessage struct {
	RJID    string
	Message string
	// QuotedID      string
	// QuotedMessage string
}

type RequestSendLocation struct {
	RJID      string
	Latitude  float64
	Longitude float64
	// QuotedID      string
	// QuotedMessage string
}
