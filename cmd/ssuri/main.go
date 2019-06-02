package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var inputFileName *string    // input file name, option -i
var outputFileName *string   // output file name, option -o
var jsonMode *bool           // run in JSON mode, option -json, default off
var dumpURI *bool            // dump URI information
var generateJSONConfig *bool // generate JSON config, option -generate-json-config

func init() {
	inputFileName = flag.String("i", "-", "input file (default: \"-\" for stdin)")
	outputFileName = flag.String("o", "-", "output file (default: \"-\" for stdout)")
	jsonMode = flag.Bool("json", false, "read JSON as input (default: off)")
	dumpURI = flag.Bool("dump-uri", false, "dump base64 encoded URI")
	generateJSONConfig = flag.Bool("generate-json-config", false, "generate JSON configurations")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [-h] [-i in_file] [-o out_file]\n", os.Args[0] /* Program name */)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	var inputFile *os.File
	var outputFile *os.File

	if *inputFileName != "-" {
		var err error

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
		var err error

		outputFile, err = os.Create(*outputFileName)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	if *jsonMode {
		data, err := ioutil.ReadAll(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}

		decodeJSONToShadowsocksBase64URI(data, outputFile)
	} else {
		scanner := bufio.NewScanner(inputFile)
		scanner.Split(bufio.ScanWords)

		var index = 1

		for scanner.Scan() {
			uri := scanner.Text()

			if *dumpURI {
				dumpShadowsocksURI(index, uri, outputFile)
			}

			if *generateJSONConfig {
				generateShadowsocksClientConfig(uri, outputFile)
			}

			index++
		}
	}
}
