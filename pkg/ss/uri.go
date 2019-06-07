package ss

import (
	"encoding/base64"
	"errors"
	"net/url"
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
	Server *RemoteServer
	Auth   *AuthInfo
	Tag    string      // optional
	Plugin *PluginInfo // optional, used in SIP002 URI scheme
}

// Server ... Returns a RemoteServer containing the given hostname and port.
func Server(hostname string, port int) *RemoteServer {
	return &RemoteServer{hostname, port}
}

// Auth ... Returns a AuthInfo containing the given method and password.
func Auth(method, password string) *AuthInfo {
	return &AuthInfo{method, password}
}

// Plugin ... Returns a PluginInfo containing the given name and options.
func Plugin(name string, options map[string]string) *PluginInfo {
	return &PluginInfo{name, options}
}

// RemoteServer ... Struct for shadowsocks remote server.
type RemoteServer struct {
	hostname string
	port     int
}

// Hostname ... Returns hostname.
func (rs *RemoteServer) Hostname() string {
	return rs.hostname
}

// Port ... Returns port.
func (rs *RemoteServer) Port() int {
	return rs.port
}

// String ... Encode remote server address.
func (rs *RemoteServer) String() string {
	if strings.Contains(rs.hostname, ":") {
		return "[" + rs.hostname + "]" + ":" + strconv.Itoa(rs.port)
	}

	return rs.hostname + ":" + strconv.Itoa(rs.port)
}

// AuthInfo ... Struct for authentication information in ShadowsocksURI.
// e.g. <method>:<password>
type AuthInfo struct {
	method   string
	password string
}

// Method ... Returns method of authentication information.
func (auth *AuthInfo) Method() string {
	return auth.method
}

// Password ... Returns password of authentication information.
func (auth *AuthInfo) Password() string {
	return auth.password
}

// String ... Return the encoded authentication information.
func (auth *AuthInfo) String() string {
	return auth.method + ":" + auth.password
}

// PluginInfo ... Struct for shadowsocks plugin (SIP002).
type PluginInfo struct {
	name    string            // Plugin name
	options map[string]string // Plugin options
}

// Name ... Returns name of plugin.
func (plugin *PluginInfo) Name() string {
	return plugin.name
}

// Options ... Returns options of plugin.
func (plugin *PluginInfo) Options() map[string]string {
	return plugin.options
}

// String ... Return the encoded plugin information.
func (plugin *PluginInfo) String() string {
	builder := url.Values{}

	var options = plugin.name

	for k, v := range plugin.options {
		options += ";" + k + "=" + v
	}

	// SIP002 URI scheme only supports one plugin.
	// See: https://shadowsocks.org/en/spec/SIP002-URI-Scheme.html
	builder.Add("plugin", options)

	return "/?" + builder.Encode()
}

// EncodeSIP002URI ... Encode shadowsocks configuration into SIP002 URI.
func (uri *ShadowsocksURI) EncodeSIP002URI() string {
	// Encode auth information.
	auth := uri.Auth.String()
	auth = base64.URLEncoding.EncodeToString([]byte(auth))

	// Add hostname, port.
	wrappedHost := uri.Server.String()

	// Encode plugin parameters.
	var plugin = ""
	if uri.Plugin != nil {
		plugin = uri.Plugin.String()
	}

	// Encode tag.
	tag := uri.encodeFragment()

	// ss://base64(method:password)@<addr>:<port> ["/"] ["?" plugin=<plugin_name>;<opt_name>=<option>+]
	return "ss://" + auth + "@" + wrappedHost + plugin + tag
}

// EncodeBase64URI ... Encode shadowsocks configuration into base64 URI (legacy).
func (uri *ShadowsocksURI) EncodeBase64URI() string {
	auth := uri.Auth.String()

	wrappedHost := uri.Server.String()

	plainURI := auth + "@" + wrappedHost

	encoded := "ss://" +
		base64.RawURLEncoding.EncodeToString([]byte(plainURI))

	tag := uri.encodeFragment()

	return encoded + tag
}

// EncodePlainURI ... Encode shadowsocks configuration into plain URI.
func (uri *ShadowsocksURI) EncodePlainURI() string {
	auth := uri.Auth.String()

	wrappedHost := uri.Server.String()

	return "ss://" + auth + "@" + wrappedHost
}

// encodeFragment ... Encode fragment.
func (uri *ShadowsocksURI) encodeFragment() string {
	if uri.Tag != "" {
		return "#" + uri.Tag
	}

	return ""
}

// DecodeSIP002URI ... Decode SIP002 shadowsocks URI.
func DecodeSIP002URI(uri string) (*ShadowsocksURI, error) {
	// Omit "ss://"
	// s := base64(<auth>)@<hostname>:<port> [ "/" ] [ "?" <plugin> ] [ "#" <tag> ]
	s, ok := checkPrefixAndTrim(uri, "ss://")
	if !ok {
		return nil, errors.New("invalid <scheme>")
	}

	// s := base64(<auth>)@<hostname>:<port> [ "/" ] [ "?" <plugin> ]
	s, tag, err := parseTag(s)
	if err != nil {
		return nil, err
	}

	// s := <hostname>:<port> [ "/" ] [ "?" <plugin> ]
	// authStr := base64(<method>:<password>)
	authStr, s, err := splitAuthAndHost(s)
	if err != nil {
		return nil, err
	}

	// decodedAuth := <method>:<password>
	decodedAuthStr, err := base64.URLEncoding.DecodeString(authStr)
	if err != nil {
		return nil, errors.New("invalid base64 encoded <auth>")
	}

	// method := <method>
	// password := <password>
	auth, err := parseAuth(string(decodedAuthStr))
	if err != nil {
		return nil, err
	}

	hostStr, pluginStr, err := splitRemoteAndPlugin(s)
	if err != nil {
		return nil, err
	}

	host, err := parseRemoteServer(hostStr)
	if err != nil {
		return nil, err
	}

	plugin, err := parsePlugin(pluginStr)
	if err != nil {
		return nil, err
	}

	return &ShadowsocksURI{
		Server: host,
		Auth:   auth,
		Tag:    tag,
		Plugin: plugin,
	}, nil
}

// DecodeBase64URI ... Decode base64 encoded URI (legacy).
func DecodeBase64URI(uri string) (*ShadowsocksURI, error) {
	// Omit "ss://"
	// s := base64(<auth>@<hostname>:<port>)#tag
	s, ok := checkPrefixAndTrim(uri, "ss://")
	if !ok {
		return nil, errors.New("invalid <scheme>")
	}

	s, tag, err := parseTag(s)
	if err != nil {
		return nil, err
	}

	// decoded := <auth>@<hostname>:<port>
	decoded, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	// authStr := <auth>
	// s := <hostname>:<port>
	authStr, s, err := splitAuthAndHost(string(decoded))
	if err != nil {
		return nil, err
	}

	auth, err := parseAuth(authStr)
	if err != nil {
		return nil, err
	}

	host, err := parseRemoteServer(s)
	if err != nil {
		return nil, err
	}

	return &ShadowsocksURI{
		Server: host,
		Auth:   auth,
		Tag:    tag,
		Plugin: nil,
	}, nil
}

// DecodePlainURI ... Decode plain shadowsocks URI.
func DecodePlainURI(uri string) (*ShadowsocksURI, error) {
	// Omit "ss://"
	s, ok := checkPrefixAndTrim(uri, "ss://")
	if !ok {
		return nil, errors.New("invalid <scheme>")
	}

	// authStr := <auth>
	// s := <hostname>:<port>
	authStr, s, err := splitAuthAndHost(s)
	if err != nil {
		return nil, err
	}

	auth, err := parseAuth(authStr)
	if err != nil {
		return nil, err
	}

	host, err := parseRemoteServer(s)
	if err != nil {
		return nil, err
	}

	return &ShadowsocksURI{
		Server: host,
		Auth:   auth,
		Tag:    "",
		Plugin: nil,
	}, nil
}

// checkPrefixAndTrim ... Check given prefix and remove it.
func checkPrefixAndTrim(s, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		return strings.TrimPrefix(s, prefix), true
	}

	return s, false
}

// parseTag ... Parse tag.
func parseTag(uri string) (string, string, error) {
	splitURI := strings.Split(uri, "#")

	if len(splitURI) == 2 {
		return splitURI[0], splitURI[1], nil
	} else if len(splitURI) == 1 {
		return splitURI[0], "", nil
	}

	return "", "", errors.New("invalid URI")
}

// splitAuthAndHost ... Split authentication information and hostname.
func splitAuthAndHost(s string) (string, string, error) {
	// s := <auth>@<hostname>:<port>
	splitIndex := strings.LastIndexByte(s, '@')

	if splitIndex == -1 || splitIndex == 0 || splitIndex == len(s)-1 {
		return "", "", errors.New("invalid <auth> and <hostname>")
	}

	return s[:splitIndex], s[splitIndex+1:], nil
}

// parseAuth ... Parse authentication information.
func parseAuth(authStr string) (*AuthInfo, error) {
	// <auth> := <method>:<password>
	splitIndex := strings.IndexByte(authStr, ':')

	if splitIndex == -1 || splitIndex == 0 || splitIndex == len(authStr)-1 {
		return nil, errors.New("invalid <auth>")
	}

	return &AuthInfo{authStr[:splitIndex], authStr[splitIndex+1:]}, nil
}

// splitRemoteAndPlugin ... Split remote server information and plugin (used in SIP002 URI).
func splitRemoteAndPlugin(uri string) (string, string, error) {
	// Add "//" prefix. See: https://golang.org/src/net/url/url.go#L508
	parsedURI, err := url.Parse("//" + uri)

	if err != nil {
		return "", "", errors.New("invalid URI")
	}

	return parsedURI.Host, parsedURI.Query().Get("plugin"), nil
}

// parseRemoteServer ... Parse remote server.
func parseRemoteServer(hostStr string) (*RemoteServer, error) {
	splitIndex := strings.LastIndexByte(hostStr, ':')

	if splitIndex == -1 || splitIndex == 0 || splitIndex == len(hostStr)-1 {
		return nil, errors.New("invalid <hostname>:<port>")
	}

	port, err := strconv.Atoi(hostStr[splitIndex+1:])
	if err != nil {
		return nil, errors.New("invalid <port>")
	}

	return &RemoteServer{hostStr[:splitIndex], port}, nil
}

// parsePlugin ... Parse plugin (used in SIP002 URI scheme).
func parsePlugin(pluginStr string) (*PluginInfo, error) {
	if pluginStr == "" {
		return nil, nil
	}

	nameAndOptions := strings.Split(pluginStr, ";")

	if len(nameAndOptions) < 2 {
		return nil, errors.New("invalid <plugin>")
	}

	var name string
	options := make(map[string]string)

	for i, o := range nameAndOptions {
		if i == 0 {
			// name of plugin
			name = o
			continue
		}

		kv := strings.Split(o, "=")
		if len(kv) == 2 {
			options[kv[0]] = kv[1]
		}
	}

	return Plugin(name, options), nil
}
