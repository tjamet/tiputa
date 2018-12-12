package credentials

import (
	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

// Getter is the interface to implement to get kubernetes credentials
type Getter interface {
	Get() (*clientauthentication.ExecCredentialStatus, error)
}

// NotProvidedError should be returned when the provider was not provided any credentials
type NotProvidedError struct{}

func (c NotProvidedError) Error() string {
	return "credentials not provided"
}
