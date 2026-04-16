---

## `docs/decisions/adr-001-modular-monolith.md`

```md
# ADR-001: Modular Monolith

## Status
Accepted

## Context
The product needs clear domain boundaries and production-grade maintainability, but does not yet require microservice operational complexity.

## Decision
The system will be implemented as a modular monolith.

## Consequences
- simpler deployment
- easier local development
- strong domain boundaries remain possible
- future extraction is still possible if necessary