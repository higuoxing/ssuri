package ss_test

import (
	"testing"

	"github.com/vgxbj/ssutils/pkg/ss"
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
				Hostname: "192.168.100.1",
				Port:     8888,
				Method:   "bf-cfb",
				Password: "test",
				Tag:      "example-server",
				PlainURI: "",
				Plugins:  nil,
			},
			"ss://YmYtY2ZiOnRlc3RAMTkyLjE2OC4xMDAuMTo4ODg4#example-server",
			"ss://bf-cfb:test@192.168.100.1:8888",
			"ss://YmYtY2ZiOnRlc3Q=@192.168.100.1:8888#example-server",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "192.168.101.101",
				Port:     8889,
				Method:   "bf-cfb",
				Password: "some_password",
				Tag:      "",
				PlainURI: "",
				Plugins:  nil,
			},
			"ss://YmYtY2ZiOnNvbWVfcGFzc3dvcmRAMTkyLjE2OC4xMDEuMTAxOjg4ODk",
			"ss://bf-cfb:some_password@192.168.101.101:8889",
			"ss://YmYtY2ZiOnNvbWVfcGFzc3dvcmQ=@192.168.101.101:8889",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "192.168.102.3",
				Port:     443,
				Method:   "bf-cfb",
				Password: "some_password@?",
				Tag:      "",
				PlainURI: "",
				Plugins:  nil,
			},
			"ss://YmYtY2ZiOnNvbWVfcGFzc3dvcmRAP0AxOTIuMTY4LjEwMi4zOjQ0Mw",
			"ss://bf-cfb:some_password@?@192.168.102.3:443",
			"ss://YmYtY2ZiOnNvbWVfcGFzc3dvcmRAPw==@192.168.102.3:443",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "some.host",
				Port:     8888,
				Method:   "bf-cfb",
				Password: "test/!@#:",
				Tag:      "",
				PlainURI: "",
				Plugins:  nil,
			},
			"ss://YmYtY2ZiOnRlc3QvIUAjOkBzb21lLmhvc3Q6ODg4OA",
			"ss://bf-cfb:test/!@#:@some.host:8888",
			"ss://YmYtY2ZiOnRlc3QvIUAjOg==@some.host:8888",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "192.168.100.1",
				Port:     8888,
				Method:   "aes-128-gcm",
				Password: "test",
				Tag:      "Example1",
				PlainURI: "",
				Plugins:  nil,
			},
			"ss://YWVzLTEyOC1nY206dGVzdEAxOTIuMTY4LjEwMC4xOjg4ODg#Example1",
			"ss://aes-128-gcm:test@192.168.100.1:8888",
			"ss://YWVzLTEyOC1nY206dGVzdA==@192.168.100.1:8888#Example1",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "192.168.100.1",
				Port:     8888,
				Method:   "rc4-md5",
				Password: "passwd",
				Tag:      "Example2",
				PlainURI: "",
				Plugins: map[string]string{
					"plugin": "obfs-local;obfs=http",
				},
			},
			"ss://cmM0LW1kNTpwYXNzd2RAMTkyLjE2OC4xMDAuMTo4ODg4#Example2",
			"ss://rc4-md5:passwd@192.168.100.1:8888",
			"ss://cmM0LW1kNTpwYXNzd2Q=@192.168.100.1:8888/?plugin=obfs-local%3Bobfs%3Dhttp#Example2",
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
		clientConfig ss.ShadowsocksClientConfig
		expected     string
	}{
		{
			ss.ShadowsocksClientConfig{
				Server:       "some_host",
				ServerPort:   8118,
				LocalAddress: "127.0.0.1",
				LocalPort:    1080,
				Password:     "test#@a",
				Timeout:      300,
				Method:       "bf-cfb",
				FastOpen:     true,
			},
			"{\n    \"server\": \"some_host\",\n    \"server_port\": 8118,\n    \"local_address\": \"127.0.0.1\",\n    \"local_port\": 1080,\n    \"password\": \"test#@a\",\n    \"timeout\": 300,\n    \"method\": \"bf-cfb\",\n    \"fast_open\": true\n}",
		},
	}

	for i, ut := range tests {
		json, err := ut.clientConfig.EncodeJSON()
		if err != nil {
			t.Errorf("%v", err)
		}

		if string(json) != ut.expected {
			t.Errorf("#%d test failed. Expected:\n%v, Got:\n%v", i, ut.expected, json)
		}
	}
}
