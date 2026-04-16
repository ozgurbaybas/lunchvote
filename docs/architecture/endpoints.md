# API Endpoints

This document describes all available HTTP endpoints for the LunchVote backend.

Base URL:
http://localhost:8080

---

# Health

## GET /health

Check if the service is running.

### Response

```json
{
  "status": "ok"
}


---

# 👤 3. Identity

```md
---

# Identity

## POST /v1/users

Create a new user.

### Request

```json
{
  "id": "11111111-1111-1111-1111-111111111111",
  "name": "User One",
  "email": "userone@example.com"
}

### Response
{
  "id": "11111111-1111-1111-1111-111111111111",
  "name": "User One",
  "email": "userone@example.com",
  "created_at": "2026-04-10T10:00:00Z"
}

POST /v1/teams

Create a team.

Request

{
  "id": "33333333-3333-3333-3333-333333333333",
  "name": "Backend Team",
  "owner_id": "11111111-1111-1111-1111-111111111111"
}


POST /v1/teams/{id}/members

Add a user to a team.

Request

{
  "user_id": "22222222-2222-2222-2222-222222222222"
}


---

# 🍔 4. Restaurant

```md
---

# Restaurant

## POST /v1/restaurants

Create a restaurant.

### Request

```json
{
  "id": "restaurant-1",
  "name": "Kebap House",
  "address": "Ataturk Caddesi No 10",
  "city": "Istanbul",
  "district": "Bakirkoy",
  "supported_meal_cards": ["ticket", "multinet"]
}

GET /v1/restaurants

List all restaurants.

Response

[
  {
    "id": "restaurant-1",
    "name": "Kebap House",
    "city": "Istanbul",
    "district": "Bakirkoy",
    "is_active": true
  }
]


---

# ⭐ 5. Rating

```md
---

# Rating

## POST /v1/ratings

Create a rating.

### Request

```json
{
  "id": "rating-1",
  "restaurant_id": "restaurant-1",
  "user_id": "11111111-1111-1111-1111-111111111111",
  "score": 5,
  "comment": "great food"
}


---

# ⭐ 5. Rating

```md
---

# Rating

## POST /v1/ratings

Create a rating.

### Request

```json
{
  "id": "rating-1",
  "restaurant_id": "restaurant-1",
  "user_id": "11111111-1111-1111-1111-111111111111",
  "score": 5,
  "comment": "great food"
}

GET /v1/restaurants/{id}/ratings

List ratings of a restaurant.

Response

[
  {
    "id": "rating-1",
    "restaurant_id": "restaurant-1",
    "user_id": "11111111-1111-1111-1111-111111111111",
    "score": 5,
    "comment": "great food"
  }
]


---

# 🗳️ 6. Poll

```md
---

# Poll

## POST /v1/polls

Create a poll.

### Request

```json
{
  "id": "poll-1",
  "team_id": "33333333-3333-3333-3333-333333333333",
  "title": "Friday Lunch",
  "restaurant_ids": ["restaurant-1", "restaurant-2"],
  "creator_user_id": "11111111-1111-1111-1111-111111111111"
}


POST /v1/polls/{id}/votes

Vote in a poll.

Request

{
  "user_id": "22222222-2222-2222-2222-222222222222",
  "restaurant_id": "restaurant-1"
}

GET /v1/polls/{id}/results

Get poll results.

Response

{
  "poll_id": "poll-1",
  "results": {
    "restaurant-1": 2,
    "restaurant-2": 1
  }
}


---

# 🤖 7. Recommendation

```md
---

# Recommendation

## GET /v1/teams/{id}/recommendations

Get restaurant recommendations for a team.

### Query Params

- limit (optional)

### Example

GET /v1/teams/{id}/recommendations?limit=3

### Response

```json
[
  {
    "restaurant_id": "restaurant-1",
    "score": 45.0,
    "reasons": [
      "team poll history",
      "strong ratings"
    ]
  }
]



---

# ⚠️ 8. Errors

```md
---

# Error Responses

## 400

```json
{
  "error": "bad request"
}


403
{
  "error": "forbidden"
}

404
{
  "error": "not found"
}

409
{
  "error": "conflict"
}

500

{
  "error": "internal server error"
}


---

# 🔚 9. Notes

```md
---

# Notes

- Users and teams use UUID
- Restaurant IDs are string
- One user can rate a restaurant once
- One user can vote in a poll once
- Only team members can vote
- Recommendation is rule-based (MVP)