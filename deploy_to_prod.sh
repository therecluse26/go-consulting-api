#!/bin/bash

platform=`uname`

if [ "$platform" == 'Darwin' ]; then

   if [ $(which jq) == "" ]; then

       if [ $(which brew) == ""]; then

            echo "Please install Homebrew and try again"
        else
             brew install jq
       fi
   fi
fi

AZUSER=FortisureAPI
CONT_GROUP_NAME=rest-api
CONT_IMG=fortisureapi.azurecr.io/rest-api-go
LOCATION=eastus
DNSLBL=fortisureapi
RES_GROUP=Container_Resources
ACR_NAME=fortisureapi.azurecr.io
AKV_NAME=FortisureKeys

echo "Environment Variable JSON (Base64 encoded):"
read ENV_VARS

# az acr login --name $AZUSER

az container create \
    --location $LOCATION  \
    --resource-group $RES_GROUP \
    --name $CONT_GROUP_NAME \
    --image $CONT_IMG \
    --cpu 1 --memory 1 \
    --dns-name-label $DNSLBL\
    --environment-variables gocfg64="$ENV_VARS" \
    --ports 80 443

EXIT_STATUS=$?

if [ $EXIT_STATUS == 0 ]; then
  echo "Container creation succeeded"
else 
  echo "Container creation failed"
fi
