/*Hola Encarna.

El modo GCM es un tanto especial (para una cosa que se llama cifrado autentificado), los más comunes (y que te recomiendo) son CTR y OFB.
Tienes razón en que requieren nonces o IV (valores de inicialización). Lo normal es generarlos de forma aleatoria y concatenarlos antes de lo cifrado, tienen que ser aleatorios y distintos cada vez pero no secretos (la seguridad o secreto reside en la clave no en el nonce/iv).

En https://golang.org/pkg/crypto/cipher/ tienes un ejemplo más claro en la parte de CTR (te lo incluyo a continuación):
*/
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func main() {
	key := []byte("example key 1234")
	plaintext := []byte("some plaintext")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// *** RAFA: el IV se pone antes de lo cifrado y se separa al descifrar. ***
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	// CTR mode is the same for both encryption and decryption, so we can
	// also decrypt that ciphertext with NewCTR.

	plaintext2 := make([]byte, len(plaintext))
	stream = cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext2, ciphertext[aes.BlockSize:])

	fmt.Printf("%s\n", plaintext2)
}
