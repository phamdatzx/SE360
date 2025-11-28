#!/bin/bash

helm install traefik traefik/traefik -f helm/traefik/values.yaml

kubectl apply -f helm/traefik/middlewares-auth.yaml

kubectl apply -f helm/traefik/middlewares-cors.yaml