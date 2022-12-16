#!/bin/sh

###############################################################################################
##               *** WARNING - INSECURE - DO NOT USE IN PRODUCTION ***                       ##
## This script is to simulate operations a Vault operator would perform and, as such,        ##
## is not a representation of best practices in production environments.                     ##
## https://learn.hashicorp.com/tutorials/vault/pattern-approle?in=vault/recommended-patterns ##
###############################################################################################

set -e

export VAULT_ADDR='http://localhost:8200'
export VAULT_FORMAT='json'

# Spawn a new process for the development Vault server and wait for it to come online
# ref: https://www.vaultproject.io/docs/concepts/dev-server
vault server -dev -dev-listen-address="0.0.0.0:8200" &
sleep 5s

# Authenticate container's local Vault CLI
# ref: https://www.vaultproject.io/docs/commands/login
vault login -no-print "${VAULT_DEV_ROOT_TOKEN_ID}"

#####################################
########## ACCESS POLICIES ##########
#####################################

# Add policies for the various roles we'll be using
# ref: https://www.vaultproject.io/docs/concepts/policies
#vault policy write dev-policy /vault/config/dev-policy.hcl
vault policy write dev-policy /vault/config/dev-policy.hcl

#####################################
######## APPROLE AUTH METHDO ########
#####################################

# Enable AppRole auth method utilized by our web application
# ref: https://www.vaultproject.io/docs/auth/approle
vault auth enable approle

# Configure a specific AppRole role with associated parameters
# ref: https://www.vaultproject.io/api/auth/approle#parameters
#
# NOTE: we use artificially low ttl values to demonstrate the credential renewal logic
vault write auth/approle/role/dev-role \
    token_policies=dev-policy \
    secret_id_ttl="24h" \
    token_ttl="24h"

# Overwrite our role id with a known value to simplify our demo
vault write auth/approle/role/dev-role/role-id role_id="${APPROLE_ROLE_ID}"

#####################################
########## STATIC SECRETS ###########
#####################################

# Enable the kv-v2 secrets engine, passing in the path parameter
# ref: https://www.vaultproject.io/docs/secrets/kv/kv-v2
vault secrets enable -path=tinyedge kv-v2

## PKI certificates

vault secrets enable pki

vault write -field=certificate pki/root/generate/internal \
    common_name="home.net" \
    issuer_name="root-2022" \
    ttl=87600h

issuer=$(vault list pki/issuers/ | jq '.[0]' | sed 's/"//g')

# create role for CA
vault write pki/roles/tinyedge allow_any_name=true
vault write pki/config/urls \
    issuing_certificates="$VAULT_ADDR/v1/pki/ca" \
    crl_distribution_points="$VAULT_ADDR/v1/pki/crl"

# generate intermediate CA
vault secrets enable -path=pki_int pki

vault secrets tune -max-lease-ttl=43800h pki_int

vault write -format=json pki_int/intermediate/generate/internal \
    common_name="home.net Intermediate Authority" \
    issuer_name="tinyedge-intermediate" \
    | jq -r '.data.csr' > pki_intermediate.csr

vault write -format=json pki/root/sign-intermediate \
    issuer_name="root-2022" \
    csr=@pki_intermediate.csr \
    format=pem_bundle ttl="43800h" \
    | jq -r '.data.certificate' > intermediate.cert.pem

# import ca intermediate to vault
vault write pki_int/intermediate/set-signed certificate=@intermediate.cert.pem

issuer_ref=$(vault read -field=default pki_int/config/issuers | sed 's/"//g')
# create role
vault write pki_int/roles/tinyedge-role \
    issuer_ref="$issuer_ref" \
    allowed_domains="home.net" \
    allow_subdomains=true \
    max_ttl="720h"
# This container is now healthy
touch /tmp/healthy

# Keep container alive
tail -f /dev/null & trap 'kill %1' TERM ; wait
