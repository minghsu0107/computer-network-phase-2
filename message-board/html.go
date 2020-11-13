package main

import "io/ioutil"

func readHTMLFile(filename string) string {
	buf, _ := ioutil.ReadFile(filename)
	return string(buf)
}
