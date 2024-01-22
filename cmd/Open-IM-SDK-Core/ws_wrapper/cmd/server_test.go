package main

import (
	"flag"
	"testing"
)

func TestRun(t *testing.T) {
	flag.Set("openIM_api_address", "http://192.168.100.99:10002")
	flag.Set("openIM_ws_address", "ws://192.168.100.99:10001")
	flag.Set("openIMDbDir", "../../../../db/sdk/")
	main()
}
