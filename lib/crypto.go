package pa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/golang/crypto/scrypt"
	"log"
	"math/big"
)

const (
	randomComplete = "`~^0OolI\"'/\\|"
	randomLetter   = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	randomNumber   = "123456789"
	randomSpecial  = "!@#$%&*()_-:;.,?+=<>[]{}"
)

func Decrypt(ciphertext, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, err
}

func Encrypt(cleartext, nonce, key []byte) ([]byte, error) {
	var block cipher.Block
	var ciphertext []byte
	var err error

	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext = aesgcm.Seal(nil, nonce, cleartext, nil)

	return ciphertext, err
}

func NewNonce() []byte {
	nonce := make([]byte, 12)
	rand.Read(nonce)
	return nonce
}

func NewKey(salt, password []byte) []byte {
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)

	if err != nil {
		log.Fatal(err)
	}

	return key
}

func RandomString(length int, complete, noNumber, noSpecial bool) string {
	var randomPool = randomLetter

	if complete {

		if noNumber || noSpecial {
			log.Fatal("Cannot use `complete` flag with `no-number` and `no-special`.")
		}

		randomPool += randomNumber
		randomPool += randomSpecial
		randomPool += randomComplete

	} else {

		if !noNumber {
			randomPool += randomNumber
		}

		if !noSpecial {
			randomPool += randomSpecial
		}
	}

	randstr := make([]byte, length) // Random string to return
	charlen := big.NewInt(int64(len(randomPool)))
	for i := 0; i < length; i++ {
		b, err := rand.Int(rand.Reader, charlen)
		if err != nil {
			log.Fatal(err)
		}
		r := int(b.Int64())
		randstr[i] = randomPool[r]
	}
	return string(randstr)
}
