<br/>
<p align="center">
<a href=" " target="_blank">
<img src="https://raw.githubusercontent.com/dataverse-os/create-dataverse-app/main/logo.svg" width="180" alt="Dataverse logo">
</a >
</p >
<br/>

# dapp-backend

Scripts to help deploy a ceramic node with dataverse backend.

## Requirements
- Git
- Docker [Docs](https://docs.docker.com/get-docker/)
- docker-compose [Docs](https://docs.docker.com/compose/install/)
- node (version >= 16) [Docs](https://nodejs.org/en/download/)

## Installation
Clone the repo

```bash
git clone https://github.com/dataverse-os/dapp-backend
cd dapp-backend
```

## Getting started

### Generate private key and Admin DID

```bash
npm install -g @composedb/cli
composedb did:generate-private-key
```
**Example Result**
```bash
✔ Generating random private key... Done!
8053f2d22cb3da5f84b6f079eb40cdc49958a7da269de3610e63e8b8078f1448
```

Keep your private key safe. You will need it to use the ceramic node.

generate DID from private key
```bash
composedb did:from-private-key 8053f2d22cb3da5f84b6f079eb40cdc49958a7da269de3610e63e8b8078f1448
```

**Example Result**
```bash
✔ Creating DID... Done!
did:key:z6MkiM1beKfKoNAS5cqHTFMrWAqqHkdb7meMqMBurDDgnTRn
```

### Config your ceramic node
modify ```daemon.config.json``` to include your admin DID.
```json
{
  ...
  "http-api": {
    "cors-allowed-origins": [".*"],
    "admin-dids": ["Your Admin DID here"]
  },
  ...
}

```

copy the config to ceramic config folder

under ```dapp-backend/```
```bash
mkdir ~/.ceramic
cp ./daemon.config.json ~/.ceramic/daemon.config.json
```

### Update docker-compose.yml

modify the docker-compose.yml file to include your private key.
```YAML
version: "3.9"
services:
  ceramic:
    image: ceramicnetwork/js-ceramic:latest
    volumes:
      - ~/.ceramic/:/root/.ceramic/
    healthcheck:
      test:
        [
          "CMD-SHELL",
          "curl -f http://localhost:7007/api/v0/node/healthcheck || exit 1"
        ]
      interval: 1m30s
      timeout: 10s
      retries: 3
      start_period: 40s

  dapp-backend:
    image: dataverseos/dapp-backend:latest
    environment:
      - DID_PRIVATE_KEY={YOUR_PRIVATE_KEY}
      - CERAMIC_URL=http://ceramic:7007
    depends_on:
      - ceramic


```

### Run ceramic node

```bash
docker-compose up -d
```

## Use the ceramic node to create dataverse apps

view details in the [create-dataverse-app docs](https://github.com/dataverse-os/create-dataverse-app#readme).


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

## Contributing

Contributions are always welcome! Open a PR or an issue!