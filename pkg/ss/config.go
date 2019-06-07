package ss

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
	Workers      int    `json:"workers"`   // available on Unix/Linux, default by 1
}
