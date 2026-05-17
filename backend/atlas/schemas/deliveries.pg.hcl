// Delivery domain: pickups, dropoffs, earnings, deliveries

// --- Pickup Locations ---
table "pickup_locations" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "name" {
    type = text
    null = false
  }

  column "lat" {
    type = double_precision
    null = false
  }

  column "lon" {
    type = double_precision
    null = false
  }

  primary_key {
    columns = [column.id]
  }
}

// --- Dropoff Locations ---
table "dropoff_locations" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "lat" {
    type = double_precision
    null = false
  }

  column "lon" {
    type = double_precision
    null = false
  }

  primary_key {
    columns = [column.id]
  }
}

// --- Earnings ---
table "earnings" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "tip_cents" {
    type = integer
    null = false
  }

  column "base_cents" {
    type = integer
    null = false
  }

  column "bonus_cents" {
    type = integer
    null = true
  }

  primary_key {
    columns = [column.id]
  }
}

// --- Deliveries ---
table "deliveries" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "user_id" {
    type = text
    null = false
  }

  column "pickup_location_id" {
    type = text
    null = false
  }

  column "dropoff_location_id" {
    type = text
    null = false
  }

  column "earnings_id" {
    type = text
    null = false
  }

  column "created_at" {
    type    = timestamptz
    null    = false
    default = sql("now()")
  }

  column "updated_at" {
    type = timestamptz
    null = true
  }

  column "start_time" {
    type = timestamptz
    null = false
  }

  column "end_time" {
    type = timestamptz
    null = false
  }

  column "note" {
    type = text
    null = true
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "deliveries_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }

  foreign_key "deliveries_pickup_location_id_fkey" {
    columns     = [column.pickup_location_id]
    ref_columns = [table.pickup_locations.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }

  foreign_key "deliveries_dropoff_location_id_fkey" {
    columns     = [column.dropoff_location_id]
    ref_columns = [table.dropoff_locations.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }

  foreign_key "deliveries_earnings_id_fkey" {
    columns     = [column.earnings_id]
    ref_columns = [table.earnings.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }

  index "deliveries_user_id_idx" {
    columns = [column.user_id]
  }

  index "deliveries_pickup_location_id_idx" {
    columns = [column.pickup_location_id]
  }

  index "deliveries_dropoff_location_id_idx" {
    columns = [column.dropoff_location_id]
  }

  index "deliveries_earnings_id_idx" {
    columns = [column.earnings_id]
  }
}
