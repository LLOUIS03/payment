# Payment API 

## Description

### Payment API
  
This API provides functionality for processing payments. It allows users to make payments, 
retrieve payment details, and perform other payment-related operations. The API supports 
various payment methods and integrates with external payment gateways.

## Table of Contents

### Installation

Once you have golang [installed go][golang-install]

Once you have [installed Docker][docker-install]

### Running payment api

`make docker-up` running the payment app at "localhost:8090/swagger/index.html"
`make docker-down` stop the payment app at "localhost:8090/swagger/index.html"
`make goose` generate new migrations
`make sqlc` generate the repos

[golang-install]:   http://golang.org/doc/install.html#releases
[docker-install]:   https://docs.docker.com/engine/install/
