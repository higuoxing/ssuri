package ss_test

import (
	"testing"

	"github.com/vgxbj/ssutils/pkg/ss"
)

func TestURIDecode(t *testing.T) {
	tests := []struct {
		expectedConfig ss.ShadowsocksURI
		b64URI         string
		plainURI       string
	}{
		{
			ss.ShadowsocksURI{
				Hostname: "192.168.100.1",
				Port:     8888,
				Method:   "bf-cfb",
				Password: "test",
				Tag:      "example-server",
				PlainURI: "ss://bf-cfb:test@192.168.100.1:8888",
			},
			"ss://YmYtY2ZiOnRlc3RAMTkyLjE2OC4xMDAuMTo4ODg4#example-server",
			"ss://bf-cfb:test@192.168.100.1:8888",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "192.168.101.101",
				Port:     8889,
				Method:   "bf-cfb",
				Password: "some_password",
				Tag:      "",
				PlainURI: "ss://bf-cfb:some_password@192.168.101.101:8889",
			},
			"ss://YmYtY2ZiOnNvbWVfcGFzc3dvcmRAMTkyLjE2OC4xMDEuMTAxOjg4ODk",
			"ss://bf-cfb:some_password@192.168.101.101:8889",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "192.168.102.3",
				Port:     443,
				Method:   "bf-cfb",
				Password: "some_password@?",
				Tag:      "",
				PlainURI: "ss://bf-cfb:some_password@?@192.168.102.3:443",
			},
			"ss://YmYtY2ZiOnNvbWVfcGFzc3dvcmRAP0AxOTIuMTY4LjEwMi4zOjQ0Mw",
			"ss://bf-cfb:some_password@?@192.168.102.3:443",
		},
		{
			ss.ShadowsocksURI{
				Hostname: "some.host",
				Port:     8888,
				Method:   "bf-cfb",
				Password: "test/!@#:",
				Tag:      "",
				PlainURI: "ss://bf-cfb:test/!@#:@some.host:8888",
			},
			"ss://YmYtY2ZiOnRlc3QvIUAjOkBzb21lLmhvc3Q6ODg4OA",
			"ss://bf-cfb:test/!@#:@some.host:8888",
		},
	}

	for i, ut := range tests {
		ssconf, err := ss.DecodeBase64URI(ut.b64URI)
		if err != nil {
			t.Errorf("#%d test failed. DecodeBase64URI() failed", i)
			continue
		}

		if !configIsEqual(*ssconf, ut.expectedConfig, true) {
			t.Errorf("#%d test failed. Expected %v, Got %v", i, ut.expectedConfig, *ssconf)
		}

		ssconf1, err := ss.DecodePlainURI(ut.plainURI)
		if err != nil {
			t.Errorf("#%d test failed. DecodeBase64URI() failed", i)
			continue
		}

		if !configIsEqual(*ssconf1, ut.expectedConfig, false) {
			t.Errorf("#%d test failed. Expected %v, Got %v", i, ut.expectedConfig, *ssconf1)
		}
	}
}

func configIsEqual(conf1 ss.ShadowsocksURI, conf2 ss.ShadowsocksURI, checkTag bool) bool {
	if conf1.Hostname != conf2.Hostname {
		return false
	}

	if conf1.Port != conf2.Port {
		return false
	}

	if conf1.Password != conf2.Password {
		return false
	}

	if conf1.Method != conf2.Method {
		return false
	}

	if conf1.Tag != conf2.Tag && checkTag {
		return false
	}

	if conf1.PlainURI != conf2.PlainURI {
		return false
	}

	return true
}
