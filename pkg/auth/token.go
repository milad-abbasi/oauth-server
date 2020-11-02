package auth

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type UserClaims struct {
	Email string `json:"email"`
}

func generateSignedToken(secret, userId, userEmail string) (string, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte(secret)},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return "", err
	}

	publicClaims := jwt.Claims{
		ID:        uuid.New().String(),
		Issuer:    "oauth-server",
		Subject:   userId,
		Audience:  jwt.Audience{"user"},
		Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	userClaims := UserClaims{
		Email: userEmail,
	}

	signedToken, err := jwt.Signed(signer).Claims(publicClaims).Claims(userClaims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateSignedToken(secret, signedToken string) (bool, error) {
	token, err := jwt.ParseSigned(signedToken)
	if err != nil {
		return false, err
	}

	publicClaims := jwt.Claims{}
	userClaims := UserClaims{}
	if err := token.Claims([]byte(secret), &publicClaims, &userClaims); err != nil {
		return false, err
	}

	err = publicClaims.Validate(jwt.Expected{
		Issuer: "oauth-server",
		// Subject:  "user-id",
		// Audience: jwt.Audience{"user"},
		// ID:       "",
		Time: time.Now(),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func generateEncryptedToken(secret, userId, userEmail string) (string, error) {
	encryption, err := jose.NewEncrypter(
		jose.A128CBC_HS256,
		jose.Recipient{Algorithm: jose.PBES2_HS256_A128KW, Key: []byte(secret)},
		(&jose.EncrypterOptions{}).WithType("JWT"),
	)
	if err != nil {
		return "", err
	}

	publicClaim := jwt.Claims{
		ID:        uuid.New().String(),
		Issuer:    "oauth-server",
		Subject:   userId,
		Audience:  jwt.Audience{"user"},
		Expiry:    jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	userClaims := UserClaims{
		Email: userEmail,
	}

	encryptedToken, err := jwt.Encrypted(encryption).Claims(publicClaim).Claims(userClaims).CompactSerialize()
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}

func validateEncryptedToken(secret, encryptedToken string) (bool, error) {
	token, err := jwt.ParseEncrypted(encryptedToken)
	if err != nil {
		return false, err
	}

	publicClaims := jwt.Claims{}
	userClaims := UserClaims{}
	if err := token.Claims([]byte(secret), &publicClaims, &userClaims); err != nil {
		return false, err
	}

	err = publicClaims.Validate(jwt.Expected{
		Issuer: "oauth-server",
		// Subject:  "user-id",
		// Audience: jwt.Audience{"user"},
		// ID:       "",
		Time: time.Now(),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
