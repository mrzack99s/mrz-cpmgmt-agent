package encryptions

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

func Decrypt(encryptedString string) string {
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher([]byte(AES_KEY))
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
