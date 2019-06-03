package ssr

// ShadowsocksRURI ... Struct for shadowsocksR URI.
// See: https://github.com/shadowsocksrr/shadowsocks-rss/wiki/SSR-QRcode-scheme
// e.g.
// ssr://base64(host:port:protocol:method:obfs:base64pass/?obfsparam=base64param&protoparam=base64param&remarks=base64remarks&group=base64group&udpport=0&uot=0)
type ShadowsocksRURI struct {
	Hostname   string
	Port       int
	Protocol   string
	Method     string
	Obfs       string
	Password   string
	ObfsParam  string // optional
	ProtoParam string // optional
	Remarks    string // optional
	Group      string // optional
	UDPPort    int    // optional
	Uot        int    // optional
}

// EncodeBase64URI ... Encode shadowsocksR configuration into base64 URI.
// func (ssru *ShadowsocksRURI) EncodeBase64URI() string {
// 	var uri string

// 	uri = ssru.Hostname +
// 		":" + strconv.Itoa(ssru.Port) +
// 		":" + ssru.Protocol +
// 		":" + ssru.Method +
// 		":" + ssru.Obfs + ":"

// 	b64Password := base64.RawStdEncoding.EncodeToString([]byte(ssru.Password))

// 	uri += b64Password

// 	if 	ObfsParam  != "" || ProtoParam != "" || Remarks != ""
// 	Group      string // optional
// 	UDPPort    int    // optional
// 	Uot        int    // optional

// 	return "ssr://" + base64.RawStdEncoding.EncodeToString([]byte(uri))
// }
