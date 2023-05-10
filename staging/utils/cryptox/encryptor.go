package cryptox

// Encryptor holds the operations of encryption.
type Encryptor interface {
	// Encrypt obtains the ciphertext bytes with the given bytes,
	// additional data can be nil.
	Encrypt(plaintext []byte, additional []byte) (ciphertext []byte, err error)

	// Decrypt obtains the plaintext bytes with the given bytes,
	// additional data can be nil.
	Decrypt(ciphertext []byte, additional []byte) (plaintext []byte, err error)
}

// Null returns an Encryptor with nothing to do.
func Null() Encryptor {
	return nullEncryptor{}
}

// nullEncryptor implements Encryptor but does nothing,
// always returns the plaintext bytes.
type nullEncryptor struct{}

func (nullEncryptor) Encrypt(p []byte, _ []byte) ([]byte, error) {
	if p == nil {
		return p, nil
	}
	c := make([]byte, len(p))
	copy(c, p)
	return c, nil
}

func (nullEncryptor) Decrypt(p []byte, _ []byte) ([]byte, error) {
	if p == nil {
		return p, nil
	}
	c := make([]byte, len(p))
	copy(c, p)
	return c, nil
}
