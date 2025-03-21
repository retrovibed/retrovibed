package jwtx

import (
	"crypto/rand"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/ssh"
)

func NewSSHSigner() jwt.SigningMethod {
	return jwtsigner{}
}

type jwtsigner struct{}

func (t jwtsigner) Verify(signingString string, signature []byte, key any) error {
	var (
		err    error
		sigb   []byte
		ok     bool
		pubkey ssh.PublicKey
		sig    ssh.Signature
	)

	if pubkey, ok = key.(ssh.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	// Decode the signature
	if sigb, err = DecodeSegment(signature); err != nil {
		return err
	}

	if err = ssh.Unmarshal(sigb, &sig); err != nil {
		return err
	}

	if err = pubkey.Verify([]byte(signingString), &sig); err != nil {
		return err
	}

	return nil
}

func DecodeSegment(signature []byte) ([]byte, error) {
	panic("unimplemented")
}

func (t jwtsigner) Sign(signingString string, key interface{}) ([]byte, error) {
	var (
		s    ssh.Signer
		sigb []byte
		ok   bool
	)

	if s, ok = key.(ssh.Signer); !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	// Sign the string and return the encoded result
	sig, err := s.Sign(rand.Reader, []byte(signingString))
	if err != nil {
		return nil, err
	}

	sigb = ssh.Marshal(sig)

	return sigb, nil
}

func (t jwtsigner) Alg() string {
	return "ssh"
}
