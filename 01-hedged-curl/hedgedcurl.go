package main

import (
	"os"
	"hedgedcurl/usage"
	"hedgedcurl/curl"
)


func main() {
	if usage.Help || len(os.Args) == 1 {
		usage.Reference()
		return
	}
	
	urls := usage.Arguments()
	code := curl.FormatHedged(urls)
	os.Exit(code)
}
