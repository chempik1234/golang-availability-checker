package webhooks

import (
	"fmt"
	"time"
)

func FormHTTPBody(name string, result bool, dateTime time.Time) string {
	return fmt.Sprintf(
		`{"name":"%s","datetime":"%s","result":%t}`,
		name, dateTime.Format(time.RFC3339), result)
}
