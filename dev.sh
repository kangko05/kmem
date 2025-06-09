#!/bin/bash

POSTGRES_PASSWORD=yourpassword \
    JWT_SECRET_KEY=jwtsecretkey \
    docker compose up
