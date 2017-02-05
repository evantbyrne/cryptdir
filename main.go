package main

import (
	"os"

	"github.com/alecthomas/kingpin"

	cryptdir "github.com/evantbyrne/cryptdir/lib"
)

var (
	apiVersion = 1

	app = kingpin.New("cryptdir", "Utility for managing folders containing files encrypted with AES-256 GCM and Scrypt.")

	appRead         = app.Command("read", "Read encrypted file.")
	appReadFileName = appRead.Arg("read_name", "File name.").Required().String()

	appUnlock = app.Command("unlock", "Unlock encrypted directory.")

	appWrite         = app.Command("write", "Write encrypted file.")
	appWriteFileName = appWrite.Arg("write_name", "File name.").Required().String()
)

func main() {
	var (
		kp = kingpin.MustParse(app.Parse(os.Args[1:]))
	)

	switch kp {

	case appRead.FullCommand():
		cryptdir.CommandRead(*appReadFileName)

	case appUnlock.FullCommand():
		cryptdir.CommandUnlock()

	case appWrite.FullCommand():
		cryptdir.CommandWrite(*appWriteFileName)

	}
}
