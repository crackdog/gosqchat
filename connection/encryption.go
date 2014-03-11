package connection

import (
	"encoding/base64"
	"math/rand"
	"time"
    "io/ioutil"
)

//Encryption
type Encryption struct {
	binkey []byte
	key    string
}

//NewEncryption is the constructor for Encryption.
func NewEncryption(base64key string) *Encryption {
	e := new(Encryption)
	if !e.SetKey(base64key) {
		e.key = ""
		e.binkey = make([]byte, 1)
	}

	return e
}

//SetKey sets a new key for an Encryption object.
func (e *Encryption) SetKey(base64key string) bool {
	e.key = base64key
    if base64key == "" {
        e.binkey = make([]byte, 1)
        e.binkey[0] = 0
        return true
    }
	var err error
	e.binkey, err = base64.StdEncoding.DecodeString(base64key)
	if err != nil {
		return false
	}
	return true
}

//LoadKeyFile loads a key from a file and saves it in e.
func (e *Encryption) LoadKeyFile(filename string) bool {
    key, err := ioutil.ReadFile(filename)
    if err != nil {
        return false
    }
    return e.SetKey(string(key))
}

//SaveKeyFile saves the key from e in a file.
func (e *Encryption) SaveKeyFile(filename string) bool {
    data := []byte(e.Key())
    err := ioutil.WriteFile(filename, data, 0644)
    if err != nil {
        panic(err)
        return false
    }
    return true
}

//Key gets the key from an Encryption object.
func (e *Encryption) Key() string {
	return e.key
}

//GenerateKey generates a key for an Encryption object.
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

//Encryption encrypts a message string.
func (e *Encryption) Encrypt(msg string) string {
	msgBuffer := []byte(msg)
	e.xor_crypt(&msgBuffer)
	return base64.StdEncoding.EncodeToString(msgBuffer)
}

//Decrypt decrypts an encrypted string.
func (e *Encryption) Decrypt(emsg string) string {
	msgBuffer, ok := base64.StdEncoding.DecodeString(emsg)
	if ok != nil {
		return ""
	}
	e.xor_crypt(&msgBuffer)
	return string(msgBuffer)
}
