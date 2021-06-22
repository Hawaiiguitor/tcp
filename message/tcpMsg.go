package message

import "encoding/json"

const (
	MAX_NET_DATA_SIZE = 1024
	OP_SENDINFO       = 1 // send file info
	OP_SENDFILE       = 2 // send file content
)

type MsgHeader struct {
	DataSize int // the size of package: header + body
	OpCode   uint8
}

type TcpMsg struct {
	Header MsgHeader
	Body   []byte
}

type fileData struct {
	filename string
	buffer   []byte
}

func ConstructMsg(msg *TcpMsg) ([]byte, error) {
	endCode := byte('\n')
	hb, err := json.Marshal(msg.Header)
	if err != nil {
		return nil, err
	}
	var buf []byte
	buf = append(buf, hb...)
	buf = append(buf, endCode)
	buf = append(buf, msg.Body...) // header + body
	return buf, nil
}

func DecodeHeader(msg []byte) (MsgHeader, error) {
	header := MsgHeader{}
	err := json.Unmarshal(msg, &header)
	if err != nil {
		return header, err
	}
	return header, nil
}
