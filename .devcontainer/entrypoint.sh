#!/bin/zsh

if [ ! -f ".env" ]; then
    cp .env.example .env
fi

rm -rf .devcontainer/.env

go mod tidy

while sleep 1000; do :; done
