package config

import "github.com/seal-io/walrus/utils/vars"

// TlsCertified indicates whether the server is TLS certified.
var TlsCertified = vars.SetOnce[bool]{}
