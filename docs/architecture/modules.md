# Module Boundaries

## identity

Responsible for:
- users
- teams
- team memberships

Main capabilities:
- create user
- create team
- add member to team

## restaurant

Responsible for:
- restaurant catalog
- supported meal cards
- active/passive restaurant state

Main capabilities:
- create restaurant
- list restaurants

## rating

Responsible for:
- restaurant ratings
- restaurant comments

Main capabilities:
- create rating
- list ratings by restaurant

Business rule:
- one user can rate one restaurant only once

## poll

Responsible for:
- team lunch polls
- poll options
- votes
- results

Main capabilities:
- create poll
- vote
- get poll results

Business rules:
- only team members can create/vote
- one user can vote once per poll
- closed polls cannot accept votes

## recommendation

Responsible for:
- team-based restaurant recommendations

Current strategy:
- rule-based scoring using:
  - active restaurants
  - rating averages
  - historical poll votes