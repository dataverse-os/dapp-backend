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

## Installation
Clone the repo

```bash
git clone https://github.com/dataverse-os/dapp-backend
cd dapp-backend
```

## Getting started

### Generate private key and Admin DID

To generate and fill the private key and admin DID into `.private-key.env` and `daemon.config.json` automatically:

under ```dapp-backend/```
```bash
docker run -it --rm \
  -v "$PWD/.private-key.env:/data/.private-key.env" \
  -v "$PWD/daemon.config.json:/data/daemon.config.json" \
  ghcr.io/dataverse-os/dapp-backend-prerun:latest
```
**Example Result**
```bash
✔ Generating random private key... Done!
8053f2d22cb3da5f84b6f079eb40cdc49958a7da269de3610e63e8b8078f1448
✔ Creating DID... Done!
did:key:z6MkiM1beKfKoNAS5cqHTFMrWAqqHkdb7meMqMBurDDgnTRn
```

Keep your private key safe. You will need it to use the ceramic node.

### Config your ceramic node
copy the config to ceramic config folder

under ```dapp-backend/```
```bash
mkdir ~/.ceramic
cp ./daemon.config.json ~/.ceramic/daemon.config.json
```

### Run ceramic node

```bash
docker-compose up -d
```
This command will start the ceramic node on port `7007` and the dapp-backend server on `8080`. You can change the port to use by changing the port in the docker-compose.yml file.


### Configure SSL certificate
To let your app connect to the ceramic node safely, you need to configure SSL certificate. You can use [Let's Encrypt](https://letsencrypt.org/) to get a free SSL certificate.

## Use the ceramic node to create dataverse apps

view details in the [create-dataverse-app docs](https://github.com/dataverse-os/create-dataverse-app#readme).


## [Optional] Run ceramic node on the mainnet

### Verify your email address

```bash
curl --request POST \
  --url https://cas.3boxlabs.com/api/v0/auth/verification \
  --header 'Content-Type: application/json' \
  --data '{"email": "youremailaddress"}'
```
Then check your email and copy the one time passcode enclosed within. It will be a string of letters and numbers similar to this: 2451cc10-5a39-494d-b8eb-1971ecd813de.
### Send a revocation request
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