/////////
//PRUEBAS
/////////

package main

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

//Genera una clave pública y otra privada
func generarClavesRSA() ([]byte, []byte) {
	claveprivada, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println(err.Error)
	}

	clavepublica := &claveprivada.PublicKey
	pemblock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(claveprivada)}

	fd, _ := x509.MarshalPKIXPublicKey(clavepublica)
	pemblockPublica := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: fd}

	return pemblock.Bytes, pemblockPublica.Bytes
}

//Cifrar con AES en modo CTR
func cifrarAES(textocifrar []byte, clave []byte) ([]byte, bool) {

	//Calculamos block con clave
	block, err := aes.NewCipher(clave)
	if err != nil {
		fmt.Println(err)
		return []byte{}, true
	}

	// IV necesita ser único aunque no seguro, se incluye al principio del textocifrado
	ciphertext := make([]byte, aes.BlockSize+len(textocifrar))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, true
	}

	//Ciframos
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], textocifrar)

	return ciphertext, false
}

//Con CTR se descifra como se cifra, con NewCTR
func descifrarAES(ciphertext []byte, clave []byte) ([]byte, bool) {

	//Calculamos block con clave
	block, err := aes.NewCipher(clave)
	if err != nil {
		fmt.Println(err)
		return []byte{}, true
	}

	//Volvemos a calcular iv (ahora sin rand, iv está al principio del textocifrado)
	iv := ciphertext[:aes.BlockSize]

	//Desciframos
	textodescifrado := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(textodescifrado, ciphertext[aes.BlockSize:])

	return textodescifrado, false
}

/*
//Cifrar con AES en modo CTR
func cifrarAES(key []byte, text []byte) ([]byte, bool) {
	// key := []byte(keyText)
	plaintext := text

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return []byte{}, true
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return []byte{}, true
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return ciphertext, false
}

//Con CTR se descifra como se cifra, con NewCTR
func descifrarAES(key []byte, cryptoText []byte) ([]byte, bool) {
	ciphertext := cryptoText

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return []byte{}, true
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		fmt.Println("Ciphertext too short")
		return []byte{}, true
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, false
}*/

//Realizamos firma digital con hash = SHA2 y cifrado clave privada de RSA
func firmaDigital(mensaje []byte, claveprivada []byte) ([]byte, bool) {

	//Hash con SHA-2
	hashmensaje := sha256.Sum256(mensaje)

	privateKey, err := x509.ParsePKCS1PrivateKey(claveprivada)
	if err != nil {
		fmt.Println(err)
		return []byte{}, false
	}

	firma, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashmensaje[:])
	if err != nil {
		fmt.Println(err.Error)
		return []byte{}, false
	}

	return firma, true
}

//Verificamos firma digital con hash = SHA2 y descifrado clave pública de RSA

func main() {

	fmt.Println("Hola")

	uno, _ := generarClavesRSA()

	fmt.Println("Mira la clave:")
	fmt.Println("Uno:::::: ", uno)

	unocifrada, _ := cifrarAES(uno, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	//fmt.Println("Mira la clave cifrada:")
	//fmt.Println("Dos:::::: ", unocifrada)

	unodescifrada, _ := descifrarAES(unocifrada, []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1})
	fmt.Println("Mira la clave descifrada:")
	fmt.Println("Tre:::::: ", unodescifrada)

	unodescifrada = bytes.Trim(unodescifrada, "\x00")
	fmt.Println("Mira la clave probando sin 0:")
	fmt.Println("Cua:::::: ", unodescifrada)

	if bytes.Compare(uno, unodescifrada) == 0 {
		fmt.Println("SON IGUALEEEEEEEEEEEEEEEEEEEEEEEES")
	} else {
		fmt.Println("NOOO SON IGUALEEEEEEEEEEEEEEEEEEEEEEEES")
	}

	firmica, _ := firmaDigital([]byte("Holaaaa"), unodescifrada)

	fmt.Println("VEMOOOOOOOOOOOS:::::: ", firmica)

}
