package jnet

import (
	"io"
	"net"
	"testing"
)

// TestNewDataPack 测试封包、拆包类的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	// 1 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Error("server listen err: ", err)
		return
	}

	// 创建一个go承载 负责从客户端处理业务
	go func() {
		// 2 从客户端读取数据，拆包处理
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Error("server accept err: ", err)
				return
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				// -------- 拆包的过程 --------
				// 定义一个拆包的对象
				dp := NewDataPack()
				for {
					// 1 先读出流中的head部分，得到ID和dataLen
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						t.Error("read head error: ", err)
						return
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						t.Error("server unpack err: ", err)
						return
					}

					if msgHead.GetMsgLen() > 0 {
						// Msg是有data数据的，需要再次读取data数据
						// 2 再根据dataLen进行第二次读取，将data读出来
						msg := msgHead.(*Message)                // 类型断言
						msg.Data = make([]byte, msg.GetMsgLen()) // 重新分配内存

						// 根据dataLen的长度，再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							t.Error("server unpack data err: ", err)
							return
						}

						// 完整的一个消息已经读取完毕
						t.Logf("==> Recv Msg: ID=%d, len=%d, data=%s", msg.Id, msg.DataLen, string(msg.Data))
					}

				}

			}(conn)
		}
	}()

	/*
		模拟的客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Error("client dial err: ", err)
		return
	}

	// 创建一个封包对象 dp
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一同发送
	// 封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l', 'o'},
	}

	// 封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
	}

	dp1, err := dp.Pack(msg1)
	if err != nil {
		t.Error("client pack msg1 err: ", err)
		return
	}
	dp2, err := dp.Pack(msg2)
	if err != nil {
		t.Error("client pack msg2 err: ", err)
		return
	}

	// 将两个包粘在一起
	dp1 = append(dp1, dp2...) // 打散dp2，将dp2的元素添加到dp1中,否则是一个嵌套的切片

	// 向服务器端写数据
	_, err = conn.Write(dp1)
	if err != nil {
		t.Error("client write err: ", err)
		return
	}

	// 客户端阻塞
	select {}
}
