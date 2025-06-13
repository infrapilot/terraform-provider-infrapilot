// Copyright (c) 2025 InfraPilot, LLC
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LicenseClaims struct {
	OrgID      string    `json:"org_id"`
	Product    string    `json:"product"`
	Tier       string    `json:"tier"`
	ExpiresAt  time.Time `json:"expires_at"`
	GraceUntil time.Time `json:"grace_until"`
	jwt.RegisteredClaims
}

type LicenseMetadata struct {
	OrgID  types.String `tfsdk:"org_id"`
	Tier   types.String `tfsdk:"tier"`
	Module types.String `tfsdk:"module"`
}
