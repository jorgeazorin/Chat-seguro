package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	//"io/ioutil"
	"log"
)

func main() {

	var ciphertext, plaintext []byte
	var err error
	// The key length can be 32, 24, 16  bytes (OR in bits: 128, 192 or 256)
	key := []byte("longer means more possible keys ")
	plaintext = []byte("This is the unecrypted data. Referring to it as plain text.")

	if ciphertext, err = encrypt(key, plaintext); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%0x\n", ciphertext)

	if plaintext, err = decrypt(key, ciphertext); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", plaintext)
}

func encrypt(key, text []byte) (ciphertext []byte, err error) {

	var block cipher.Block

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	ciphertext = make([]byte, aes.BlockSize+len(string(text)))

	// iv =  initialization vector
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)

	return
}

func decrypt(key, ciphertext []byte) (plaintext []byte, err error) {

	var block cipher.Block

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	if len(ciphertext) < aes.BlockSize {
		err = errors.New("ciphertext too short")
		return
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	plaintext = ciphertext

	return
}
