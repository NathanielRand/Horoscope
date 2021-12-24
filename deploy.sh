#!/bin/bash

# Export Go to PATH env var.
export PATH=$PATH:/usr/local/go/bin

# Turn off Go modules.
export GO111MODULE=off

# Pull most recent updates from Github repo.
git pull https://github.com/NathanielRand/Horoscope

# Build go program.
go build

# Kill previously running background process.
kill $(pgrep Horoscope)

# Run and detach updated go program into a new process.
nohup ./Horoscope &