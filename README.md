## SSURI

> A command line tool for manipulating Shadowsocks URI and configuration.

### Installation

```sh
$ go get github.com/vgxbj/ssuri/cmd/ssuri
```

### Usage

```
Usage: ssuri [-h] [-i in_file] [-o out_file]
  -dump-uri
        dump shadowsocks URI
  -generate-json-config
        generate JSON configurations
  -generate-qr
        generate QR code
  -i string
        input file (default: "-" for stdin) (default "-")
  -json
        read JSON as input (default: off)
  -legacy
        dump shadowsocks URI in legacy mode (default: off)
  -o string
        output file (default: "-" for stdout) (default "-")
```

### Example

- Read shadowsocks URI, dump it, generate QR code and JSON configuration.

```sh
$ echo "ss://YmYtY2ZiOnRlc3Q=@192.168.100.1:8888" | ssuri -dump-uri -generate-json-config -generate-qr
```

- Read JSON configuration, generate shadowsocks URI and QR code.

```
$ ssuri -i config.json -generate-json-config -generate-uri
```

### Features

- [x] Support SIP002 URI scheme and legacy base64 encoded URI scheme.
- [x] Read JSON configuration or URI. Generate JSON configuration, URI or QR code.

### TODO

- [ ] Support manipulating shadowsocksR URI or configuration.
