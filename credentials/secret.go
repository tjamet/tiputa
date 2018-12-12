package credentials

import (
	"strings"

	"github.com/tjamet/tiputa/secret"
	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

type Secret struct {
	Reader                  secret.Reader
	TokenSecret             string
	ClientCertificateSecret string
	ClientKeySecret         string
}

func (s Secret) Get() (*clientauthentication.ExecCredentialStatus, error) {
	if s.TokenSecret != "" {
		t, err := s.Reader.Read(s.TokenSecret)
		if err != nil {
			return nil, err
		}
		return &clientauthentication.ExecCredentialStatus{
			Token: strings.Trim(t, "\n"),
		}, nil
	}
	if s.ClientCertificateSecret != "" && s.ClientKeySecret != "" {
		cert, err := s.Reader.Read(s.ClientCertificateSecret)
		if err != nil {
			return nil, err
		}
		key, err := s.Reader.Read(s.ClientKeySecret)
		if err != nil {
			return nil, err
		}
		return &clientauthentication.ExecCredentialStatus{
			ClientCertificateData: strings.Trim(cert, "\n"),
			ClientKeyData:         strings.Trim(key, "\n"),
		}, nil
	}
	return nil, NotProvidedError{}
}
