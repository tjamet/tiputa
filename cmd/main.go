package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/tjamet/tiputa/credentials"
	"github.com/tjamet/tiputa/secret"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthentication "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

var getters = []credentials.Getter{}
var output io.Writer = os.Stdout

var rootCmd = &cobra.Command{
	Use:   "tiputa",
	Short: "Retrieve kubernetes credentials from a password manager",
	Long:  `Tiputa is a PoC to use pass to encrypt user authentication of kuberntes clients. It implements Kubernetes client-go credential plugins available in beta since kubernetes 1.11.`,
	Run: func(cmd *cobra.Command, args []string) {
		WriteCredentials(getters, output)
	},
}

func WriteCredentials(getters []credentials.Getter, w io.Writer) {
	for _, g := range getters {
		status, err := g.Get()
		if err != nil {
			if _, ok := err.(credentials.NotProvidedError); ok {
				continue
			}
			fmt.Println(err)
			os.Exit(1)
		}
		credentials := clientauthentication.ExecCredential{
			TypeMeta: metav1.TypeMeta{
				APIVersion: clientauthentication.SchemeGroupVersion.String(),
				Kind:       "ExecCredential",
			},
			Status: status,
		}
		e := json.NewEncoder(w)
		err = e.Encode(credentials)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

// AddCredentialGetter allows extending the features to provide a custom credentials getter from any other means
func AddCredentialGetter(g credentials.Getter) {
	getters = append(getters, g)
}

// RegisterSecretArgs adds command line options to instruct retrieving secrets from the secret reader
func RegisterSecretArgs(description, prefix string, reader secret.Reader) {
	getter := credentials.Secret{
		Reader: reader,
	}
	Flags().StringVar(&getter.TokenSecret, fmt.Sprintf("%s-token", prefix), "", fmt.Sprintf("the %s secret containing the kubernetes access token", description))
	Flags().StringVar(&getter.ClientCertificateSecret, fmt.Sprintf("%s-client-certificate", prefix), "", fmt.Sprintf("the %s secret containing the kubernetes client certificate", description))
	Flags().StringVar(&getter.ClientKeySecret, fmt.Sprintf("%s-client-key", prefix), "", fmt.Sprintf("the %s secret containing the kubernetes client private key", description))
	AddCredentialGetter(&getter)
}

// Flags exports the root command flags to allow hooking into the commang line to add other providers
func Flags() *flag.FlagSet {
	return rootCmd.Flags()
}

func init() {
	RegisterSecretArgs("password store", "pass", secret.DefaultPasswordStore)
}

// Execute runs the actual credential getter
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
