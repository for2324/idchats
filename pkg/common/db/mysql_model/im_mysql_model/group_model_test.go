package im_mysql_model

import (
	"fmt"
	"testing"
)

func TestColligateSearch(t *testing.T) {
	// {
	// 	resp, err := ColligateSearch("13225")
	// 	if err != nil {
	// 		t.Errorf("ColligateSearch() error = %v", err)
	// 		return
	// 	}
	// 	fmt.Println(resp)
	// }
	{
		resp, err := ColligateSearch("0x")
		if err != nil {
			t.Errorf("ColligateSearch() error = %v", err)
			return
		}
		fmt.Println(resp)
	}

}
