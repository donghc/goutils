package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

var (
	config  clientv3.Config
	err     error
	client  *clientv3.Client
	kv      clientv3.KV
	putResp *clientv3.PutResponse
	path    = "/config/engines/v1"
)

func t2() {

	//获取时间戳

	timestamp := time.Now().Unix()
	fmt.Println(timestamp)
	//格式化为字符串,tm为Time类型
	tm := time.Unix(timestamp, 0)
	fmt.Println(tm.Format("2006-01-02 15:04:05"))
}

func main() {
	t2()

}

func t1() {
	//配置
	config = clientv3.Config{
		Endpoints:   []string{"192.168.102.155:2379"},
		Username:    "app",
		Password:    "psKdqSe",
		DialTimeout: time.Second * 5,
	}
	//连接 床见一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("clientv3.New : ", err)
		return
	}
	//用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	put(kv)
	get(kv)
}

func get(kv clientv3.KV) {
	resp, err := kv.Get(context.TODO(), path)
	if err != nil {
		fmt.Println(" kv.Get : ", err)
		return
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			fmt.Println(string(v))
		}
	}
	if resp.Count == 0 {
		fmt.Println("没有获取到数据")
		return
	}
	fmt.Println(string(resp.Kvs[0].Value))
}

func put(kv clientv3.KV) {
	en := make(map[string]bool)
	en["en1"] = true
	en["en2"] = true
	en["en3"] = true
	en["en4"] = true
	en["en5"] = true
	en["en6"] = true
	marshal, _ := json.Marshal(en)
	putResp, err = kv.Put(context.TODO(), path, string(marshal), clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(" kv.Put : ", err)
	} else {
		//获取版本信息
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("key:", string(putResp.PrevKv.Key))
			fmt.Println("Value:", string(putResp.PrevKv.Value))
			fmt.Println("Version:", string(putResp.PrevKv.Version))
		}
	}
}
