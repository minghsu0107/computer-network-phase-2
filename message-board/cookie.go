package main

import (
	"fmt"
)

func setCookie(respHeader, key, val string) string {
	return fmt.Sprintf("%sSet-Cookie: %s=%s\r\n", respHeader, key, val)
}
