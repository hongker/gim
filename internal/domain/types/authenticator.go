package types

import (
	"context"
	"encoding/base64"
)

var (
	DefaultAuthenticator = NewBase64Authenticator
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (uid string, err error)
	GenerateToken(ctx context.Context, uid string) (token string, err error)
}

func NewBase64Authenticator() *Base64Authenticator {
	return &Base64Authenticator{encoder: base64.StdEncoding}
}

type Base64Authenticator struct {
	encoder *base64.Encoding
}

func (auth Base64Authenticator) Authenticate(ctx context.Context, token string) (uid string, err error) {
	bytes, err := auth.encoder.DecodeString(token)
	return string(bytes), err
}

func (auth Base64Authenticator) GenerateToken(ctx context.Context, uid string) (token string, err error) {
	token = auth.encoder.EncodeToString([]byte(uid))
	return
}
