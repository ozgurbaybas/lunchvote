# LunchVote

LunchVote is an internal company backend service for restaurant discovery, ratings, team polls, and rule-based lunch recommendations.

## Architecture

- Go
- PostgreSQL
- DDD
- Clean Architecture
- Modular Monolith

## Modules

- identity
- restaurant
- rating
- poll
- recommendation

## Implemented Features

- create users
- create teams
- add team members
- create restaurants
- list restaurants
- create ratings
- list restaurant ratings
- create polls
- vote on polls
- get poll results
- get team recommendations

## Project Structure

```text
cmd/
internal/
platform/
modules/
migrations/
docs/