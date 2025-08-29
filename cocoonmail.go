package cocoonmail

import (
	"errors"
	"net/url"

	"github.com/cocoonmail/cocoonmail-go/rest"
)

// CocoonmailOptions for CreateRequest
type CocoonmailOptions struct {
	Key      string
	Endpoint string
	Host     string
	Subuser  string
}

// cocoonmail host map for different regions
var allowedRegionsHostMap = map[string]string{
	"eu":     "https://api.eu.cocoonmail.com",
	"global": "https://webhook.cocoonmail.com",
}

// GetRequest
// @return [Request] a default request object
func GetRequest(key, endpoint, host string) rest.Request {
	return createCocoonmailRequest(CocoonmailOptions{key, endpoint, host, ""})
}

// GetRequestSubuser like GetRequest but with On-Behalf of Subuser
// @return [Request] a default request object
func GetRequestSubuser(key, endpoint, host, subuser string) rest.Request {
	return createCocoonmailRequest(CocoonmailOptions{key, endpoint, host, subuser})
}

// createCocoonmailRequest create Request
// @return [Request] a default request object
func createCocoonmailRequest(sgOptions CocoonmailOptions) rest.Request {
	options := options{
		"Bearer " + sgOptions.Key,
		sgOptions.Endpoint,
		sgOptions.Host,
		sgOptions.Subuser,
	}

	if options.Host == "" {
		options.Host = "https://webhook.cocoonmail.com"
	}

	return requestNew(options)
}

// NewSendClient constructs a new Cocoonmail client given an API key
func NewSendClient(key string) *Client {
	request := GetRequest(key, "/webhook/mail/send", "")
	request.Method = "POST"
	return &Client{request}
}

// extractEndpoint extracts the endpoint from a baseURL
func extractEndpoint(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "", err
	}

	return parsedURL.Path, nil
}

// SetDataResidency modifies the host as per the region
/*
 * This allows support for global and eu regions only. This set will likely expand in the future.
 * Global should be the default
 * Global region means the message should be sent through:
 * HTTP: api.cocoonmail.com
 * EU region means the message should be sent through:
 * HTTP: api.eu.cocoonmail.com
 */
// @return [Request] the modified request object
func SetDataResidency(request rest.Request, region string) (rest.Request, error) {
	regionalHost, present := allowedRegionsHostMap[region]
	if !present {
		return request, errors.New("error: region can only be \"eu\" or \"global\"")
	}
	endpoint, err := extractEndpoint(request.BaseURL)
	if err != nil {
		return request, err
	}
	request.BaseURL = regionalHost + endpoint
	return request, nil
}
