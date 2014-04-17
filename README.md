# julius

julius is a simple command line utility to encrypt and decrypt messages using the [Caesar cipher][caesar-wiki]. It is written in Go and makes use of the [caesar.go package][caesar.go]. 

julius applies the Caesar cipher to all letters in the English alphabet. Non-English characters **will be** preserved in encrypted output however.

[caesar-wiki]: https://en.wikipedia.org/wiki/Caesar_cipher
[caesar.go]: https://github.com/alexjohnj/caesar

## Installation

With the Go development tools installed and your `$GOPATH` set up, just run these commands:

```bash
go get github.com/alexjohnj/julius
go install github.com/alexjohnj/julius
```

## Usage

### Encryption

To encrypt a message using the default key of 13 (equivalent to applying a ROT13) you'd use the `encrypt` subcommand:

```bash
julius encrypt "Romani ite domum"
```

If you wanted to use a different key, pass the `--key`/`-k` flag:

```bash
julius encrypt --key=10 "Romani ite domum"
```

### Decryption

To decrypt a message use the `decrypt` subcommand. Again, using the default key of 13 you'd do the following:

```bash
julius decrypt "Ebznav vgr qbzhz"
```

Or for a custom key:

```bash
julius decrypt --key=10 "Bywkxs sdo nywew"
```

### Piped Input

julius accepts piped input and writes all of its text to stdout so these sort of commands are possibe:

```bash
cat secret-message.txt | julius encrypt --key=9> encrypted-secret-message.txt
```

### Headers

You can include a [PGP style][pgp-header] header in encrypted output using the `--include-header`/`-b` flag. julius will automatically strip the header when decrypting a message.

[pgp-header]: http://xkcd.com/1181/

## Acknowledgements

julius makes use of the [cli.go][cli.go] package by [codegangsta][codegangsta-profile].

[cli.go]: https://github.com/codegangsta/cli
[codegangsta-profile]: https://github.com/codegangsta

