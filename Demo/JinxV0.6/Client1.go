package main

import (
	"Jinx/jnet"
	"fmt"
	"io"
	"net"
	"time"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	// 1. 直接链接远程服务器，得到一个conn连接, Dial()就是C++中的Connect()的封装
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client1 start err, exit!")
		return
	}

	for {
		// 发送封包的message消息
		dp := jnet.NewDataPack()

		binaryMsg, err := dp.Pack(jnet.NewMsgPackage(1, []byte("Jinx client Test Message")))
		if err != nil {
			fmt.Println("Pack error err", err)
			return
		}

		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error err", err)
			return
		}

		// 服务器应该先回复一个MsgID:1的消息ping

		// 先读出流中的head部分 得到ID和dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, binaryHead)
		if err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		// 再根据head中的dataLen再次读取data，放在data中
		fmt.Println("binaryHead: ", string(binaryHead))
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack error ", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*jnet.Message)
			msg.SetData(make([]byte, msg.GetMsgLen()))

			// 根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.GetData())
			if err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		fmt.Println("==> Receive Server Msg: ID=", msgHead.GetMsgId(),
			", len=", msgHead.GetMsgLen(),
			", data=", string(msgHead.GetData()))

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
