# Parking Lot API

[![codecov](https://codecov.io/gh/teachmind/Parcel-Service/branch/master/graph/badge.svg?token=HivKkjhfjl)](https://codecov.io/gh/teachmind/Parcel-Service)
[![Go Report Card](https://goreportcard.com/badge/github.com/teachmind/Parcel-Service)](https://goreportcard.com/report/github.com/teachmind/Parcel-Service)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/934b654ea9eb4f72b98138b21b5aea94)](https://www.codacy.com/gh/teachmind/Parcel-Service/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=teachmind/Parcel-Service&amp;utm_campaign=Badge_Grade)
[![](https://godoc.org/github.com/teachmind/Parcel-Service?status.svg)](https://godoc.org/github.com/teachmind/Parcel-Service)

## Features 
-   Create Parking Lot
-   Parking Lot Status
-   Toggle Maintenance Mode
-   Park Vehicle
-   Unpark Vehicle
-   Daily Parking Report

## Project Structure
    .
    |-- migration           # Contains migration files
    |-- .env.example        # example/structure of .env file
    |-- Dockerfile          # Used to build docker image.
    |-- go.mode             # Define's the module's import path used for root directory
    |-- go.sum              # Contains the expected cryptographic checksums of the content of specific module versions
    |-- readme.md           # Explains project installation and other informations

## Tools and Technology
-   Golang
-   PostgreSQL

## Installation
-   **Step-1:** Copy/rename `.env.example` file as `.env`. Change the `APP_PORT`, `DB_PORT`, `DB_NAME`,`DB_HOST`, `DB_USER`, `DB_PASSWORD` value as per your DB and Project setup. 
    
    For local environment, you can use the following commands to set env: 
    
    `export DB_NAME=parkingapp`
    `export DB_HOST=localhost`
    `export DB_PORT=5432`
    
-   **Step-2:** Import `parkingapp.sql` in the database
-   **Step-4:** To start server run `go run main.go`
-   **Step-5:** For API collection, Import `Parking Lot API.postman_collection.json` in postman.
