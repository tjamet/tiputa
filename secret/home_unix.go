package secret

import "os"

func homeDir() string {
	return os.Getenv("HOME")
}
