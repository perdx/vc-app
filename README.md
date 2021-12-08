# Verifiable Credentials Backend

## Quick Start

Install Microsoft Authenticator on your phone.
Install (and configure) [ngrok](https://ngrok.com/download).

Run the app locally (in debug mode or `go run .`).

Use `ngrok` to create an external URL for callbacks:

    ngrok http 8080

Copy the `ngrok` HTTPS URL and use in Postman or React to test the end to end process locally.

## Infrastructure Setup

Use the Azure CLI scripts located in the `/.scripts/azure.sh`

Otherwise, for a step by step guide, see [Issue Verifiable Credentials from Azure Active Directory](https://docs.microsoft.com/en-nz/azure/active-directory/verifiable-credentials/enable-your-tenant-verifiable-credentials#update-the-sample-app).
