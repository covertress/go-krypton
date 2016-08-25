#!/bin/bash -x

# creates the necessary docker images to run testrunner.sh locally

docker build --tag="krypton/cppjit-testrunner" docker-cppjit
docker build --tag="krypton/python-testrunner" docker-python
docker build --tag="krypton/go-testrunner" docker-go
