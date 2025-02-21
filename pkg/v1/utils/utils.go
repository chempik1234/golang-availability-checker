package utils

import "fmt"

func Url(protocol string, host string, port string) string {
	return fmt.Sprintf("%s://%s:%s/", protocol, host, port)
}
