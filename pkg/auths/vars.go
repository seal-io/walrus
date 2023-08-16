package auths

import (
	"time"

	"github.com/seal-io/walrus/utils/cryptox"
	"github.com/seal-io/walrus/utils/vars"
)

var (
	// MaxIdleDurationConfig holds the config of the max idle duration.
	MaxIdleDurationConfig = vars.SetOnce[time.Duration]{}

	// SecureConfig holds the config of securing.
	SecureConfig = vars.SetOnce[bool]{}

	// EncryptorConfig holds the config of the token encryptor.
	EncryptorConfig = vars.NewSetOnce[cryptox.Encryptor](cryptox.Null())
)
