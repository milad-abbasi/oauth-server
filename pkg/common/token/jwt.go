package token

import (
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type JWT struct {
	ID            string
	Issuer        string
	Subject       string
	Audience      []string
	Expiry        time.Duration
	NotBefore     time.Time
	IssuedAt      time.Time
	PrivateClaims []interface{}
}

type Expectation struct {
	Secret   string
	ID       string
	Issuer   string
	Subject  string
	Audience []string
	Time     time.Time
}

func (j *JWT) Sign(secret string) (string, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte(secret)},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return "", err
	}

	publicClaims := jwt.Claims{
		ID:        j.ID,
		Issuer:    j.Issuer,
		Subject:   j.Subject,
		Audience:  jwt.Audience(j.Audience),
		Expiry:    jwt.NewNumericDate(time.Now().Add(j.Expiry)),
		NotBefore: jwt.NewNumericDate(j.NotBefore),
		IssuedAt:  jwt.NewNumericDate(j.IssuedAt),
	}

	tokenBuilder := jwt.Signed(signer).Claims(publicClaims)
	for _, claim := range j.PrivateClaims {
		tokenBuilder = tokenBuilder.Claims(claim)
	}

	signedToken, err := tokenBuilder.CompactSerialize()
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWT) Encrypt(secret string) (string, error) {
	encryption, err := jose.NewEncrypter(
		jose.A128CBC_HS256,
		jose.Recipient{Algorithm: jose.PBES2_HS256_A128KW, Key: []byte(secret)},
		(&jose.EncrypterOptions{}).WithType("JWT"),
	)
	if err != nil {
		return "", err
	}

	publicClaims := jwt.Claims{
		ID:        j.ID,
		Issuer:    j.Issuer,
		Subject:   j.Subject,
		Audience:  jwt.Audience(j.Audience),
		Expiry:    jwt.NewNumericDate(time.Now().Add(j.Expiry)),
		NotBefore: jwt.NewNumericDate(j.NotBefore),
		IssuedAt:  jwt.NewNumericDate(j.IssuedAt),
	}

	tokenBuilder := jwt.Encrypted(encryption).Claims(publicClaims)
	for _, claim := range j.PrivateClaims {
		tokenBuilder = tokenBuilder.Claims(claim)
	}

	encryptedToken, err := tokenBuilder.CompactSerialize()
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}

func ValidateSignedJwt(signedToken string, e *Expectation, publicClaims *jwt.Claims, privateClaims ...interface{}) (bool, error) {
	token, err := jwt.ParseSigned(signedToken)
	if err != nil {
		return false, err
	}

	claims := []interface{}{publicClaims}
	claims = append(claims, privateClaims...)
	if err := token.Claims([]byte(e.Secret), claims...); err != nil {
		return false, err
	}

	err = publicClaims.Validate(jwt.Expected{
		ID:       e.ID,
		Issuer:   e.Issuer,
		Subject:  e.Subject,
		Audience: jwt.Audience(e.Audience),
		Time:     e.Time,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func ValidateEncryptedJwt(signedToken string, e *Expectation, publicClaims *jwt.Claims, privateClaims ...interface{}) (bool, error) {
	token, err := jwt.ParseEncrypted(signedToken)
	if err != nil {
		return false, err
	}

	claims := []interface{}{publicClaims}
	claims = append(claims, privateClaims...)
	if err := token.Claims([]byte(e.Secret), claims...); err != nil {
		return false, err
	}

	err = publicClaims.Validate(jwt.Expected{
		ID:       e.ID,
		Issuer:   e.Issuer,
		Subject:  e.Subject,
		Audience: jwt.Audience(e.Audience),
		Time:     e.Time,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func ExtractBearerJwt(token string) string {
	return strings.TrimPrefix(token, "Bearer ")
}
