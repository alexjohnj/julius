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

// Source for letter frequency:
// http://en.wikipedia.org/wiki/Letter_frequency#Relative_frequencies_of_letters_in_the_English_language
const englishFrequencyList = "etaoinshrdlcumwfgypbvkjxqz"

type Message struct {
	key        int
	plaintext  string
	ciphertext string
}

func main() {
	app := cli.NewApp()
	app.Name = "julius"
	app.Version = "0.2.0--dev"
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
			Usage:       "julius brute [options] [message]",
			Description: "Brute forces the key for a ciphertext by trying all possibly keys.",
			Action:      bruteForceMessage,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "first, f",
					Usage: "Prints only the most likely brute forced plaintext.",
				},
			},
		},
	}

	app.Run(os.Args)
}

/*----------------------------------------------
									CLI FUNCTIONS
-----------------------------------------------*/

func encryptMessage(c *cli.Context) {
	inputMessage := new(Message)
	inputMessage.plaintext = getUserMessage(c)
	inputMessage.key = c.Int("key")

	inputMessage.ciphertext = caesar.EncryptPlaintext(inputMessage.plaintext, inputMessage.key)

	if c.Bool("include-header") {
		fmt.Printf("%s%s%s", encryptedMessageHeader, inputMessage.ciphertext, encryptedMessageFooter)
	} else {
		fmt.Printf("%s\n", inputMessage.ciphertext)
	}
}

func decryptMessage(c *cli.Context) {
	inputMessage := new(Message)
	inputMessage.ciphertext = stripJuliusHeader(getUserMessage(c))
	inputMessage.key = c.Int("key")

	inputMessage.plaintext = caesar.DecryptCiphertext(inputMessage.ciphertext, inputMessage.key)

	fmt.Printf("%s\n", inputMessage.plaintext)
}

func bruteForceMessage(c *cli.Context) {
	frequencyMap := make(map[rune]int)

	// Get the user's input
	inputMessage := new(Message)
	inputMessage.ciphertext = stripJuliusHeader(getUserMessage(c))
	var potentialMessages [26]Message

	// Calculate the frequency of each letter in the ciphertext
	for _, letter := range strings.ToLower(inputMessage.ciphertext) {
		if letter >= 'a' && letter <= 'z' {
			frequencyMap[letter]++
		}
	}

	// Find the most frequent letter
	var mostFrequentLetter rune
	biggestFrequency := 0
	for letter, frequency := range frequencyMap {
		if frequency > biggestFrequency {
			biggestFrequency = frequency
			mostFrequentLetter = letter
		}
	}

	if c.Bool("first") {
		potentialMessage := new(Message)
		potentialMessage.key = int((26 + (mostFrequentLetter - rune(englishFrequencyList[0]))) % 26)
		potentialMessage.plaintext = caesar.DecryptCiphertext(inputMessage.ciphertext, potentialMessage.key)
		fmt.Printf("[Key: %d]: %s\n", potentialMessage.key, potentialMessage.plaintext)
	} else {
		// Determine the most probable keys based on the frequency of letters in the English Alphabet
		for index, letter := range englishFrequencyList {
			potentialMessage := new(Message)
			potentialMessage.key = int((26 + (mostFrequentLetter - letter)) % 26)
			potentialMessage.plaintext = caesar.DecryptCiphertext(inputMessage.ciphertext, potentialMessage.key)
			potentialMessages[index] = *potentialMessage
		}

		for _, message := range potentialMessages {
			fmt.Printf("[Key: %d]: %s\n", message.key, message.plaintext)
		}
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
func stripJuliusHeader(message string) string {
	message = strings.Replace(message, encryptedMessageHeader, "", 1)
	message = strings.Replace(message, encryptedMessageFooter, "", 1)

	return message
}
