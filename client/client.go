package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/util"
	"zinx/znet"
)

func main() {
	util.Init()
	fmt.Println("Client Test xty start")
	time.Sleep(2 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start error, exit!")
		return
	}
	dp := znet.NewDataPack()

	var id uint32
	id = 0
	for id <= 1 {
		data := []byte("hello!!!")
		sendMsg, err := dp.Pack(znet.NewMsgPackage(id, data))
		if err != nil {
			fmt.Println("send msg error ", err)
			break
		}
		conn.Write(sendMsg)
		head := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, head)
		if err != nil {
			fmt.Println("read buf error ", err)
			break
		}
		msg, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("msg head unpack error ", err)
			break
		}
		if msg.GetDataLen() > 0 {
			dataBuf := make([]byte, msg.GetDataLen())
			cnt, err := io.ReadFull(conn, dataBuf)
			if err != nil {
				fmt.Println("get msg data error ", err)
				break
			}
			msg.SetData(dataBuf)
			fmt.Printf("received message %d: %s, cnt=%d\n", msg.GetMsgID(), msg.GetData(), cnt)
		}
		id++
		time.Sleep(1 * time.Second)
	}
}
