package main

import (
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/vgxbj/ssuri/pkg/ss"
)

// decodeJSONConfig ... Decode JSON configuration.
func decodeJSONConfig(data []byte) (*ss.ShadowsocksClientConfig, error) {
	return ss.DecodeJSON(data)
}

// decodeURI ... Decode shadowsocks URI.
func decodeURI(uri string, legacy bool) (*ss.ShadowsocksURI, error) {
	if legacy {
		return ss.DecodeBase64URI(uri)
	}

	return ss.DecodeSIP002URI(uri)
}

// generateShadowsocksClientConfig ... Generate shadowsocks client configuration.
func generateShadowsocksClientConfig(uri *ss.ShadowsocksURI) *ss.ShadowsocksClientConfig {
	return ss.ToShadowsocksClientConfig(uri)
}

// generateShadowsocksURI ... Generate shadowsocks URI scheme.
func generateShadowsocksURI(scc *ss.ShadowsocksClientConfig) *ss.ShadowsocksURI {
	return ss.ToShadowsocksURI(scc)
}

// dumpShadowsocksURI ... dump shadowsocks base64 encoded URI.
func dumpShadowsocksURI(ssu *ss.ShadowsocksURI, outputFile *os.File) {
	if ssu.Tag != "" {
		fmt.Fprintf(outputFile, "Server #%s:\n", ssu.Tag)
	} else {
		fmt.Fprintf(outputFile, "Server #%s:\n", ssu.Remote.String())
	}

	fmt.Fprintf(outputFile, "Hostname          : %v\n", ssu.Remote.Hostname())
	fmt.Fprintf(outputFile, "Port              : %v\n", ssu.Remote.Port())
	fmt.Fprintf(outputFile, "Encryption Method : %v\n", ssu.Auth.Method())
	fmt.Fprintf(outputFile, "Password          : %v\n", ssu.Auth.Password())
	fmt.Fprintf(outputFile, "\n")
}

// generateClientJSONConfig ... Generate JSON configuration.
func generateClientJSONConfig(scc *ss.ShadowsocksClientConfig, outputFile *os.File) {
	json, err := ss.EncodeClientJSON(scc, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	fmt.Fprintf(outputFile, "%s\n", string(json))
	fmt.Fprintf(outputFile, "\n")
}

// generateShadowsocksQRCode ... Generate QR code.
func generateShadowsocksQRCode(ssu *ss.ShadowsocksURI, legacy bool, outputFile *os.File) {
	var uri string

	if legacy {
		uri = ssu.EncodeBase64URI()
	} else {
		uri = ssu.EncodeSIP002URI()
	}

	qrterminal.Generate(uri, qrterminal.M, outputFile)
	fmt.Fprintf(outputFile, "\n")
}
