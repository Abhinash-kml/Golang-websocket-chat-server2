package models

type Message struct {
	Sender  string      `json:"user"`
	Type    string      `json:"type"`
	Channel string      `json:"channel"`
	Payload interface{} `json:"payload"`
}

type TextMessage struct {
	Data string `json:"data"`
}

type BinaryMessage struct {
	Data string `json:"binarydata"`
}

type ByteArray struct {
	Data []byte `json:"bytearray"`
}
