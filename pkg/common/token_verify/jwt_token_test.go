package token_verify

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetClaimFromToken(t *testing.T) {

	got, err := GetClaimFromToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIweDIwYjE3ZGI5YTY1ZDhhNDhmMTQzMzhlZmQ4MjlhNWM3N2VmODMwMmQiLCJQbGF0Zm9ybSI6IldlYiIsImV4cCI6MTY4NDgyMzMxNSwibmJmIjoxNjg0ODIzMDE1LCJpYXQiOjE2ODQ4MjMzMTV9.tgNtPZay7-hGeRHHZO1FUllZ--IACayqopk4ZblWgxE")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		bytedata, _ := json.Marshal(got)
		fmt.Println(string(bytedata))
	}

}
