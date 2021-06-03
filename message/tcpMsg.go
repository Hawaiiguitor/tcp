package message


type msgHeader struct {
    datasize int
	end byte
}

type tcpMsg struct {
	header msgHeader
	body []byte
}