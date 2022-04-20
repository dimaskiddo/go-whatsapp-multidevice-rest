package types

type RequestLogin struct {
	Output string
}

type RequestSendMessage struct {
	RJID    string
	Message string
}

type RequestSendLocation struct {
	RJID      string
	Latitude  float64
	Longitude float64
}
