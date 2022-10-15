package jnet

import (
	"Jinx/jinterface"
	"Jinx/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

// DataPack 封包拆包类
type DataPack struct {
}

// NewDataPack 初始化方法
func NewDataPack() *DataPack {
	dp := &DataPack{}
	return dp
}

// GetHeadLen 获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// Pack 封包方法
func (dp *DataPack) Pack(msg jinterface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 注意写入的顺序
	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		fmt.Println("Pack error: ", err)
		return nil, err
	}

	// 将msgID写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		fmt.Println("Pack error: ", err)
		return nil, err
	}

	// 将data写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		fmt.Println("Pack error: ", err)
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法 (将包的head信息读出来,再根据head的长度再读取data数据)
func (dp *DataPack) Unpack(binaryData []byte) (jinterface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息,得到dataLen和msgID
	msg := &Message{}

	// 读dataLen, 给msg.DataLen赋值
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		fmt.Println("Unpack error: ", err)
		return nil, err
	}

	// 读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		fmt.Println("Unpack error: ", err)
		return nil, err
	}

	// 判断dataLen是否已经超出了我们允许的最大包长度(>0是因为前提条件是设置了最大包长度)
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, fmt.Errorf("too large msg data received")
	}

	return msg, nil
}
