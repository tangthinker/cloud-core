package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 填充数据
	blockSize := block.BlockSize()
	data = PKCS7Padding(data, blockSize)

	// 生成随机 IV
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 加密
	ciphertext := make([]byte, blockSize+len(data))
	copy(ciphertext[:blockSize], iv) // 前 blockSize 字节存储 IV
	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(ciphertext[blockSize:], data)

	// 返回 base64 编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := 0; i < padding; i++ {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}
