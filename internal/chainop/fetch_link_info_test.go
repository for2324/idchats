package chainop

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_fetchLinkInfo(t *testing.T) {
	got, _ := fetchLinkInfo("https://twitter.com/weirdlilguys/status/1675228298937237504?s=20")
	data, _ := json.Marshal(got)
	fmt.Printf(string(data))
}
