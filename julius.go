package main

import (
  "bufio"
  "code.google.com/p/go.crypto/ssh/terminal"
  "fmt"
  "github.com/alexjohnj/caesar"
  "github.com/codegangsta/cli"
  "log"
  "os"
  "strings"
)

const encryptedMessageHeader = "-----BEGIN JULIUS MESSAGE-----\n\n"
const encryptedMessageFooter = "\n\n-----END JULIUS MESSAGE-----"

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
        cli.BoolFlag{
          Name:  "include-header, b",
          Usage: "include a pgp style header in output",
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
  plaintext := getUserMessage(c)
  key := c.Int("key")

  ciphertext := caesar.EncryptPlaintext(plaintext, key)

  if c.Bool("include-header") {
    fmt.Printf("%s%s%s", encryptedMessageHeader, ciphertext, encryptedMessageFooter)
  } else {
    fmt.Println(ciphertext)
  }
}

func decryptMessage(c *cli.Context) {
  ciphertext := getUserMessage(c)
  ciphertext = stripJuliusHeader(c, ciphertext)
  key := c.Int("key")

  plaintext := caesar.DecryptCiphertext(ciphertext, key)

  fmt.Println(plaintext)
}

func readFromFile(f *os.File) string {
  var fileContent string
  fileScanner := bufio.NewScanner(f)

  // Read the first line from the File
  fileScanner.Scan()
  fileContent = fileScanner.Text()

  // Read the remaining lines
  for fileScanner.Scan() {
    fileContent = strings.Join([]string{fileContent, fileScanner.Text()}, "\n")

    if err := fileScanner.Err(); err != nil {
      log.Fatal(err)
    }
  }
  return fileContent
}

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

func stripJuliusHeader(c *cli.Context, message string) string {
  message = strings.Replace(message, encryptedMessageHeader, "", 1)
  message = strings.Replace(message, encryptedMessageFooter, "", 1)

  return message
}
