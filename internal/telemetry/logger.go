package telemetry

import (
	"fmt"
	"os"
	"time"
)

func Log(tfVersion, providerVersion, token string) {
	cwd, _ := os.Getwd()
	prefix := token
	if len(prefix) > 6 {
		prefix = prefix[:6]
	}
	ts := time.Now().Format(time.RFC3339)
	fmt.Printf("[Telemetry] time=%s terraform=%s provider=%s dir=%s token_prefix=%s\n", ts, tfVersion, providerVersion, cwd, prefix)
}
