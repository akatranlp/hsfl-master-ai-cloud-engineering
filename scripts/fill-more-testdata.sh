#!/bin/env bash

# Auth needs to be disabled for this to work

for i in {1..100}; do
    curl -X POST -H "Content-Type: application/json" -d '{"name":"The Name of the book'$i'","description": "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."}' http://vv.hsfl.de:32131/api/v1/books
done
