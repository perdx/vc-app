# AT Verifiable Credentials Backend

## Quick Start

Install Microsoft Authenticator on your phone.
Install (and configure) [ngrok](https://ngrok.com/download).

Run the app locally (in debug mode or `go run .`).

Use `ngrok` to create an external URL for callbacks:

    ngrok http 8080

Copy the `ngrok` HTTPS URL and use in Postman or React to test the end to end process locally.

## Infrastructure Setup

For a step by step guide, see [Issue Verifiable Credentials from Azure Active Directory](https://docs.microsoft.com/en-nz/azure/active-directory/verifiable-credentials/enable-your-tenant-verifiable-credentials#update-the-sample-app).

### Verifiable Credential Request Service

Create a service principal for the Verifiable Credential Request Service API. `bbb94529-53a3-4be5-a069-7eaf2712b826` is the service's app id.

    az ad sp create --id "bbb94529-53a3-4be5-a069-7eaf2712b826"

### Key Vault

The key vault is used to store the issuer private keys (signing key and recovery key).

1. Add `Sign` permission to the admin user.
2. Add an access policy for VC service principal created above. Grant `Key Permissions/Get`, `Key Permissions/Sign`, and `Secret Permissions/Get`

### External Application Access

Register an application with API permission to the `Verifiable Credential Request Service` so that apps can get access tokens.
Use Azure AD/App registrations to create an app with `API permissions` to the `Verifiable Credential Request Service`.

### Add a Verifiable Credential

Search for verifiable credentials in Azure Portal.

Add an organisation name, domain, and select the key vault created earlier.
Add a credential with a `Display file` (claims) and `Rules file` (rules to get credential).

