# TODO:
# 1. Deployment service principal for Github Actions
# 2. Container App

# https://docs.microsoft.com/en-us/azure/active-directory/verifiable-credentials/verifiable-credentials-configure-tenant

# --------------------------------------------------
# Environment
# --------------------------------------------------

# env vars
SUBSCRIPTION="<subscription>"
RESOURCE_GROUP=rg-vc
LOCATION=australiaeast
VC_APP_ID=bbb94529-53a3-4be5-a069-7eaf2712b826
VC_IMAGE=ghcr.io/perdx/vc-app:latest

# secrets and variables
GHCR_USER="<github machine user>"
GHCR_PWD="<github PAT>"

# service names
CA_ENV_NAME="<container app environment name>"
CA_NAME="<container app name>"

az account set --subscription $SUBSCRIPTION

# --------------------------------------------------
# Resource Group(s)
# --------------------------------------------------
az group create --name $RESOURCE_GROUP --location $LOCATION

# --------------------------------------------------
# Service Principals
# --------------------------------------------------
# verifiable credential service principal 
# - adds Verifiable Creditential service to this tenant
az ad sp create --id $VC_APP_ID

# app (our backend) service principal
az ad sp create-for-rbac \
    --name "sp-deployment" \
    --role contributor \
    --scopes /subscriptions/$SUBSCRIPTION/resourceGroups/$RESOURCE_GROUP \
    --sdk-auth

# N.B. save result for use in Github Actions

# --------------------------------------------------
# Key Vault
#  - used by Verifiable Creditential service
# --------------------------------------------------
KEY_VAULT=kv-vc
VC_REQUEST_OBJID=e81217f2-7480-4413-b672-88ba9a27d40d
VC_ISSUER_OBJID=70d7e704-cc68-4da7-9ad1-0b7c0bcb05fa

# create key vault
az keyvault create \
    --name $KEY_VAULT \
    --resource-group $RESOURCE_GROUP \
    --location $LOCATION

# Verifiable Credentials Request Service permissions
az keyvault set-policy \
    --name $KEY_VAULT \
    --resource-group $RESOURCE_GROUP \
    --object-id $VC_REQUEST_OBJID \
    --key-permissions get sign \
    --secret-permissions get

# Verifiable Credentials Issuer Service permissions
az keyvault set-policy \
    --name $KEY_VAULT \
    --resource-group $RESOURCE_GROUP \
    --object-id $VC_ISSUER_OBJID \
    --key-permissions get sign \
    --secret-permissions get

# az keyvault purge --subscription Innovation -n kv-perdx

# --------------------------------------------------
# Storage
#  - used by Verifiable Creditential service
# --------------------------------------------------

az storage account create \
    --name stverifiablecredentials \ 
    --resource-group $RESOURCE_GROUP \  
    --location $LOCATION \
    --sku Standard_LRS

# --------------------------------------------------
# App Registration
#  - used by Go app to connect to Verifiable Credential service 
# --------------------------------------------------
API_APP_NAME="Verifiable Credentials API"

# add app registration - used by our backend to access verifiable credentials
az ad app create \
    --display-name $API_APP_NAME \
    --available-to-other-tenants false

# add API permissions for Verifiable Credentials Request Service
API_APP_ID=$(az ad app list --display-name $API_APP_NAME --query '[].appId' -o tsv)
ACCESS_ID=$(az ad sp show --id $VC_APP_ID --query "appRoles[?displayName == 'VerifiableCredential.Create.All'].id" -o tsv)

az ad app permission add \
    --id $API_APP_ID \
    --api $VC_APP_ID \
    --api-permissions $ACCESS_ID=Role

az ad app permission admin-consent --id $API_APP_ID

# create secret ($API_APP_SECRET)
az ad app credential reset \
    --id $API_APP_ID \
    --credential-description "App Authn" \
    --years 1 \
    --append

# --------------------------------------------------
# Configure Verifiable Credentials
#  - N.B. no az cli commands available for now
# --------------------------------------------------

# Search for verifiable credentials in Azure Portal
# 1. Create new credential
# 2. Link to display.json and rules.json (in storage container)

# --------------------------------------------------
# Configure Azure CLI for Container Apps
# --------------------------------------------------

# install Container Apps extension
az extension add --source https://workerappscliextension.blob.core.windows.net/azure-cli-extension/containerapp-0.2.0-py2.py3-none-any.whl

# register Microsoft.Web provider
az provider register --namespace Microsoft.Web

# --------------------------------------------------
# Container App Environment
# --------------------------------------------------

LA_NAME=vc-logs
CA_LOCATION=canadacentral

# logging
az monitor log-analytics workspace create \
  --resource-group $RESOURCE_GROUP \
  --workspace-name $LA_NAME

# container app environment
LA_CLIENT_ID=`az monitor log-analytics workspace show --query customerId -g $RESOURCE_GROUP -n $LA_NAME --out tsv`
LA_CLIENT_SECRET=`az monitor log-analytics workspace get-shared-keys --query primarySharedKey -g $RESOURCE_GROUP -n $LA_NAME --out tsv`

az containerapp env create \
  --location $CA_LOCATION \
  --resource-group $RESOURCE_GROUP \
  --name $CA_ENV_NAME \
  --logs-workspace-id $LA_CLIENT_ID \
  --logs-workspace-key $LA_CLIENT_SECRET

# --------------------------------------------------
# Container App
# --------------------------------------------------

TENANT=9125264c-86cb-45fe-baa2-e022db0590d6
AUTHORITY=did:ion:EiBuwQM4Yu-r3NV2qQsaeu2ziZ03D4TUTKRCAZjDVVteIg:eyJkZWx0YSI6eyJwYXRjaGVzIjpbeyJhY3Rpb24iOiJyZXBsYWNlIiwiZG9jdW1lbnQiOnsicHVibGljS2V5cyI6W3siaWQiOiJzaWdfM2ZlZjk4ZDQiLCJwdWJsaWNLZXlKd2siOnsiY3J2Ijoic2VjcDI1NmsxIiwia3R5IjoiRUMiLCJ4IjoiSjZDeEE5U2QzeUV4Z2hTTDJ6OUx0YzYzMXZxbEJfRFV6bEo4QlU3WWZORSIsInkiOiJwWmhYRG9LbVNNc2FlcnY4N3V3ME5zOWZLZkFEN1hLZmZQMjJreXBPVXZNIn0sInB1cnBvc2VzIjpbImF1dGhlbnRpY2F0aW9uIiwiYXNzZXJ0aW9uTWV0aG9kIl0sInR5cGUiOiJFY2RzYVNlY3AyNTZrMVZlcmlmaWNhdGlvbktleTIwMTkifV0sInNlcnZpY2VzIjpbeyJpZCI6ImxpbmtlZGRvbWFpbnMiLCJzZXJ2aWNlRW5kcG9pbnQiOnsib3JpZ2lucyI6WyJodHRwczovL3BlcmR4LmlvLyJdfSwidHlwZSI6IkxpbmtlZERvbWFpbnMifV19fV0sInVwZGF0ZUNvbW1pdG1lbnQiOiJFaUJJeTFTTXhzWjRwLS1uNkI1MVRXQjFDZWl4bDZqa3V6QVNUZmc2QWF4bFdBIn0sInN1ZmZpeERhdGEiOnsiZGVsdGFIYXNoIjoiRWlEcjlUbVlhSTEwamEtMFpzWkJ5ODJuWmcyT2F1MGpSbUFWSHE3Z09renBRQSIsInJlY292ZXJ5Q29tbWl0bWVudCI6IkVpQlJldTY5RXN0UVRLeDVCV1VfUzhyNDRoXzZpY20tVlEyWXZTNU54NkwtV3cifX0
CLIENT_ID=$API_APP_ID
CLIENT_SECRET="<API_APP_SECRET from above>"

az containerapp create \
  --resource-group $RESOURCE_GROUP \
  --name $CA_NAME \
  --environment $CA_ENV_NAME \
  --registry-login-server ghcr.io \
  --registry-username "$GHCR_USER" \
  --registry-password "$GHCR_PWD" \
  --image $VC_IMAGE \
  --cpu 0.5 --memory 1.0Gi \
  --max-replicas 5 \
  --ingress external \
  --target-port 8080 \
  --secrets "client-id=$CLIENT_ID,client-secret=$CLIENT_SECRET" \
  --environment-variables \
    "TENANT=$TENANT,AUTHORITY=$AUTHORITY,CLIENT_ID=secretref:client-id,CLIENT_SECRET=secretref:client-secret"