#!/bin/bash

kubectl create configmap traefik-dynamic-conf  --from-file=middlewares.yaml= ./helm/traefik/middlewares.yaml

helm install traefik traefik/traefik -f ./helm/traefik/values.yaml
