package ss_test

import (
	"testing"

	"github.com/vgxbj/ssuri/pkg/ss"
)

func TestURIDecode(t *testing.T) {
	tests := []struct {
		expectedConfig ss.ShadowsocksURI
		b64URI         string
		plainURI       string
		sip002URI      string
	}{
		{
			ss.ShadowsocksURI{
				Remote: ss.NewServer("192.168.100.1", 8888),
				Auth:   ss.NewAuthInfo("bf-cfb", "test"),
				Tag:    "example-server",
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
	}

	for i, ut := range tests {
		sip002, err := ss.DecodeSIP002URI(ut.sip002URI)
		if err != nil {
			t.Errorf("#%d test failed. DecodeSIP002URI() failed", i)
			continue
		}

		if !checkSIP002URI(sip002, &ut.expectedConfig) {
			t.Errorf("#%d test failed.\nExpected: %v\nGot     : %v", i, ut.expectedConfig, *sip002)
		}

		b64, err := ss.DecodeBase64URI(ut.b64URI)
		if err != nil {
			t.Errorf("#%d test failed. DecodeSIP002URI() failed", i)
			continue
		}

		if !checkBase64EncodedURI(b64, &ut.expectedConfig) {
			t.Errorf("#%d test failed.\nExpected: %v\nGot     : %v", i, ut.expectedConfig, *b64)
		}

		plain, err := ss.DecodePlainURI(ut.plainURI)
		if err != nil {
			t.Errorf("#%d test failed. DecodeSIP002URI() failed", i)
			continue
		}

		if !checkPlainURI(plain, &ut.expectedConfig) {
			t.Errorf("#%d test failed.\nExpected: %v\nGot     : %v", i, ut.expectedConfig, *plain)
		}
	}
}

func checkBasic(uri1, uri2 *ss.ShadowsocksURI) bool {
	s1 := uri1.Remote
	s2 := uri2.Remote

	if s1.Hostname() != s2.Hostname() || s1.Port() != s2.Port() {
		return false
	}

	a1 := uri1.Auth
	a2 := uri2.Auth

	if a1.Method() != a2.Method() || a1.Password() != a2.Password() {
		return false
	}

	return true
}

func checkSIP002URI(uri1, uri2 *ss.ShadowsocksURI) bool {
	if !checkBasic(uri1, uri2) {
		return false
	}

	t1 := uri1.Tag
	t2 := uri2.Tag

	if t1 != t2 {
		return false
	}

	p1 := uri1.Plugin
	p2 := uri2.Plugin

	if p1 != nil || p2 != nil {
		if p1 == nil || p2 == nil {
			return false
		}

		if p1.Name() != p2.Name() {
			return false
		}

		if len(p1.Options()) != len(p2.Options()) {
			return false
		}

		for k, v := range p1.Options() {
			if p2.Options()[k] != v {
				return false
			}
		}
	}

	return true
}

func checkBase64EncodedURI(uri1, uri2 *ss.ShadowsocksURI) bool {
	if !checkBasic(uri1, uri2) {
		return false
	}

	t1 := uri1.Tag
	t2 := uri2.Tag

	if t1 != t2 {
		return false
	}

	return true
}

func checkPlainURI(uri1, uri2 *ss.ShadowsocksURI) bool {
	if !checkBasic(uri1, uri2) {
		return false
	}

	return true
}
