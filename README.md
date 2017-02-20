[![Travis CI](https://api.travis-ci.org/evantbyrne/cryptdir.svg?branch=master)](https://travis-ci.org/evantbyrne/cryptdir) [![Go report card](https://goreportcard.com/badge/github.com/evantbyrne/cryptdir)](https://goreportcard.com/report/github.com/evantbyrne/cryptdir)

# cryptdir

Utility for managing folders containing files encrypted with AES-256 GCM and Scrypt.


## Install

Make sure that [Go](https://golang.org/) is installed and your PATH includes GOBIN. Then run the following:

    $ git clone --recursive https://github.com/evantbyrne/cryptdir.git $GOPATH"/src/github.com/evantbyrne/cryptdir"
    $ go install github.com/evantbyrne/cryptdir

**Note:** The `$` at the beginning of newlines in this document represents the bash shell prompt, and is not a part of the actual commands. Likewise, `cryptdir>` represents an unlocked shell prompt.


## Usage

    $ cryptdir
    usage: cryptdir [<flags>] <command> [<args> ...]

    Utility for managing folders containing files encrypted with AES-256 GCM and
    Scrypt.

    Flags:
      --help  Show context-sensitive help (also try --help-long and --help-man).

    Commands:
      help [<command>...]
        Show help.

      ls [<flags>]
        List encrypted files.

      read <read_name>
        Read encrypted file.

      unlock
        Unlock encrypted directory.

      write <write_name>
        Write encrypted file.

Open a new shell with the encrypted folder unlocked:

    $ cd path/to/my/folder
    $ cryptdir unlock
    Password:
    cryptdir>

Read and write encrypted files from an unlocked shell, lock shell by exiting:

    cryptdir> cryptdir read hello.txt
    File not found.
    cryptdir> echo "Hello, world" | cryptdir write hello.txt
    cryptdir> cryptdir read hello.txt
    Hello, world
    cryptdir> exit
    exit
    $ cryptdir read hello.txt
    2017/02/05 13:42:27 The encrypted directory is locked. Please run `cryptdir unlock` to unlock.

Listing encrypted files from an unlocked shell:

    cryptdir> cryptdir ls
    foo.txt
    hello.txt
    zebra.png
    cryptdir> cryptdir ls -ms
    ZzvBpUZDpXVJmbLi foo.txt
    nUyWajppDtwLrLxj hello.txt
    xShsSwNrGNmnFyeC zebra.png


## Encrypted data format

Each file is encrypted in the following format, with a random salt and nonce generated per file:

    +----------------+---------------------+-------------------+-------------------+
    | 4 byte version | 12 byte scrypt salt | 12 byte gcm nonce | encrypted data... |
    +----------------+---------------------+-------------------+-------------------+

Here is an example of what the contents of an encrypted folder might look like:

![Example encrypted folder](https://raw.githubusercontent.com/evantbyrne/cryptdir/master/cryptdir-folder.png)

The raw file names are randomly generated 250 character strings. The `.cryptdir` file contains an encrypted mapping of human-readable filenames to the randomly named ones used on a filesystem level.
