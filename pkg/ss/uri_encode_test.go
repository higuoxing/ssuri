package ss_test

import (
	"testing"

	"github.com/vgxbj/ssuri/pkg/ss"
)

func TestURIEncode(t *testing.T) {
	tests := []struct {
		config         ss.ShadowsocksURI
		b64Expected    string
		plainExpected  string
		sip002Expected string
	}{
		{
			ss.ShadowsocksURI{
				Remote: ss.NewServer("192.168.100.1", 8888),
				Auth:   ss.NewAuthInfo("bf-cfb", "test"),
				Tag:    "example-server",
				Plugin: nil,
			},
			"ss://YmYtY2ZiOnRlc3RAMTkyLjE2OC4xMDAuMTo4ODg4#example-server",
			"ss://bf-cfb:test@192.168.100.1:8888",
			"ss://YmYtY2ZiOnRlc3Q=@192.168.100.1:8888#example-server",
		},
		{
			ss.ShadowsocksURI{
				Remote: ss.NewServer("192.168.100.1", 8888),
				Auth:   ss.NewAuthInfo("rc4-md5", "passwd"),
				Tag:    "example-server",
				Plugin: ss.NewPlugin("obfs-local", map[string]string{"obfs": "http"}),
			},
			"ss://cmM0LW1kNTpwYXNzd2RAMTkyLjE2OC4xMDAuMTo4ODg4#example-server",
			"ss://rc4-md5:passwd@192.168.100.1:8888",
			"ss://cmM0LW1kNTpwYXNzd2Q=@192.168.100.1:8888/?plugin=obfs-local%3Bobfs%3Dhttp#example-server",
		},
		{
			ss.ShadowsocksURI{
				Remote: ss.NewServer("test.example.com", 8888),
				Auth:   ss.NewAuthInfo("rc4-md5", "passwd"),
				Tag:    "example-server",
				Plugin: ss.NewPlugin("obfs-local", map[string]string{"obfs": "http"}),
			},
			"ss://cmM0LW1kNTpwYXNzd2RAdGVzdC5leGFtcGxlLmNvbTo4ODg4#example-server",
			"ss://rc4-md5:passwd@test.example.com:8888",
			"ss://cmM0LW1kNTpwYXNzd2Q=@test.example.com:8888/?plugin=obfs-local%3Bobfs%3Dhttp#example-server",
		},
		{
			ss.ShadowsocksURI{
				Remote: ss.NewServer("192.168.100.1", 8888),
				Auth:   ss.NewAuthInfo("bf-cfb", "test"),
				Tag:    "",
				Plugin: nil,
			},
			"ss://YmYtY2ZiOnRlc3RAMTkyLjE2OC4xMDAuMTo4ODg4",
			"ss://bf-cfb:test@192.168.100.1:8888",
			"ss://YmYtY2ZiOnRlc3Q=@192.168.100.1:8888",
		},
	}

	for i, ut := range tests {
		if b64URI := ut.config.EncodeBase64URI(); b64URI != ut.b64Expected {
			t.Errorf("#%d test failed. Expected: %v, Got: %v", i, ut.b64Expected, b64URI)
		}

		if plainURI := ut.config.EncodePlainURI(); plainURI != ut.plainExpected {
			t.Errorf("#%d test failed. Expected: %v, Got: %v", i, ut.plainExpected, plainURI)
		}

		if sip002URI := ut.config.EncodeSIP002URI(); sip002URI != ut.sip002Expected {
			t.Errorf("#%d test failed. Expected: %v, Got: %v", i, ut.sip002Expected, sip002URI)
		}
	}
}

func TestShadowsocksClientConfigEncode(t *testing.T) {
	tests := []struct {
		clientConfig   ss.ShadowsocksClientConfig
		legacy         bool
		expected       string
		expectedSIP002 string
	}{
		{
			ss.ShadowsocksClientConfig{
				Remote:   ss.NewServer("some_host", 8118),
				Local:    ss.NewServer("127.0.0.1", 1080),
				Auth:     ss.NewAuthInfo("bf-cfb", "test#@a"),
				Timeout:  300,
				FastOpen: true,
				Workers:  1,
			},
			true,
			"{\n    \"server\": \"some_host\",\n    \"server_port\": 8118,\n    \"local_address\": \"127.0.0.1\",\n    \"local_port\": 1080,\n    \"password\": \"test#@a\",\n    \"timeout\": 300,\n    \"method\": \"bf-cfb\",\n    \"fast_open\": true,\n    \"workers\": 1\n}",
			"ss://YmYtY2ZiOnRlc3QjQGE=@some_host:8118#some_host:8118",
		},
		{
			ss.ShadowsocksClientConfig{
				Remote:   ss.NewServer("some_host", 8118),
				Local:    ss.NewServer("127.0.0.1", 1080),
				Auth:     ss.NewAuthInfo("bf-cfb", "test#@a"),
				Timeout:  300,
				FastOpen: true,
				Workers:  1,
				Plugin:   ss.NewPlugin("obfs-local", map[string]string{"obfs": "http", "obfs-host": "www.baidu.com"}),
			},
			false,
			"{\n    \"server\": \"some_host\",\n    \"server_port\": 8118,\n    \"local_address\": \"127.0.0.1\",\n    \"local_port\": 1080,\n    \"password\": \"test#@a\",\n    \"timeout\": 300,\n    \"method\": \"bf-cfb\",\n    \"fast_open\": true,\n    \"workers\": 1,\n    \"plugin\": \"obfs-local\",\n    \"plugin_opts\": \"obfs=http;obfs-host=www.baidu.com\"\n}",
			"ss://YmYtY2ZiOnRlc3QjQGE=@some_host:8118/?plugin=obfs-local%3Bobfs%3Dhttp%3Bobfs-host%3Dwww.baidu.com#some_host:8118",
		},
	}

	for i, ut := range tests {
		json, err := ss.EncodeClientJSON(&ut.clientConfig, ut.legacy)
		if err != nil {
			t.Errorf("%v", err)
		}

		if string(json) != ut.expected {
			t.Errorf("#%d test failed. Expected:\n%v, Got:\n%v", i, ut.expected, string(json))
		}

		cc, err := ss.DecodeJSON(json)
		if err != nil {
			t.Errorf("%v", err)
		}

		if !checkClientConfig(cc, &ut.clientConfig) {
			t.Errorf("#%d test failed. Expected:\n%v,\nGot:\n%v", i, ut.clientConfig, *cc)
		}

		sip002 := ss.ToShadowsocksURI(cc)

		if sip002.EncodeSIP002URI() != ut.expectedSIP002 {
			t.Errorf("#%d test failed. Expected:\n%v, Got:\n%v", i, ut.expectedSIP002, sip002.EncodeSIP002URI())
		}
	}
}

func checkClientConfig(cc1, cc2 *ss.ShadowsocksClientConfig) bool {
	if cc1.Remote.Hostname() != cc2.Remote.Hostname() || cc1.Remote.Port() != cc2.Remote.Port() {
		return false
	}

	if cc1.Local.Hostname() != cc2.Local.Hostname() || cc1.Local.Port() != cc2.Local.Port() {
		return false
	}

	if cc1.Auth.Method() != cc2.Auth.Method() || cc1.Auth.Password() != cc2.Auth.Password() {
		return false
	}

	if cc1.FastOpen != cc2.FastOpen || cc1.Timeout != cc2.Timeout || cc1.Workers != cc2.Workers {
		return false
	}

	if cc1.Plugin != nil || cc2.Plugin != nil {
		if cc1.Plugin == nil || cc2.Plugin == nil {
			return false
		}

		if cc1.Plugin.String() != cc2.Plugin.String() {
			return false
		}
	}

	return true
}
