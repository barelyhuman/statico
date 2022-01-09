#!/bin/bash

curl -fsSL https://bina.egoist.sh/barelyhuman/statico | sh
statico -c doc.statico.yml 
# Make sure you have node and npx installed
npx purgecss --config purgecss.config.js