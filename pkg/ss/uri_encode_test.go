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
				Server: ss.Server("192.168.100.1", 8888),

				Auth:   ss.Auth("bf-cfb", "test"),
				Tag:    "example-server",
				Plugin: nil,
			},
			"ss://YmYtY2ZiOnRlc3RAMTkyLjE2OC4xMDAuMTo4ODg4#example-server",
			"ss://bf-cfb:test@192.168.100.1:8888",
			"ss://YmYtY2ZiOnRlc3Q=@192.168.100.1:8888#example-server",
		},
		{
			ss.ShadowsocksURI{
				Server: ss.Server("192.168.100.1", 8888),
				Auth:   ss.Auth("rc4-md5", "passwd"),
				Tag:    "example-server",
				Plugin: ss.Plugin("obfs-local", map[string]string{"obfs": "http"}),
			},
			"ss://cmM0LW1kNTpwYXNzd2RAMTkyLjE2OC4xMDAuMTo4ODg4#example-server",
			"ss://rc4-md5:passwd@192.168.100.1:8888",
			"ss://cmM0LW1kNTpwYXNzd2Q=@192.168.100.1:8888/?plugin=obfs-local%3Bobfs%3Dhttp#example-server",
		},
		{
			ss.ShadowsocksURI{
				Server: ss.Server("test.example.com", 8888),
				Auth:   ss.Auth("rc4-md5", "passwd"),
				Tag:    "example-server",
				Plugin: ss.Plugin("obfs-local", map[string]string{"obfs": "http"}),
			},
			"ss://cmM0LW1kNTpwYXNzd2RAdGVzdC5leGFtcGxlLmNvbTo4ODg4#example-server",
			"ss://rc4-md5:passwd@test.example.com:8888",
			"ss://cmM0LW1kNTpwYXNzd2Q=@test.example.com:8888/?plugin=obfs-local%3Bobfs%3Dhttp#example-server",
		},
		{
			ss.ShadowsocksURI{
				Server: ss.Server("192.168.100.1", 8888),
				Auth:   ss.Auth("bf-cfb", "test"),
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
