# dapp-backend

Scripts to help deploy a ceramic node with dataverse backend.

## Install

### Docker

### Build from source

## Getting started

### Generate private key and Admin DID

```bash
npm install -g @composedb/cli
composedb did:generate-private-key
```
Example Result
```bash
✔ Generating random private key... Done!
8053f2d22cb3da5f84b6f079eb40cdc49958a7da269de3610e63e8b8078f1448
```

generate DID from private key
```bash
composedb did:from-private-key 8053f2d22cb3da5f84b6f079eb40cdc49958a7da269de3610e63e8b8078f1448
```

Example Result
```bash
✔ Creating DID... Done!
did:key:z6MkiM1beKfKoNAS5cqHTFMrWAqqHkdb7meMqMBurDDgnTRn
```

### Config your ceramic node

Modify the config file `~/.ceramic/config.json` to add the admin DID to the `admin-dids` array.

```json
{
  "anchor": {
  },
  "http-api": {
    "cors-allowed-origins": [
      ".*"
    ],
    "admin-dids": [
            "did:key:z6MkiM1beKfKoNAS5cqHTFMrWAqqHkdb7meMqMBurDDgnTRn#z6MkiM1beKfKoNAS5cqHTFMrWAqqHkdb7meMqMBurDDgnTRn",
    ]
  },
  "ipfs": {
    "mode": "bundled"
  },
  "logger": {
    "log-level": 2,
    "log-to-files": false
  },
  "metrics": {
    "metrics-exporter-enabled": false
  },
  "network": {
    "name": "mainnet"
  },
  "node": {},
  "state-store": {
    "mode": "fs",
    "local-directory": "~/.ceramic/statestore/"
  },
  "indexing": {
    "db": "sqlite://~/.ceramic/indexing.sqlite",
    "allow-queries-before-historical-sync": true
  }
}
```

### Run ceramic node

```bash
```

### [Optional] Run ceramic node on the mainnet

#### Verify your email address

```bash
curl --request POST \
  --url https://cas.3boxlabs.com/api/v0/auth/verification \
  --header 'Content-Type: application/json' \
  --data '{"email": "youremailaddress"}'
```
Then check your email and copy the one time passcode enclosed within. It will be a string of letters and numbers similar to this: 2451cc10-5a39-494d-b8eb-1971ecd813de.
#### Send a revocation request
```bash
 curl --request POST \
  --url https://cas.3boxlabs.com/api/v0/auth/did \
  --header 'Content-Type: application/json' \
  --data '{
    "email": "youremailaddress",
      "otp": "youronetimepasscode",
      "dids": [
          "yourdid"
      ]
  }'
  ```