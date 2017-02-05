package pa

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Files map[string]string `json:"files"`
}

const currentVersion uint32 = 1
const dirMode = 0700
const fileMode = 0600

func configJoin(ciphertext, nonce, salt []byte) []byte {
	out := make([]byte, 4)
	binary.LittleEndian.PutUint32(out, currentVersion)

	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	return out
}

func configSplit(data, password []byte) (ciphertext, nonce, key, salt []byte, err error) {
	if len(data) < 25 {
		return nil, nil, nil, nil, errors.New("Invalid config data length.")
	}

	// version = data[:4] uint32 reserved for later use.
	salt = data[4:16]
	nonce = data[16:28]
	ciphertext = data[28:]
	key = NewKey(salt, password)

	return ciphertext, nonce, key, salt, nil
}

func ConfigRead(configDir string, password []byte) (config Config, nonce, key, salt []byte) {
	var (
		err       error
		cleartext []byte
	)

	cleartext, nonce, key, salt = DecryptFile(configDir+"/.cryptdir", password)

	if string(cleartext) != "" {

		if err = json.Unmarshal(cleartext, &config); err != nil {
			fmt.Println("Could not parse ctyptdir config file JSON.")
			os.Exit(1)
		}

		return config, nonce, key, salt
	}

	return Config{Files: make(map[string]string)}, nonce, key, salt
}

func ConfigWrite(configDir string, config Config, nonce, key, salt []byte) {
	var (
		configMarshalled []byte
		err              error
	)

	if configMarshalled, err = json.MarshalIndent(config, "", "\t"); err != nil {
		fmt.Printf("Could not convert config to JSON format.")
		os.Exit(1)
	}

	EncryptFile(configDir+"/.cryptdir", configMarshalled, nonce, key, salt)
}

func DecryptFile(path string, password []byte) (cleartext, nonce, key, salt []byte) {

	if _, err := os.Stat(path); err == nil {

		var (
			ciphertext []byte
			cleartext  []byte
			configData []byte
		)

		if configData, err = ioutil.ReadFile(path); err != nil {
			fmt.Printf("Could not read file at '%s'.\n", path)
			os.Exit(1)
		}

		if ciphertext, nonce, key, salt, err = configSplit(configData, password); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if cleartext, err = Decrypt(ciphertext, nonce, key); err != nil {
			fmt.Printf("Could not decrypt '%s' with given password.\n", path)
			os.Exit(1)
		}

		return cleartext, nonce, key, salt
	}

	nonce = NewNonce()
	salt = NewNonce()
	key = NewKey(salt, password)

	return cleartext, nonce, key, salt
}

func EncryptFile(path string, cleartext, nonce, key, salt []byte) {
	var (
		ciphertext []byte
		err        error
	)

	if ciphertext, err = Encrypt(cleartext, nonce, key); err != nil {
		fmt.Println("Unexpected error: could not encrypt password vault.")
		os.Exit(1)
	}

	if err := ioutil.WriteFile(path, configJoin(ciphertext, nonce, salt), fileMode); err != nil {
		fmt.Printf("Could not write file at '%s'.", path)
		os.Exit(1)
	}
}
