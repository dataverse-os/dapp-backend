#!/bin/sh
NETWORK=$1
# PUBKEY=$2

if [ -z "$NETWORK" ]
then
      echo "NETWORK is empty, using testnet-clay! \n"
      NETWORK='testnet-clay'
else
      echo "using $NETWORK! \n"
fi


PUBKEY="did:key:z6MkiM1beKfKoNAS5cqHTFMrWAqqHkdb7meMqMBurDDgnTRn#z6MkiM1beKfKoNAS5cqHTFMrWAqqHkdb7meMqMBurDDgnTRn"
if [ -z "$PUBKEY" ]
then
      echo "public key is empty! \n"
else
      cat daemon.config.template.json | jq  .'"http-api"'.'"admin-dids" += ["'$PUBKEY'"]' | jq  '.network.name = "mainnet"'  > daemon.config.json
      echo "deamon.config.json done! \n"
fi



