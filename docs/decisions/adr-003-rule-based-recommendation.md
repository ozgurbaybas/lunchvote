# ADR-003: Rule-Based Recommendation First

## Status
Accepted

## Context
The product needs an initial recommendation capability before enough behavioral data exists for advanced personalization or ML-based approaches.

## Decision
The first recommendation engine will be rule-based.

## Current Inputs
- active restaurants
- historical poll votes
- restaurant rating averages

## Consequences
- fast to implement
- explainable behavior
- easy to iterate
- limited personalization in the first version