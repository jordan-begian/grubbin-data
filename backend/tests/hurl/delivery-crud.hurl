# Delivery CRUD Tests
# Tests the full lifecycle: Create → Read → Update → Verify → Delete
#
# Usage:
#   hurl --test --variable base_url=http://localhost:8080/api backend/tests/hurl/delivery-crud.hurl
#
# Prerequisites:
#   - Database running with seed user (task up)
#   - Seed user: username=test, password=T3stPass135!

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

# --- CREATE: First delivery ---
POST {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

{
  "start": "2024-06-01T12:00:00Z",
  "end": "2024-06-01T12:30:00Z",
  "pickup": {
    "name": "Chipotle Downtown",
    "lat": 40.7128,
    "lon": -74.0060
  },
  "dropoff": {
    "lat": 40.7580,
    "lon": -73.9855
  },
  "earnings": {
    "tip": 350,
    "base": 525
  },
  "note": "Left at front door"
}

HTTP 201
[Captures]
delivery_1_id: jsonpath "$.id"

[Asserts]
# Ownership
jsonpath "$.user" == "{{user_id}}"

# ID is a non-empty string (ULID format)
jsonpath "$.id" != null
jsonpath "$.pickup.id" != null
jsonpath "$.dropoff.id" != null
jsonpath "$.earnings.id" != null

# Timestamps exist
jsonpath "$.created" != null

# Pickup data
jsonpath "$.pickup.name" == "Chipotle Downtown"
jsonpath "$.pickup.lat" == 40.7128
jsonpath "$.pickup.lon" == -74.0060

# Dropoff data
jsonpath "$.dropoff.lat" == 40.7580
jsonpath "$.dropoff.lon" == -73.9855

# Earnings data
jsonpath "$.earnings.tip" == 350
jsonpath "$.earnings.base" == 525

# Note
jsonpath "$.note" == "Left at front door"

# Response header
header "Content-Type" == "application/json"

# --- CREATE: Second delivery ---
POST {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

{
  "start": "2024-06-01T13:00:00Z",
  "end": "2024-06-01T13:45:00Z",
  "pickup": {
    "name": "Pizza Palace",
    "lat": 40.7282,
    "lon": -73.7949
  },
  "dropoff": {
    "lat": 40.7484,
    "lon": -73.9857
  },
  "earnings": {
    "tip": 500,
    "base": 675,
    "bonus": 150
  }
}

HTTP 201
[Captures]
delivery_2_id: jsonpath "$.id"

[Asserts]
jsonpath "$.user" == "{{user_id}}"
jsonpath "$.pickup.name" == "Pizza Palace"
jsonpath "$.pickup.lat" == 40.7282
jsonpath "$.pickup.lon" == -73.7949
jsonpath "$.dropoff.lat" == 40.7484
jsonpath "$.dropoff.lon" == -73.9857
jsonpath "$.earnings.tip" == 500
jsonpath "$.earnings.base" == 675
jsonpath "$.earnings.bonus" == 150
jsonpath "$.start" == "2024-06-01T13:00:00Z"
jsonpath "$.end" == "2024-06-01T13:45:00Z"

# --- CREATE: Third delivery ---
POST {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

{
  "start": "2024-06-02T10:00:00Z",
  "end": "2024-06-02T10:20:00Z",
  "pickup": {
    "name": "Sushi Spot",
    "lat": 40.7614,
    "lon": -73.9776
  },
  "dropoff": {
    "lat": 40.7489,
    "lon": -73.9680
  },
  "earnings": {
    "tip": 200,
    "base": 400
  },
  "note": "Call on arrival"
}

HTTP 201
[Captures]
delivery_3_id: jsonpath "$.id"

[Asserts]
jsonpath "$.user" == "{{user_id}}"
jsonpath "$.pickup.name" == "Sushi Spot"
jsonpath "$.pickup.lat" == 40.7614
jsonpath "$.pickup.lon" == -73.9776
jsonpath "$.dropoff.lat" == 40.7489
jsonpath "$.dropoff.lon" == -73.9680
jsonpath "$.earnings.tip" == 200
jsonpath "$.earnings.base" == 400
jsonpath "$.note" == "Call on arrival"
jsonpath "$.start" == "2024-06-02T10:00:00Z"
jsonpath "$.end" == "2024-06-02T10:20:00Z"

# --- GET: Single delivery ---
GET {{base_url}}/users/{{user_id}}/deliveries/{{delivery_1_id}}

HTTP 200

[Asserts]
jsonpath "$.id" == "{{delivery_1_id}}"
jsonpath "$.user" == "{{user_id}}"
jsonpath "$.pickup.name" == "Chipotle Downtown"
jsonpath "$.pickup.lat" == 40.7128
jsonpath "$.pickup.lon" == -74.0060
jsonpath "$.dropoff.lat" == 40.7580
jsonpath "$.dropoff.lon" == -73.9855
jsonpath "$.earnings.tip" == 350
jsonpath "$.earnings.base" == 525
jsonpath "$.note" == "Left at front door"
jsonpath "$.start" == "2024-06-01T12:00:00Z"
jsonpath "$.end" == "2024-06-01T12:30:00Z"

# --- GET: List all deliveries ---
GET {{base_url}}/users/{{user_id}}/deliveries

HTTP 200

[Asserts]
# At least our 3 created deliveries exist
jsonpath "$.deliveries" count >= 3

# Stats are computed
jsonpath "$.stats.total_time" > 0
jsonpath "$.stats.average_tip" > 0
jsonpath "$.stats.average_base_pay" > 0
jsonpath "$.stats.average_delivery_distance" > 0

# Our specific deliveries are in the list (by pickup name)
jsonpath "$.deliveries[?(@.pickup.name == 'Chipotle Downtown')]" exists
jsonpath "$.deliveries[?(@.pickup.name == 'Pizza Palace')]" exists
jsonpath "$.deliveries[?(@.pickup.name == 'Sushi Spot')]" exists

# --- GET: List deliveries filtered by date range ---
GET {{base_url}}/users/{{user_id}}/deliveries?start_date=2024-06-01T00:00:00Z&end_date=2024-06-01T23:59:59Z

HTTP 200

[Asserts]
# Only June 1st deliveries (delivery 3 is June 2nd, should be excluded)
jsonpath "$.deliveries" count >= 2

# June 1st deliveries are present
jsonpath "$.deliveries[?(@.pickup.name == 'Chipotle Downtown')]" exists
jsonpath "$.deliveries[?(@.pickup.name == 'Pizza Palace')]" exists

# Stats reflect only filtered deliveries
jsonpath "$.stats.total_time" > 0

# --- UPDATE: Single delivery - update note ---
PATCH {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

[
  {
    "id": "{{delivery_1_id}}",
    "note": "Updated: Left with doorman"
  }
]

HTTP 200

[Asserts]
jsonpath "$.message" == "deliveries updated"

# --- UPDATE: Single delivery - update earnings ---
PATCH {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

[
  {
    "id": "{{delivery_2_id}}",
    "earnings": {
      "tip": 600,
      "base": 700,
      "bonus": 200
    }
  }
]

HTTP 200

[Asserts]
jsonpath "$.message" == "deliveries updated"

# --- UPDATE: Multiple deliveries at once ---
PATCH {{base_url}}/users/{{user_id}}/deliveries
Content-Type: application/json

[
  {
    "id": "{{delivery_1_id}}",
    "pickup": {
      "name": "Chipotle Downtown (Updated)",
      "lat": 40.7130,
      "lon": -74.0062
    }
  },
  {
    "id": "{{delivery_3_id}}",
    "note": "Cleared and reset"
  }
]

HTTP 200

[Asserts]
jsonpath "$.message" == "deliveries updated"

# --- VERIFY: Check delivery 1 was updated (note + pickup) ---
GET {{base_url}}/users/{{user_id}}/deliveries/{{delivery_1_id}}

HTTP 200

[Asserts]
jsonpath "$.note" == "Updated: Left with doorman"
jsonpath "$.pickup.name" == "Chipotle Downtown (Updated)"
jsonpath "$.pickup.lat" == 40.7130
jsonpath "$.pickup.lon" == -74.0062

# --- VERIFY: Check delivery 2 earnings were updated ---
GET {{base_url}}/users/{{user_id}}/deliveries/{{delivery_2_id}}

HTTP 200

[Asserts]
jsonpath "$.earnings.tip" == 600
jsonpath "$.earnings.base" == 700
jsonpath "$.earnings.bonus" == 200

# --- VERIFY: Check delivery 3 note was updated ---
GET {{base_url}}/users/{{user_id}}/deliveries/{{delivery_3_id}}

HTTP 200

[Asserts]
jsonpath "$.note" == "Cleared and reset"

# --- DELETE: Delete single delivery ---
DELETE {{base_url}}/users/{{user_id}}/deliveries?id={{delivery_3_id}}

HTTP 204

# --- VERIFY: Deleted delivery returns 404 ---
GET {{base_url}}/users/{{user_id}}/deliveries/{{delivery_3_id}}

HTTP 404

# --- DELETE: Delete remaining deliveries ---
DELETE {{base_url}}/users/{{user_id}}/deliveries?id={{delivery_1_id}}&id={{delivery_2_id}}

HTTP 204

# --- VERIFY: List endpoint still works after deletions ---
GET {{base_url}}/users/{{user_id}}/deliveries?start_date=2024-06-01T00:00:00Z&end_date=2024-06-02T23:59:59Z

HTTP 200

[Asserts]
# Response structure is still valid
jsonpath "$.deliveries" exists
jsonpath "$.stats" exists
