// Copyright (c) 2025 InfraPilot, LLC
// SPDX-License-Identifier: MPL-2.0

package jwt

import (
	"fmt"
	"net/http"
	"time"

	"terraform-provider-infrapilot/internal/model"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

// ValidateToken fetches the JWKS from the license server and validates the token.
func ValidateToken(tokenStr string, jwksURL string) (*model.LicenseClaims, error) {
	// Setup JWKS options (note: no DisableValidation field here)
	options := keyfunc.Options{
		RefreshErrorHandler: func(err error) {
			fmt.Printf("JWKS refresh error: %v\n", err)
		},
		RefreshInterval: time.Hour,
		RefreshTimeout:  5 * time.Second,
		Client:          http.DefaultClient,
	}

	jwks, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to load JWKS: %w", err)
	}

	// Use jwt.NewParser to disable automatic validation of claims like 'exp'
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())

	token, err := parser.ParseWithClaims(tokenStr, &model.LicenseClaims{}, jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("token parse error: %w", err)
	}

	claims, ok := token.Claims.(*model.LicenseClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid or malformed token claims")
	}

	if claims.RegisteredClaims.ExpiresAt == nil {
		return nil, fmt.Errorf("token missing 'exp' field")
	}
	if time.Now().After(claims.RegisteredClaims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired at %s", claims.RegisteredClaims.ExpiresAt.Format(time.RFC3339))
	}

	return claims, nil
}
