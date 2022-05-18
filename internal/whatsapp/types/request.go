package types

type RequestLogin struct {
	Output string
}

type RequestSendMessage struct {
	RJID     string
	Message  string
	ViewOnce bool
}

type RequestSendLocation struct {
	RJID      string
	Latitude  float64
	Longitude float64
}

type RequestSendContact struct {
	RJID  string
	Name  string
	Phone string
}

type RequestSendLink struct {
	RJID    string
	Caption string
	URL     string
}
