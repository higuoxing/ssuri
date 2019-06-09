package ss

import "encoding/json"

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
//     "fast_open": false,
//     "workers": 1
// }
type ShadowsocksClientConfig struct {
	Remote   *Server     // Remote server
	Auth     *AuthInfo   // Authentication information
	Local    *Server     // Local server, default by localhost:1080
	Timeout  int         // Connection timeout, default by 300
	FastOpen bool        // Fast open, default by false
	Workers  int         // Available on Unix/Linux, default by 1
	Plugin   *PluginInfo // Plugin, only used in defined in SIP003
}

// ShadowsocksClientJSON ... Shadowsocks client configuration in JSON format.
type ShadowsocksClientJSON struct {
	Server       string `json:"server"`
	ServerPort   int    `json:"server_port"`
	LocalAddress string `json:"local_address"`
	LocalPort    int    `json:"local_port"`
	Password     string `json:"password"`
	Timeout      int    `json:"timeout"`
	Method       string `json:"method"`
	FastOpen     bool   `json:"fast_open"`
	Workers      int    `json:"workers"`
	Plugin       string `json:"plugin"`
	PluginOpts   string `json:"plugin_opts"`
}

// NewShadowsocksClientJSON ... Generate new client configuration in JSON format.
func NewShadowsocksClientJSON(scc *ShadowsocksClientConfig) *ShadowsocksClientJSON {
	pluginName, pluginOpts := "", ""

	if scc.Plugin != nil {
		pluginName = scc.Plugin.Name()
		pluginOpts = scc.Plugin.OptionsString()
	}

	return &ShadowsocksClientJSON{
		Server:       scc.Remote.Hostname(),
		ServerPort:   scc.Remote.Port(),
		LocalAddress: scc.Local.Hostname(),
		LocalPort:    scc.Local.Port(),
		Password:     scc.Auth.Password(),
		Timeout:      scc.Timeout,
		Method:       scc.Auth.Method(),
		FastOpen:     scc.FastOpen,
		Workers:      scc.Workers,
		Plugin:       pluginName,
		PluginOpts:   pluginOpts,
	}
}

// EncodeClientJSON ... Encode ShadowsocksClietnConfig to JSON.
func EncodeClientJSON(scc *ShadowsocksClientConfig, legacy bool) ([]byte, error) {
	clientJSON := NewShadowsocksClientJSON(scc)

	if legacy {
		data, err := json.MarshalIndent(struct {
			Server       string `json:"server"`
			ServerPort   int    `json:"server_port"`
			LocalAddress string `json:"local_address"`
			LocalPort    int    `json:"local_port"`
			Password     string `json:"password"`
			Timeout      int    `json:"timeout"`
			Method       string `json:"method"`
			FastOpen     bool   `json:"fast_open"`
			Workers      int    `json:"workers"`
		}{
			Server:       clientJSON.Server,
			ServerPort:   clientJSON.ServerPort,
			LocalAddress: clientJSON.LocalAddress,
			LocalPort:    clientJSON.LocalPort,
			Password:     clientJSON.Password,
			Timeout:      clientJSON.Timeout,
			Method:       clientJSON.Method,
			FastOpen:     clientJSON.FastOpen,
			Workers:      clientJSON.Workers,
		}, "", "    ")

		return data, err
	}

	data, err := json.MarshalIndent(clientJSON, "", "    ")

	return data, err
}

// DecodeJSON ... Decode JSON to ShadowsocksClietnConfig.
func DecodeJSON(data []byte) (*ShadowsocksClientConfig, error) {
	var clientJSON ShadowsocksClientJSON

	err := json.Unmarshal(data, &clientJSON)

	if err != nil {
		return nil, err
	}

	opts, err := ParsePluginOpts(clientJSON.PluginOpts)
	if err != nil {
		return nil, err
	}

	var plugin *PluginInfo
	if len(opts) != 0 {
		plugin = NewPlugin(clientJSON.Plugin, opts)
	}

	scc := &ShadowsocksClientConfig{
		Remote:   NewServer(clientJSON.Server, clientJSON.ServerPort),
		Auth:     NewAuthInfo(clientJSON.Method, clientJSON.Password),
		Local:    NewServer(clientJSON.LocalAddress, clientJSON.LocalPort),
		Timeout:  clientJSON.Timeout,
		FastOpen: clientJSON.FastOpen,
		Workers:  clientJSON.Workers,
		Plugin:   plugin,
	}

	return scc, nil
}

// ToShadowsocksURI ... Convert JSON configuration to shadowsocks URI.
func ToShadowsocksURI(scc *ShadowsocksClientConfig) *ShadowsocksURI {
	return &ShadowsocksURI{
		Remote: scc.Remote,
		Auth:   scc.Auth,
		Tag:    scc.Remote.String(),
		Plugin: scc.Plugin,
	}
}
