// Copyright (c) 2025 InfraPilot, LLC
// SPDX-License-Identifier: MPL-2.0

package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"terraform-provider-infrapilot/internal/jwt"
	"terraform-provider-infrapilot/internal/model"

	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

// createTestJWKS sets up a temporary JWKS server and returns the server URL and a token signed with the key.
func createTestJWKS(t *testing.T) (string, string) {
	t.Helper()

	// Generate RSA key pair
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pubKey := privKey.PublicKey
	n := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString([]byte{0x01, 0x00, 0x01}) // 65537

	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"kid": "test-key",
				"use": "sig",
				"alg": "RS256",
				"n":   n,
				"e":   e,
			},
		},
	}

	jwksBytes, err := json.Marshal(jwks)
	require.NoError(t, err)

	jwksServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(jwksBytes)
	}))

	token := jwtv4.NewWithClaims(jwtv4.SigningMethodRS256, model.LicenseClaims{
		OrgID:   "test-org",
		Product: "infra-pilot",
		Tier:    "pro",
		RegisteredClaims: jwtv4.RegisteredClaims{
			ExpiresAt: jwtv4.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwtv4.NewNumericDate(time.Now()),
		},
	})
	token.Header["kid"] = "test-key"

	signedToken, err := token.SignedString(privKey)
	require.NoError(t, err)

	return jwksServer.URL, signedToken
}

func TestValidateToken_Success(t *testing.T) {
	jwksURL, token := createTestJWKS(t)

	claims, err := jwt.ValidateToken(token, jwksURL)
	require.NoError(t, err)
	require.Equal(t, "test-org", claims.OrgID)
	require.Equal(t, "infra-pilot", claims.Product)
	require.Equal(t, "pro", claims.Tier)
}
