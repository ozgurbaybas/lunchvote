
---

## `docs/architecture/overview.md`

```md
# Architecture Overview

LunchVote is implemented as a modular monolith using Go.

## Architectural Style

- Domain Driven Design
- Clean Architecture
- Modular Monolith

## Core Principles

- business rules live in domain and application layers
- infrastructure details are isolated
- modules are separated by business capability
- PostgreSQL is the only persistence layer
- REST API is the external interface

## Modules

- identity
- restaurant
- rating
- poll
- recommendation

## Runtime Flow

HTTP request -> interfaces/http -> application -> domain -> infrastructure/postgres