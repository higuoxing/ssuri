package main

import (
	"fmt"
	"os"

	"github.com/vgxbj/ssutils/pkg/ss"
)

// dumpShadowsocksURI ... dump shadowsocks base64 encoded URI
func dumpShadowsocksURI(index int, uri string, outputFile *os.File) {
	ssConfig, err := ss.DecodeBase64URI(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if ssConfig.Tag != "" {
		fmt.Fprintf(outputFile, "Server #%s:\n", ssConfig.Tag)
	} else {
		fmt.Fprintf(outputFile, "Server #%d:\n", index)
	}

	fmt.Fprintf(outputFile, "Plain URI         : %v\n", ssConfig.PlainURI)
	fmt.Fprintf(outputFile, "Hostname          : %v\n", ssConfig.Hostname)
	fmt.Fprintf(outputFile, "Port              : %v\n", ssConfig.Port)
	fmt.Fprintf(outputFile, "Encryption Method : %v\n", ssConfig.Method)
	fmt.Fprintf(outputFile, "Password          : %v\n", ssConfig.Password)
	fmt.Fprintf(outputFile, "\n")
}

// generateShadowsocksClientConfig ... Generate shadowsocks client configuration.
func generateShadowsocksClientConfig(uri string, outputFile *os.File) {
	ssConfig, err := ss.DecodeBase64URI(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	ssJSON := ssConfig.ToShadowsocksClientConfig()

	ssJSONStr, err := ssJSON.EncodeJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	fmt.Fprintf(outputFile, "%s\n", ssJSONStr)
	fmt.Fprintf(outputFile, "\n")
}

func decodeJSONToShadowsocksBase64URI(data []byte, outputFile *os.File) {
	sscc, err := ss.DecodeJSON(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	ssURI := sscc.ToShadowsocksURI()
	uri := ssURI.EncodeBase64URI()

	fmt.Fprintf(outputFile, "%s\n", uri)
}
