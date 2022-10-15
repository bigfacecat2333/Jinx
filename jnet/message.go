package jnet

// Message TLV: Type, Length, Value 来解决tcp粘包问题
type Message struct {
	Id      uint32 // 消息的ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的内容
}

// NewMsgPackage 初始化一个Message
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}