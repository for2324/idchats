package order

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

/*
*
定义两个指针：虚拟指针和顺序指针

一开始使用顺序指针，但是发现顺序指针会导致链表的顺序发生变化，所以改用虚拟指针
*/
type Dummy struct {
	Head int `json:"head"`
	Size int `json:"size"`
}

type Order struct {
	Head int `json:"head"`
	Size int `json:"size"`
}

// 定义链表元素结构体
type Element struct {
	Data int `json:"data"`
	Next int `json:"next"`
}

const (
	CHAIN_TABLE_KEY = "pay_nonce_chain_table"
)

var (
	ErrNoElement = errors.New("no element")
	// 重复插入元素
	ErrElementExist           = errors.New("element exist")
	ErrCasVersionInconformity = errors.New("cas compare version inconformity")
)

type ChainTable struct {
	NameSpace            string
	SpaceSize            int
	ConcurrencyPrecision int
	LeaseDelay           time.Duration
}

type ChainTableElement struct {
	OrderCursor int
	RentList    []string
	BackList    []int
}

func (c *ChainTable) WatchDelay() {
	cli, err := GetEtcdClient()
	if err != nil {
		return
	}
	// 监听前缀为"/my-prefix/"的key的变化
	watchChan := cli.Watch(context.Background(), fmt.Sprintf("%s/%s/", c.GetPrefix(), "delay"), clientv3.WithPrefix())
	// 处理key的变化
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if event.Type == clientv3.EventTypeDelete {
				log.NewInfo("WatchDelay", c.GetPrefix()+"Key has been deleted", string(event.Kv.Key))

				list := strings.Split(string(event.Kv.Key), ":")
				if len(list) == 3 {
					data, err := strconv.Atoi(list[2])
					if err == nil {
						c.giveBack(data)
					}
				}
			}
		}
	}
}

func NewChainTable(nameSpace string, concurrencyPrecision int) *ChainTable {
	LeaseDelay := time.Duration(config.Config.Pay.OrderExpireTime) * time.Minute
	// 延迟1分钟归还，防止不可预知的情况
	LeaseDelay = LeaseDelay + time.Minute*1
	spaceSize := int(math.Pow10(concurrencyPrecision))
	c := &ChainTable{
		NameSpace:            nameSpace,
		SpaceSize:            spaceSize,
		LeaseDelay:           LeaseDelay,
		ConcurrencyPrecision: concurrencyPrecision,
	}
	c.FilterRentDataExpired()
	go c.WatchDelay()
	return c
}

func (c *ChainTable) FilterRentDataExpired() {
	cli, err := GetEtcdClient()
	if err != nil {
		return
	}
	defer cli.Close()
	chainTableElement, _, err := c.GetEtcdChainTableElement(cli)
	if err != nil {
		return
	}
	for _, v := range chainTableElement.RentList {
		list := strings.Split(v, ":")
		if len(list) != 3 {
			continue
		}
		val := list[2]
		// 判断是否存在
		resp, err := cli.Get(context.Background(), v)
		if err != nil {
			continue
		}
		if len(resp.Kvs) == 0 {
			val, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			c.giveBack(val)
			continue
		}
	}
}

func (c *ChainTable) GetPrefix() string {
	return fmt.Sprintf("%s:///%s/", CHAIN_TABLE_KEY, c.NameSpace)
}

func (c *ChainTable) GiveBack(giveBackId string) (err error) {
	etcdClient, err := GetEtcdClient()
	if err != nil {
		return err
	}
	// 删除key
	_, err = etcdClient.Delete(context.Background(), giveBackId)
	return err
}

func (c *ChainTable) GetEtcdChainTableElement(cli *clientv3.Client) (chainTableElement *ChainTableElement, curVersion int64, err error) {
	// 获取链表头元素
	var resp *clientv3.GetResponse
	resp, err = cli.Get(context.Background(), c.GetPrefix())
	if err != nil {
		return
	}
	if len(resp.Kvs) == 0 {
		err = fmt.Errorf("list is empty")
		return
	}
	value := resp.Kvs[0].Value
	curVersion = resp.Kvs[0].Version
	err = json.Unmarshal(value, &chainTableElement)
	return
}

// 归还元素
func (c *ChainTable) giveBack(data int) (err error) {
	var cli *clientv3.Client
	cli, err = GetEtcdClient()
	if err != nil {
		return
	}
	defer cli.Close()
	chainTableElement, curVersion, err := c.GetEtcdChainTableElement(cli)
	if err != nil {
		return err
	}
	// 判断元素是否存在（防止重复归还）
	for _, v := range chainTableElement.BackList {
		if v == data {
			return ErrElementExist
		}
	}
	RentList := []string{}
	for _, v := range chainTableElement.RentList {
		list := strings.Split(v, ":")
		if len(list) != 3 {
			continue
		}
		val, err := strconv.Atoi(list[2])
		if err != nil {
			continue
		}
		if val == data {
			continue
		}
		RentList = append(RentList, v)
	}
	chainTableElement.RentList = RentList
	chainTableElement.BackList = append(chainTableElement.BackList, data)
	// 更新链表头元素
	chainTableElementBytes, err := json.Marshal(chainTableElement)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	txnResp, err := cli.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(c.GetPrefix()), "=", curVersion)).
		Then(clientv3.OpPut(c.GetPrefix(), string(chainTableElementBytes))).
		Commit()
	log.NewInfo("give back", "set new value txnResp version", txnResp.Header.Revision)
	return err
}

func (c *ChainTable) Next() (nextValue int, giveBackId string, Rerr error) {
	etcdClient, err := GetEtcdClient()
	if err != nil {
		return 0, "", err
	}
	// 获取链表头元素
	resp, err := etcdClient.Get(context.Background(), c.GetPrefix())
	if err != nil {
		return 0, "", err
	}
	var chainTableElement ChainTableElement
	var curVersion int64 = 0
	defer func() {
		if Rerr == nil {
			// 获取随机数 string
			randStr := strconv.Itoa(rand.Int())
			giveBackId = fmt.Sprintf("%s/%s/%s:%d", c.GetPrefix(), "delay", randStr, nextValue)
			// 创建租约
			err := c.SetDelayOperation(giveBackId, nextValue)
			if err != nil {
				Rerr = err
				return
			}
			chainTableElement.RentList = append(chainTableElement.RentList, giveBackId)
			// 更新链表头元素
			chainTableElementBytes, err := json.Marshal(chainTableElement)
			if err != nil {
				Rerr = err
				log.NewError("next", "json marshal error", err.Error())
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			var txnResp *clientv3.TxnResponse
			txnResp, Rerr = etcdClient.Txn(ctx).
				If(clientv3.Compare(clientv3.Version(c.GetPrefix()), "=", curVersion)).
				Then(clientv3.OpPut(c.GetPrefix(), string(chainTableElementBytes))).
				Commit()
			if Rerr != nil {
				log.NewError("next", "set new value error", Rerr.Error())
				return
			}
			if !txnResp.Succeeded {
				Rerr = ErrCasVersionInconformity
				return
			}
			// 翻转 nextValue
			nextValue = reverseNumber(nextValue, c.ConcurrencyPrecision)
		}
	}()
	if len(resp.Kvs) == 0 {
		chainTableElement = ChainTableElement{
			OrderCursor: 1,
			BackList:    []int{},
		}
		return 1, "", nil
	}
	value := resp.Kvs[0].Value
	curVersion = resp.Kvs[0].Version

	err = json.Unmarshal(value, &chainTableElement)
	if err != nil {
		return 0, "", err
	}
	// 直接使用池子里的元素
	if len(chainTableElement.BackList) > 0 {
		nextValue = chainTableElement.BackList[len(chainTableElement.BackList)-1]
		chainTableElement.BackList = chainTableElement.BackList[:len(chainTableElement.BackList)-1]
		return nextValue, "", nil
	}
	if chainTableElement.OrderCursor < c.SpaceSize-1 {
		chainTableElement.OrderCursor++
		return chainTableElement.OrderCursor, "", nil
	}
	return 0, "", ErrNoElement
}

func GetEtcdClient() (*clientv3.Client, error) {
	etcdAddr := strings.Join(config.Config.Etcd.EtcdAddr, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcdAddr, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c *ChainTable) SetDelayOperation(key string, value int) error {
	// 获取 etcd 客户端
	client, err := GetEtcdClient()
	if err != nil {
		return err
	}
	defer client.Close()
	// 创建租约
	resp, err := client.Grant(context.Background(), int64(c.LeaseDelay/time.Second))
	if err != nil {
		return err
	}
	// 将键与租约关联
	_, err = client.Put(context.Background(), key, fmt.Sprint(value), clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	return nil
}

func reverseNumber(n int, numCount int) int {
	// Convert the number to a string
	s := strconv.Itoa(n)

	// Create an array of length 5 with all elements set to '0'
	arr := []byte{}
	for i := 0; i < numCount; i++ {
		arr = append(arr, '0')
	}

	// Iterate over the string and set the corresponding array elements
	for i := 0; i < len(s); i++ {
		arr[i] = s[len(s)-1-i]
	}
	nStr := string(arr[:])
	// Convert the array to a string
	newN, _ := strconv.ParseInt(nStr, 10, 64)
	return int(newN)
}
