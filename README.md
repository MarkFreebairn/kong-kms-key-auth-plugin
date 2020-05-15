# kong-kms-key-auth-plugin
A kong plugin written in go to provide key authentication with AWS kms

## Usage
Example of kong declarative config
```
_format_version: "1.1"
services:
- url: https://www.google.com
  routes:
  - paths:
    - "/test"
  plugins:
  - name: kms-key-auth
    config:
      encapikey: <encrypted kms blob>
      kmsregion: us-east-1
```

## Build
Build kong go plugins with the foillowing command:
```
go build --buildmode plugin
```