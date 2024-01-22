package main

import (
	"flag"
	"fmt"

	"open_im_sdk/ws_wrapper/test"
	"open_im_sdk/ws_wrapper/test/client"
)

var jssdkURL = flag.String("url", "ws://192.168.100.188:10003/", "jssdk URL")
var imAPI = flag.String("api", "http://192.168.100.99:10002", "openIM api")
var connNum = flag.Int("connNum", 400, "conn num")

func main() {

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIweGNiZDAzM2VhM2MwNWRjOTUwNDYxMDA2MWM4NmM3YWUxOTFjNWM5MTMiLCJQbGF0Zm9ybSI6IldlYiIsImV4cCI6MTY5MjA3ODAyOCwibmJmIjoxNjg0MzAxNzI4LCJpYXQiOjE2ODQzMDIwMjh9.waXSsw5q7BZn5I739h_4zza9ICu2PleN8xHst-jNPhc"
	userId := "0xcbd033ea3c05dc9504610061c86c7ae191c5c913"

	fmt.Printf("simulation js client, user num: %d, jssdkURL:%s, apiURL:%s \n\n", *connNum, *jssdkURL, *imAPI)
	client.NewIMClient(token, "openIMAdmin", *imAPI, *jssdkURL, 1)

	test.StartSimulationJSClient(token, *imAPI, *jssdkURL, userId)

	// var err error
	// admin.Token, err = admin.GetToken()
	// if err != nil {
	// 	panic(err)
	// }
	// uidList, err := admin.GetALLUserIDList()
	// if err != nil {
	// 	panic(err)
	// }
	// l := uidList[0:*connNum]
	// // l = []string{"MTc3MjYzNzg0Mjg="}
	// for num, userID := range l {
	// 	time.Sleep(time.Millisecond * 500)
	// 	go test.StartSimulationJSClient(*imAPI, *jssdkURL, userID, num, l)
	// }

	// for {
	// 	time.Sleep(time.Second * 150)
	// 	fmt.Println("jssdk simulation is running, total num:", test.TotalSendMsgNum)
	// }

}
