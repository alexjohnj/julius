package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "julius"
	app.Version = "0.1.0"
	app.Usage = "Encrypt and decrypt ASCII strings using a Caesar cipher"
	app.Commands = []cli.Command{
		{
			Name:        "encrypt",
			ShortName:   "e",
			Usage:       "julius encrypt [options] [message]",
			Description: "Encrypts a message using a key given with the -k flag. Defaults to 13 if no key is given.",
			Action:      encryptMessage,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "key, k",
					Value: 13,
					Usage: "The key to use for the cipher.",
				},
			},
		},

		{
			Name:        "decrypt",
			ShortName:   "d",
			Usage:       "julius decrypt [options] [message]",
			Description: "Decrypts a message using a key given with the -k flag. Defaults to 13 if no key is given",
			Action:      decryptMessage,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "key, k",
					Value: 13,
					Usage: "The key used to decrypt the message.",
				},
			},
		},
	}

	app.Run(os.Args)
}

func encryptMessage(c *cli.Context) {

}

func decryptMessage(c *cli.Context) {

}
