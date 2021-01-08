package encryptions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(plantext string) string {

	block, err := aes.NewCipher([]byte(AES_KEY))
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plantext), nil)
	return fmt.Sprintf("%x", ciphertext)
}
