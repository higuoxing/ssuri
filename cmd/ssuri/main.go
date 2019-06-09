package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/vgxbj/ssuri/pkg/ss"
)

var inputFileName *string    // input file name, option -i, default stdin
var outputFileName *string   // output file name, option -o, default stdout
var jsonMode *bool           // run in JSON mode, option -json, default off
var dumpURI *bool            // dump URI information
var legacyMode *bool         // dump shadowsocks URI in legacy mode, option -legacy, default off
var generateJSONConfig *bool // generate JSON config, option -generate-json-config
var generateQRCode *bool     // generate QR code, option -generate-qr.
var generateURI *bool        // generate URI.

func init() {
	inputFileName = flag.String("i", "-", "input file (default: \"-\" for stdin)")
	outputFileName = flag.String("o", "-", "output file (default: \"-\" for stdout)")
	jsonMode = flag.Bool("json", false, "read JSON as input (default: off)")
	dumpURI = flag.Bool("dump-uri", false, "dump shadowsocks URI")
	legacyMode = flag.Bool("legacy", false, "dump shadowsocks URI in legacy mode (default: off)")
	generateJSONConfig = flag.Bool("generate-json-config", false, "generate JSON configurations")
	generateQRCode = flag.Bool("generate-qr", false, "generate QR code")
	generateURI = flag.Bool("generate-uri", false, "generate URI")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [-h] [-i in_file] [-o out_file]\n", os.Args[0] /* Program name */)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	var inputFile *os.File
	var outputFile *os.File

	var err error

	if *inputFileName != "-" {
		inputFile, err = os.Open(*inputFileName)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		defer inputFile.Close()
	} else {
		inputFile = os.Stdin
	}

	if *outputFileName != "-" {
		outputFile, err = os.Create(*outputFileName)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	// Process input file.
	data, err := ioutil.ReadAll(inputFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	s := strings.TrimSpace(string(data))

	var clientConfig *ss.ShadowsocksClientConfig
	var uri *ss.ShadowsocksURI

	if *jsonMode {
		// Read JSON configuration.
		clientConfig, err = decodeJSONConfig([]byte(s))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		uri = generateShadowsocksURI(clientConfig)
	} else {
		// Read shadowsocks URI.
		uri, err = decodeURI(s, *legacyMode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		clientConfig = generateShadowsocksClientConfig(uri)
	}

	if *dumpURI {
		dumpShadowsocksURI(uri, outputFile)
	}

	if *generateJSONConfig {
		generateClientJSONConfig(clientConfig, outputFile)
	}

	if *generateQRCode {
		generateShadowsocksQRCode(uri, false, outputFile)
	}

	if *generateURI {
		fmt.Fprintf(outputFile, "%s", uri.EncodeSIP002URI())
	}
}
