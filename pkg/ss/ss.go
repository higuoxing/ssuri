package ss

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ShadowsocksURI ... Struct for shadowsocks URI.
// See: https://shadowsocks.org/en/config/quick-guide.html
// Shadowsocks for Android/iOS accepts base64 encoded URI format configs.
// e.g. ss://BASE64-ENCODED-STRING-WITHOUT-PADDING#TAG
// Where the plain URI should be
// ss://method:password@hostname:port
type ShadowsocksURI struct {
	Hostname string
	Port     int
	Method   string
	Password string
	Tag      string
	PlainURI string
}

// DecodeBase64URI ... Decode base64 encoded URI.
func DecodeBase64URI(uri string) (*ShadowsocksURI, error) {
	uri = strings.TrimSpace(uri)
	encoded, tag, err := parseTag(uri)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(encoded, "ss://") {
		return nil, fmt.Errorf("URI should have prefix \"ss://\"")
	}

	// Omit "ss://"
	encoded = strings.TrimPrefix(encoded, "ss://")

	decoded, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("Cannot decode base64 part")
	}

	// Remove spaces.
	// s := <method>:<password>@<hostname>:<port>
	s := strings.TrimSpace(string(decoded))
	plainURI := "ss://" + s

	ssconf, err := DecodePlainURI(plainURI)
	if err != nil {
		return nil, err
	}

	// Append tag
	ssconf.Tag = tag

	return ssconf, err
}

// DecodePlainURI ... Decode plain shadowsocks URI.
func DecodePlainURI(uri string) (*ShadowsocksURI, error) {
	if !strings.HasPrefix(uri, "ss://") {
		return nil, fmt.Errorf("URI should have prefix \"ss://\"")
	}

	// Omit "ss://"
	s := strings.TrimPrefix(uri, "ss://")

	// s := <method>:<password>@<hostname>
	s, port, err := parsePort(s)
	if err != nil {
		return nil, err
	}

	// s := <method>:<password>
	s, hostname, err := parseHostname(s)
	if err != nil {
		return nil, err
	}

	// s := <password>
	s, method, err := parseMethod(s)
	if err != nil {
		return nil, err
	}

	password, err := parsePassword(s)
	if err != nil {
		return nil, err
	}

	return &ShadowsocksURI{
		Hostname: hostname,
		Port:     port,
		Method:   method,
		Password: password,
		Tag:      "",
		PlainURI: uri}, nil
}

// EncodeBase64URI ... Encode shadowsocks configuration into base64 URI.
func (ssc *ShadowsocksURI) EncodeBase64URI() string {
	plainURI := ssc.Method + ":" +
		ssc.Password + "@" +
		ssc.Hostname + ":" +
		strconv.Itoa(ssc.Port)

	encoded := "ss://" +
		base64.RawStdEncoding.EncodeToString([]byte(plainURI))

	if ssc.Tag != "" {
		return encoded + "#" + ssc.Tag
	}

	return encoded
}

// EncodePlainURI ... Encode shadowsocks configuration into plain URI.
func (ssc *ShadowsocksURI) EncodePlainURI() string {
	return "ss://" +
		ssc.Method + ":" +
		ssc.Password + "@" +
		ssc.Hostname + ":" +
		strconv.Itoa(ssc.Port)
}

// ToShadowsocksClientConfig ... Convert ShadowsocksURI to ShadowsocksClientConfig.
func (ssc *ShadowsocksURI) ToShadowsocksClientConfig() ShadowsocksClientConfig {
	return ShadowsocksClientConfig{
		Server:       ssc.Hostname,
		ServerPort:   ssc.Port,
		LocalAddress: "127.0.0.1",
		LocalPort:    1080,
		Password:     ssc.Password,
		Timeout:      300,
		Method:       ssc.Method,
		FastOpen:     false,
	}
}

// ShadowsocksClientConfig ... Struct for shadowsocks client configuration.
// See: https://github.com/shadowsocks/shadowsocks/wiki/Configuration-via-Config-File
// e.g.
// {
//     "server":"my_server_ip",
//     "server_port":8388,
//     "local_address": "127.0.0.1",
//     "local_port":1080,
//     "password":"mypassword",
//     "timeout":300,
//     "method":"aes-256-cfb",
//     "fast_open": false
// }
type ShadowsocksClientConfig struct {
	Server       string `json:"server"`
	ServerPort   int    `json:"server_port"`
	LocalAddress string `json:"local_address"` // default by 127.0.0.1
	LocalPort    int    `json:"local_port"`    // default by 1080
	Password     string `json:"password"`
	Timeout      int    `json:"timeout"` // default by 300
	Method       string `json:"method"`
	FastOpen     bool   `json:"fast_open"` // default by false
}

// EncodeJSON ... Encode ShadowsocksClientConfig into JSON file.
func (sscc *ShadowsocksClientConfig) EncodeJSON() (string, error) {
	json, err := json.MarshalIndent(sscc, "", "    ")
	return string(json), err
}

// DecodeJSON ... Decode JSON configuration to ShadowsocksClientConfig.
func DecodeJSON(data []byte) (*ShadowsocksClientConfig, error) {
	var sscc ShadowsocksClientConfig

	err := json.Unmarshal(data, &sscc)
	if err != nil {
		return nil, err
	}

	return &sscc, nil
}

// ToShadowsocksURI ... Convert ShadowsocksClientConfig to ShadowsocksURI.
func (sscc *ShadowsocksClientConfig) ToShadowsocksURI() ShadowsocksURI {
	return ShadowsocksURI{
		Hostname: sscc.Server,
		Port:     sscc.ServerPort,
		Method:   sscc.Method,
		Password: sscc.Password,
		Tag:      "",
		PlainURI: "",
	}
}

// parseTag ... Parse tag.
func parseTag(uri string) (string, string, error) {
	splitURI := strings.Split(uri, "#")

	if len(splitURI) == 2 {
		return splitURI[0], splitURI[1], nil
	} else if len(splitURI) == 1 {
		return splitURI[0], "", nil
	}

	return "", "", fmt.Errorf("Malformed Shadowsocks URI")
}

// parsePort ... Parse port.
func parsePort(s string) (string, int, error) {
	// Index of ':'
	sepIndex := strings.LastIndexByte(s, ':')

	// Make sure there are characters after ':'
	if sepIndex == -1 || sepIndex >= len(s)-1 {
		return "", 0, fmt.Errorf("Cannot parse <port>")
	}

	// Safely convert port to integer.
	port, err := strconv.Atoi(s[sepIndex+1:])
	if err != nil {
		return "", 0, fmt.Errorf("Invalid <port>")
	}

	return s[:sepIndex], port, nil
}

// parseHostname ... Parse hostname.
func parseHostname(s string) (string, string, error) {
	// Index of '@'
	sepIndex := strings.LastIndexByte(s, '@')

	// Make sure there are characters after '@'
	if sepIndex == -1 || sepIndex >= len(s)-1 {
		return "", "", fmt.Errorf("Cannot parse <hostname>")
	}

	return s[:sepIndex], s[sepIndex+1:], nil
}

// parseMethod ... Parse encryption method.
func parseMethod(s string) (string, string, error) {
	// Index of ':'
	sepIndex := strings.IndexByte(s, ':')

	// Make sure there are charaters before ':'
	if sepIndex == 0 {
		return "", "", fmt.Errorf("Cannot parse <method>")
	}

	return s[sepIndex+1:], s[:sepIndex], nil
}

// parsePassword ... Parse password.
func parsePassword(s string) (string, error) {
	if len(s) == 0 {
		return "", fmt.Errorf("Cannot parse <password>")
	}

	return s, nil
}
