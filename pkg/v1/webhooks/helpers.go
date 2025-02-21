package webhooks

import "fmt"

func FormHTTPBody(name string, result bool) string {
	return fmt.Sprintf("{\"name\":\"%s\", \"result\":%t}", name, result)
}
