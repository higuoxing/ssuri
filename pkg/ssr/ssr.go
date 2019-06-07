package ssr

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
