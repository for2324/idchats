package order

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.etcd.io/etcd/client/v3/concurrency"
)

func TestChainTable(t *testing.T) {
	chainTable := NewChainTable("test", 1)
	backList := make([]string, 0)
	for i := 0; i < 11; i++ {
		nonce, backId, err := chainTable.Next()
		if err != nil {
			fmt.Printf("get nonce fail, %v \n", err)
			continue
		}
		fmt.Printf("nonce: %d \n", nonce)
		backList = append(backList, backId)
	}
	if len(backList) < 2 {
		return
	}
	chainTable.GiveBack(backList[0])
	chainTable.GiveBack(backList[0])
	chainTable.GiveBack(backList[1])
	for i := 0; i < 3; i++ {
		nonce, _, err := chainTable.Next()
		if err != nil {
			fmt.Printf("re get nonce fail, %v \n", err)
			continue
		}
		fmt.Printf("re nonce: %d \n", nonce)
	}
}

func TestEtcdLock(t *testing.T) {
	etcdClient, err := GetEtcdClient()
	if err != nil {
		fmt.Printf("get etcd client fail, %v \n", err)
		return
	}
	go func() {
		session, _ := concurrency.NewSession(etcdClient)
		m := concurrency.NewMutex(session, "testLock")
		err := m.Lock(context.TODO())
		if err != nil {
			fmt.Printf("lock fail, %v \n", err)
			return
		}
		fmt.Println("get lock 1 success")
		time.Sleep(60 * time.Second)
		defer m.Unlock(context.TODO())
	}()
	go func() {
		session, _ := concurrency.NewSession(etcdClient)
		m := concurrency.NewMutex(session, "testLock")
		err := m.Lock(context.TODO())
		if err != nil {
			fmt.Printf("lock fail, %v \n", err)
			return
		}
		fmt.Println("get lock 2 success")
		time.Sleep(60 * time.Second)
		defer m.Unlock(context.TODO())
	}()
	go func() {
		session, _ := concurrency.NewSession(etcdClient)
		m := concurrency.NewMutex(session, "testLock")
		err := m.Lock(context.TODO())
		if err != nil {
			fmt.Printf("lock fail, %v \n", err)
			return
		}
		fmt.Println("get lock 3 success")
		time.Sleep(60 * time.Second)
		defer m.Unlock(context.TODO())
	}()
	live := make(chan struct{})
	<-live
}
