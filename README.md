# Gin Custom Template

## Overview

This repository servers as a template for applications build with Gin web framework

This template consists of all basic functionality that every project has.
1. Makefile for improving developer experience;
2. Dockerized application and managed by Docker Compose (including PostgreSQl and Redis);
3. Easily customizable config file;
4. Ready-to-use GOOSE tool for migration management;
5. Monitoring with Grafana, Prometheus and Grafana Loki;
6. Middlewares:
   1. Rate limit middleware to prevent being overwhelmed by excessive requests;
   2. Logger middleware with Prometheus;
   3. Error middleware to improve user experience;
   4. Authentication middleware via JWT;
   5. Caching middleware with Redis;
7. Ready-to-use Swagger;
8. Ready-to-use PostgreSQL connection with defined SQL queries;
9. Ready-to-use CRUD operations for User entity;
10. **TBD**: Test environment.

The application uses **Go 1.24** and all newest versions as in **2025-07-31**.

There are two database declared in the `docker-compose.yml` file:
1. **PostgreSQL** as a primary database;
2. **Redis** for caching HTTP requests.

## Prerequisites

- Docker & Docker Compose
- Make

## How to run

Create `.env` file in the project root directory and fill it up according to `.env.sample` file.

To spin up the whole application, you can enter:

```shell
make start
```

To shut it down, you can enter:

```shell
make down
```

## Explore other available Make commands

To display all available commands and their description, you can enter:

```shell
make help
```