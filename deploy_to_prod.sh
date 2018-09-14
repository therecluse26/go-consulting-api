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
elif [ "$platform" == 'Linux' ]; then

    if [ $(which jq) == "" ]; then

        apt install jq

    fi
else

    echo "Must be run from a Unix-based bash shell. If you are on Windows, please install the Ubuntu package and try again from there."

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

export CFGJSON=$(echo $ENV_VARS | base64 -D)

# az acr login --name $AZUSER

az container create \
    --location $LOCATION  \
    --resource-group $RES_GROUP \
    --name $CONT_GROUP_NAME \
    --image $CONT_IMG \
    --cpu 1 --memory 1 \
    --dns-name-label $DNSLBL\
    --environment-variables CFGJSON="$CFGJSON" \
    --ports 80 443

EXIT_STATUS=$?

if [ $EXIT_STATUS == 0 ]; then
  echo "Container creation succeeded"
else 
  echo "Container creation failed"
fi
