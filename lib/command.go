package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	envPassword       = "CRYPTDIR_PASSWORD"
	messageLocked     = "The encrypted directory is locked. Please run `cryptdir unlock` to unlock."
	messageWorkingDir = "Unexpected error: could not get working directory."
	rawFileNameLength = 250
	shellPrompt       = "cryptdir> "
)

func CommandList(mapped, short bool) {
	var (
		config    Config
		configDir string
		password  []byte
	)

	password = []byte(mustGetEnvPassword())
	configDir = mustGetWorkingDir()
	config, _, _, _ = ConfigRead(configDir, password)

	for fileName, rawFileName := range config.Files {
		if mapped && short {
			fmt.Printf("%s %s\n", rawFileName[0:16], fileName)
		} else if mapped {
			fmt.Printf("%s %s\n", rawFileName, fileName)
		} else {
			fmt.Println(fileName)
		}
	}
}

func CommandRead(fileName string) {
	var (
		cleartext []byte
		config    Config
		configDir string
		password  []byte
	)

	password = []byte(mustGetEnvPassword())
	configDir = mustGetWorkingDir()
	config, _, _, _ = ConfigRead(configDir, password)

	if _, ok := config.Files[fileName]; !ok {
		fmt.Println("File not found.")
		os.Exit(1)
	}

	cleartext, _, _, _ = DecryptFile(configDir+"/"+config.Files[fileName], password)
	os.Stdout.Write(cleartext)
}

func CommandUnlock() {
	var (
		configDir string
		err       error
		password  []byte
	)

	fmt.Print("Password: ")
	if password, err = terminal.ReadPassword(int(syscall.Stdin)); err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n")

	configDir = mustGetWorkingDir()
	ConfigRead(configDir, password)

	os.Setenv(envPassword, string(password))
	os.Setenv("PS1", shellPrompt)
	syscall.Exec(os.Getenv("SHELL"), []string{os.Getenv("SHELL")}, syscall.Environ())
}

func CommandWrite(fileName string) {
	var (
		config      Config
		configDir   string
		configNonce []byte
		configKey   []byte
		configSalt  []byte
		content     []byte
		fileNonce   []byte
		fileKey     []byte
		fileSalt    []byte
		password    []byte
		rawFileName string
	)

	password = []byte(mustGetEnvPassword())
	content, _ = ioutil.ReadAll(os.Stdin)
	configDir = mustGetWorkingDir()
	config, configNonce, configKey, configSalt = ConfigRead(configDir, password)
	rawFileName = getRawFileName(config, fileName)
	config.Files[fileName] = rawFileName

	ConfigWrite(configDir, config, configNonce, configKey, configSalt)

	fileNonce = NewNonce()
	fileSalt = NewNonce()
	fileKey = NewKey(fileSalt, password)
	EncryptFile(configDir+"/"+rawFileName, content, fileNonce, fileKey, fileSalt)
}
