#!/bin/bash

helm install user-service ./helm/user-service -f ./helm/user-service/values.yaml -f ./helm/user-service/values-secret.yaml