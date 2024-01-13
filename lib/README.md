# LIB

This is a collection of various packages which are used in more than one of our services.

## Auth-Middleware

The auth-middleware is used to export a controller which authenticates the current request with the user-service and blocks not authenticated requests.

## Client

Client is a package for an interface and it's mock for testing.

## containerhelpers

Containerhelpers provides a function to easily start a postgres testcontainer.

## database

The database-package provides the Postgres-Config for all our services.

## gprc

In the GRPC-Package are all protobuf files and generated interfaces.

## health

Here is the controller with the healthcheck-endpoint which all services use.

## router

The router package provides an http-router with middleware, which can match url of an incoming request and call the specified handler for this request. If a middleware for this url is specified, it will be called before the handler. The handler then only will be called if the next-function in the middleware was called.

## shared-types

shared-types provides structs for the http-communication between services.

## utils

In utils we have the map and filter function to map or filter over an array like in javascript.
