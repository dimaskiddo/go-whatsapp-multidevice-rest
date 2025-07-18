package util

import "net/url"

func IsValidURL(rawUrl string) bool {
	parsedURL, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return false
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}

func MaskedJID(jid string) string {
	if len(jid) < 4 {
		return jid
	}
	return jid[:len(jid)-4] + "xxxx"
}
