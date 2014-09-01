package main

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt"
	"github.com/alexjohnj/caesar"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const encryptedMessageHeader = "-----BEGIN JULIUS MESSAGE-----\n\n"
const encryptedMessageFooter = "\n\n-----END JULIUS MESSAGE-----"

func main() {
	app := cli.NewApp()
	app.Name = "julius"
	app.Version = "0.1.1"
	app.Usage = "Encrypt and decrypt messages using the Caesar cipher."
	app.Commands = []cli.Command{
		{
			Name:        "encrypt",
			ShortName:   "e",
			Usage:       "julius encrypt [options] [message]",
			Description: "Encrypts a plaintext message. The default key is 13, use the --key flag to change it.",
			Action:      encryptMessage,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "key, k",
					Value: 13,
					Usage: "The key to use for the cipher.",
				},
				cli.BoolFlag{
					Name:  "include-header, b",
					Usage: "Include a PGP style header in the encrypted output.",
				},
			},
		},

		{
			Name:        "decrypt",
			ShortName:   "d",
			Usage:       "julius decrypt [options] [message]",
			Description: "Decrypts ciphertext. By default it uses a key of 13. use the --key flag to change it.",
			Action:      decryptMessage,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "key, k",
					Value: 13,
					Usage: "The key to use to decrypt the message.",
				},
			},
		},

		{
			Name:        "brute",
			ShortName:   "b",
			Usage:       "julius brute [message]",
			Description: "Brute forces the key for a ciphertext by trying all possibly keys.",
			Action:      bruteForceMessage,
		},
	}

	app.Run(os.Args)
}

/*----------------------------------------------
									CLI FUNCTIONS
-----------------------------------------------*/

func encryptMessage(c *cli.Context) {
	plaintext := getUserMessage(c)
	key := c.Int("key")

	ciphertext := caesar.EncryptPlaintext(plaintext, key)

	if c.Bool("include-header") {
		fmt.Printf("%s%s%s", encryptedMessageHeader, ciphertext, encryptedMessageFooter)
	} else {
		fmt.Printf("%s\n", ciphertext)
	}
}

func decryptMessage(c *cli.Context) {
	ciphertext := getUserMessage(c)
	ciphertext = stripJuliusHeader(c, ciphertext)
	key := c.Int("key")

	plaintext := caesar.DecryptCiphertext(ciphertext, key)

	fmt.Printf("%s\n", plaintext)
}

func bruteForceMessage(c *cli.Context) {
	ciphertext := getUserMessage(c)
	ciphertext = stripJuliusHeader(c, ciphertext)
	var plaintexts [26]string

	for key := 0; key < 26; key++ {
		plaintexts[key] = caesar.DecryptCiphertext(ciphertext, key)
	}

	for key := 0; key < 26; key++ {
		fmt.Printf("[Key: %d]: %s\n", key, plaintexts[key])
	}
}

/*----------------------------------------------
								HELPER FUNCTIONS
-----------------------------------------------*/

// getUserMessage tries to obtain the user's message from either the command arguments, piped stdin or by prompting the user for it.
// It returns the message as a string
func getUserMessage(c *cli.Context) string {
	var messageArgument string

	// Try to determine if the user provided a message as an argument, piped one in or just didn't bother
	if len(c.Args()) < 1 && !terminal.IsTerminal(int(os.Stdin.Fd())) {
		messageArgument = readFromFile(os.Stdin) // Read the piped input
	} else if len(c.Args()) < 1 && terminal.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Printf("Enter a message (CTRL+D to end entry):\n") // Prompt the user to enter something
		messageArgument = readFromFile(os.Stdin)
	} else {
		messageArgument = c.Args()[0] // The user passed some text as an argument
	}
	return messageArgument
}

// readFromFile reads a file line-by-line and returns its contents in a single string
func readFromFile(f *os.File) string {
	fileContent, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatal(err)
	}
	fileContent = bytes.TrimSuffix(fileContent, []byte("\n"))
	return string(fileContent)
}

// stripJuliusHeader returns a string with the standard julius header/footer text removed
func stripJuliusHeader(c *cli.Context, message string) string {
	message = strings.Replace(message, encryptedMessageHeader, "", 1)
	message = strings.Replace(message, encryptedMessageFooter, "", 1)

	return message
}
