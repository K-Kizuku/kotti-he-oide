package valueobject

import (
	"fmt"
	"net/url"
	"strings"
)

type PushEndpoint struct {
	value string
}

func NewPushEndpoint(endpoint string) (PushEndpoint, error) {
	if endpoint == "" {
		return PushEndpoint{}, fmt.Errorf("endpoint cannot be empty")
	}

	// Validate URL format
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return PushEndpoint{}, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	// Check if it's HTTPS (required for Web Push)
	if parsedURL.Scheme != "https" {
		return PushEndpoint{}, fmt.Errorf("endpoint must use HTTPS")
	}

	// Check if it's from a known push service provider
	host := strings.ToLower(parsedURL.Host)
	validProviders := []string{
		"fcm.googleapis.com",
		"android.googleapis.com", 
		"updates.push.services.mozilla.com",
		"web.push.apple.com",
	}

	isValid := false
	for _, provider := range validProviders {
		if strings.Contains(host, provider) {
			isValid = true
			break
		}
	}

	if !isValid {
		return PushEndpoint{}, fmt.Errorf("unknown push service provider: %s", host)
	}

	return PushEndpoint{value: endpoint}, nil
}

func (e PushEndpoint) Value() string {
	return e.value
}

func (e PushEndpoint) String() string {
	return e.value
}

func (e PushEndpoint) Equals(other PushEndpoint) bool {
	return e.value == other.value
}