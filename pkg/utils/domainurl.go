package utils

import (
	"fmt"
	"net"
)

func GetDomainTxtList(domainName string) ([]string, error) {
	// 查询 TXT 记录
	txtRecords, err := net.LookupTXT(domainName)
	if err != nil {
		fmt.Println("Failed to lookup TXT records:", err)
		return nil, err
	}
	return txtRecords, nil
}
