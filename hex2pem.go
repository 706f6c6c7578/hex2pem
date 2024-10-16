package main

import (
	"encoding/hex"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  Convert PEM to hex:  hex2pem -h < infile > outfile")
	fmt.Println("  Convert hex to PEM:  hex2pem -p < infile > outfile")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
}

func main() {
	hexFlag := flag.Bool("h", false, "Convert PEM to hex")
	pemFlag := flag.Bool("p", false, "Convert hex to PEM")
	flag.Parse()

	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	if *hexFlag == *pemFlag {
		fmt.Println("Error: Please specify either -h or -p flag")
		printUsage()
		os.Exit(1)
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	if *hexFlag {
		// PEM to hex
		block, _ := pem.Decode(input)
		if block == nil {
			fmt.Println("Failed to decode PEM block")
			os.Exit(1)
		}
		fmt.Println(hex.EncodeToString(block.Bytes))
	} else {
		// hex to PEM
		hexString := strings.TrimSpace(string(input))
		data, err := hex.DecodeString(hexString)
		if err != nil {
			fmt.Println("Error decoding hex:", err)
			os.Exit(1)
		}

		var pemType string
		if len(data) == 32 {
			pemType = "PUBLIC KEY"
		} else if len(data) == 64 {
			pemType = "PRIVATE KEY"
		} else {
			fmt.Println("Invalid key length")
			os.Exit(1)
		}

		pemBlock := &pem.Block{
			Type:  pemType,
			Bytes: data,
		}
		err = pem.Encode(os.Stdout, pemBlock)
		if err != nil {
			fmt.Println("Error encoding PEM:", err)
			os.Exit(1)
		}
	}
}
