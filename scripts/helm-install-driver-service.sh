#!/bin/bash

helm install driver-service ./helm/driver-service -f ./helm/driver-service/values.yaml -f ./helm/driver-service/values-secret.yaml