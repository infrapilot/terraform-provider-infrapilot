# InfraPilot Terraform Provider

The **InfraPilot Terraform Provider** enables licensed access to the InfraPilot module ecosystem, enforcing license validation and surfacing metadata like organization ID and subscription tier.

This provider is used internally by InfraPilot modules to verify usage rights, validate licenses, and enable tier-specific functionality. It acts as a **lightweight gateway** to the InfraPilot ecosystem and is required when consuming any InfraPilot-managed Terraform module.

---

## 🔐 Why InfraPilot?

InfraPilot offers a curated, enterprise-grade library of Terraform modules that:

* Follow best practices and opinionated design patterns
* Are versioned, tested, and continuously validated
* Include built-in CI/CD, IAM, and policy hooks
* Require license validation to ensure proper use

This provider validates a license token (JWT) issued by the InfraPilot license server and enables downstream modules to introspect tier metadata (e.g., `pro`, `enterprise`, etc.).

---

## 🚀 Getting Started

### Installation

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    infrapilot = {
      source  = "infra-pilot/infrapilot"
      version = ">= 0.1.0"
    }
  }
}

provider "infrapilot" {
  token = var.infrapilot_license_token
}
```

You can also set the license token via environment variable:

```bash
export INFRAPILOT_TOKEN="your-jwt-license-token"
```

---

### Example Usage

```hcl
data "infrapilot_license_check" "this" {
  module = "your_module_name"
}

output "org_id" {
  value = data.infrapilot_license_check.this.org_id
}

output "tier" {
  value = data.infrapilot_license_check.this.tier
}
```

---

## 🧪 Developing the Provider

### Requirements

* [Go](https://golang.org/dl/) 1.23+
* [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.0+

### Build Locally

```bash
go install
```

### Run Tests

```bash
make testacc
```

Tests use the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) and include unit and acceptance coverage.

> Note: Acceptance tests do not create cloud resources but do simulate full provider evaluation behavior.

### Generate Docs

```bash
make generate
```

This updates the `docs/` directory based on your schema definitions.

---

## 📦 Publishing

The provider is published to the [Terraform Registry](https://registry.terraform.io/providers/infra-pilot/infrapilot/latest). If you are building your own fork, update the `main.go` with your appropriate `ServeOpts.Address`.

---

## 📝 License

This project is licensed under the [Mozilla Public License 2.0 (MPL-2.0)](LICENSE). It incorporates portions of the HashiCorp provider scaffolding, also licensed under MPL-2.0.

---

## 🤝 Support

If you're an InfraPilot customer, please refer to your support plan or license documentation. For general issues or feature requests, open a GitHub issue in this repository.
