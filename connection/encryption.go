package connection

import (
	"encoding/base64"
	"math/rand"
	"time"
)

type Encryption struct {
	data   []byte
	binkey []byte
	key    string
}

func NewEncryption(base64key string) *Encryption {
	e := new(Encryption)
	e.data = make([]byte, 1)
	if !e.SetKey(base64key) {
		e.key = ""
		e.binkey = make([]byte, 1)
	}

	return e
}

func (e *Encryption) SetKey(base64key string) bool {
	e.key = base64key
	var ok error
	e.binkey, ok = base64.StdEncoding.DecodeString(base64key)
	if ok != nil {
		return false
	}
	return true
}

func (e *Encryption) Key() string {
	return e.key
}

func GenerateKey(len int64) string {
	k := make([]byte, len)
	rand.Seed(time.Now().Unix())
	var i int64
	for i = 0; i < len; i++ {
		k[i] = byte(rand.Intn(256))
	}
	return base64.StdEncoding.EncodeToString(k)
}

func (e *Encryption) xor_crypt(buffer *[]byte) {
	key := e.binkey
	buf := *buffer
	j := 0
	for i := 0; i < len(buf); i++ {
		buf[i] = buf[i] ^ key[j]
		j++
		if j >= len(key) {
			j = 0
		}
	}
	return
}

func (e *Encryption) Encrypt(msg string) string {
	msgBuffer := []byte(msg)
	e.xor_crypt(&msgBuffer)
	return base64.StdEncoding.EncodeToString(msgBuffer)
}

func (e *Encryption) Decrypt(emsg string) string {
    msgBuffer, ok := base64.StdEncoding.DecodeString(emsg)
    if ok != nil {
        return ""
    }
    e.xor_crypt(&msgBuffer)
    return string(msgBuffer)
}
