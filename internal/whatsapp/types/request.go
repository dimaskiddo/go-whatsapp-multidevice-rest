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

type RequestSendPoll struct {
	RJID        string
	Question    string
	Options     string
	MultiAnswer bool
}

type RequestMessage struct {
	RJID    string
	MSGID   string
	Message string
	Emoji   string
}

type RequestGroupJoin struct {
	Link string
}

type RequestGroupLeave struct {
	GID string
}
