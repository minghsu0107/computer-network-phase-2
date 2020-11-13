package main

import "fmt"

func setHeader(resp, respHeader string) string {
	return fmt.Sprintf("%s%s\r\n", resp, respHeader)
}

func setBody(resp, body string) string {
	return fmt.Sprintf("%s\r\n%s", resp, body)
}
