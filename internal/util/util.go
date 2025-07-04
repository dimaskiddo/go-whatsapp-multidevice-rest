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
