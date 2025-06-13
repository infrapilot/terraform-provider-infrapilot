# 0.1.0 (2025-06-12)

Initial release of the InfraPilot Terraform provider.

FEATURES:

- New `infrapilot` provider with support for JWT-based license validation
- `infrapilot_license_check` data source:
  - Validates token authenticity and expiration using public JWKS endpoint
  - Exposes computed attributes `org_id` and `tier`
  - Accepts an optional `module` attribute for module-specific feature control and future telemetry
- Environment variable support: `INFRAPILOT_TOKEN`
