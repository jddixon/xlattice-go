package reg

// xlattice_go/reg/aes_cnx.go

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xt "github.com/jddixon/xlattice_go/transport"
	// "sync"
)

var _ = fmt.Print

const (
	MSG_BUF_LEN = 16 * 1024
)

type AesCnxHandler struct {
	State                              int
	Cnx                                *xt.TcpConnection
	engine                             cipher.Block
	encrypter                          cipher.BlockMode
	decrypter                          cipher.BlockMode
	iv1, key1, iv2, key2, salt1, salt2 []byte
}

func (h *AesCnxHandler) SetUpSessionKey() (err error) {
	h.engine, err = aes.NewCipher(h.key2)
	if err == nil {
		h.encrypter = cipher.NewCBCEncrypter(h.engine, h.iv2)
		h.decrypter = cipher.NewCBCDecrypter(h.engine, h.iv2)
	}
	return
}

// Read data from the connection.  XXX This will not handle partial
// reads correctly.
//
func (h *AesCnxHandler) ReadData() (data []byte, err error) {
	data = make([]byte, MSG_BUF_LEN)
	count, err := h.Cnx.Read(data)
	// DEBUG
	//fmt.Printf("ReadData: count is %d, err is %v\n", count, err)
	// END
	if err == nil && count > 0 {
		data = data[:count]
		return
	}
	return nil, err
}

func (h *AesCnxHandler) WriteData(data []byte) (err error) {
	count, err := h.Cnx.Write(data)
	// XXX handle cases where not all bytes written
	_ = count
	return
}
func DecodePacket(data []byte) (*XLRegMsg, error) {
	var m XLRegMsg
	err := proto.Unmarshal(data, &m)
	// XXX do some filtering, eg for nil op
	return &m, err
}

func EncodePacket(msg *XLRegMsg) (data []byte, err error) {
	return proto.Marshal(msg)
}

func EncodePadEncrypt(msg *XLRegMsg, engine cipher.BlockMode) (
	ciphertext []byte, err error) {

	var paddedData []byte

	cData, err := EncodePacket(msg)
	if err == nil {
		paddedData, err = xc.AddPKCS7Padding(cData, aes.BlockSize)
	}
	if err == nil {
		msgLen := len(paddedData)
		nBlocks := (msgLen + aes.BlockSize - 2) / aes.BlockSize
		ciphertext = make([]byte, nBlocks*aes.BlockSize)
		engine.CryptBlocks(ciphertext, paddedData) // dest <- src
	}
	return
}

func DecryptUnpadDecode(ciphertext []byte, engine cipher.BlockMode) (
	msg *XLRegMsg, err error) {

	plaintext := make([]byte, len(ciphertext))
	engine.CryptBlocks(plaintext, ciphertext) // dest <- src

	unpaddedCData, err := xc.StripPKCS7Padding(plaintext, aes.BlockSize)
	if err == nil {
		msg, err = DecodePacket(unpaddedCData)
	}
	return
}

//// Read the next message over the connection
//func (mc *Client) readMsg() (m *XLRegMsg, err error) {
//	inBuf, err := mc.h.ReadData()
//	if err == nil && inBuf != nil {
//		m, err = DecryptUnpadDecode(inBuf, mc.decrypter)
//	}
//	return
//}
//
//// Write a message out over the connection
//func (mc *Client) writeMsg(m *XLRegMsg) (err error) {
//	var data []byte
//	// serialize, marshal the message
//	data, err = EncodePadEncrypt(m, mc.encrypter)
//	if err == nil {
//		err = mc.h.WriteData(data)
//	}
//	return
//}
