#!/bin/bash

helm install trip-service ./helm/trip-service -f ./helm/trip-service/values.yaml -f ./helm/trip-service/values-secret.yaml