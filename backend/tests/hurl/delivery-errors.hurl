# Delivery Error & Edge Case Tests
# Tests validation, missing data, and error responses
#
# Usage:
#   hurl --test --variable base_url=http://localhost:8080/api backend/tests/hurl/delivery-errors.hurl

# --- Login to get test user ID ---
POST {{base_url}}/auth/login
Content-Type: application/json

{
  "username": "test",
  "password": "T3stPass135!"
}

HTTP 200
[Captures]
user_id: jsonpath "$.id"

# --- ERROR: Create delivery with missing required fields ---
POST {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

{
  "pickup": {
    "name": "Test",
    "lat": 40.0,
    "lon": -74.0
  },
  "dropoff": {
    "lat": 40.1,
    "lon": -73.9
  },
  "earnings": {
    "tip": 0,
    "base": 0
  }
}

HTTP 400

[Asserts]
jsonpath "$.error" contains "start time is required"

# --- ERROR: Create delivery with end before start ---
POST {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

{
  "start": "2024-06-01T13:00:00Z",
  "end": "2024-06-01T12:00:00Z",
  "pickup": {
    "name": "Test",
    "lat": 40.0,
    "lon": -74.0
  },
  "dropoff": {
    "lat": 40.1,
    "lon": -73.9
  },
  "earnings": {
    "tip": 0,
    "base": 0
  }
}

HTTP 400

[Asserts]
jsonpath "$.error" contains "end time must be after start time"

# --- ERROR: Create delivery with negative earnings ---
POST {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

{
  "start": "2024-06-01T12:00:00Z",
  "end": "2024-06-01T12:30:00Z",
  "pickup": {
    "name": "Test",
    "lat": 40.0,
    "lon": -74.0
  },
  "dropoff": {
    "lat": 40.1,
    "lon": -73.9
  },
  "earnings": {
    "tip": -100,
    "base": 0
  }
}

HTTP 400

[Asserts]
jsonpath "$.error" contains "tip cannot be negative"

# --- ERROR: Get delivery with non-existent ID ---
GET {{base_url}}/users/{{user_id}}/deliveries/nonexistent-id

HTTP 404

# --- ERROR: Update with empty array ---
PATCH {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

[]

HTTP 400

[Asserts]
jsonpath "$.error" contains "at least one delivery update is required"

# --- ERROR: Update with missing ID ---
PATCH {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

[
  {
    "note": "Should fail"
  }
]

HTTP 400

[Asserts]
jsonpath "$.error" contains "delivery id is required"

# --- ERROR: Update with no fields to update ---
PATCH {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

[
  {
    "id": "some-id"
  }
]

HTTP 400

[Asserts]
jsonpath "$.error" contains "at least one field must be provided"

# --- ERROR: Delete with no IDs ---
DELETE {{base_url}}/users/{{user_id}}/deliveries

HTTP 400

[Asserts]
jsonpath "$.error" contains "at least one delivery id is required"

# --- ERROR: Delete non-existent delivery ---
DELETE {{base_url}}/users/{{user_id}}/deliveries?id=nonexistent-id

HTTP 404

# --- ERROR: Invalid date format for list ---
GET {{base_url}}/users/{{user_id}}/deliveries?start_date=not-a-date

HTTP 400

[Asserts]
jsonpath "$.error" contains "invalid start_date format"
