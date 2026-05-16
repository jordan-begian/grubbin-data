// User domain: users, profiles, vehicles

// --- Users ---
table "users" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "username" {
    type = text
    null = false
  }

  column "password_hash" {
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

  primary_key {
    columns = [column.id]
  }

  unique "users_username_key" {
    columns = [column.username]
  }
}

// --- Profiles (1:1 with users) ---
table "profiles" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "first_name" {
    type = text
    null = false
  }

  column "last_name" {
    type = text
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "profiles_user_id_fkey" {
    columns     = [column.id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}

// --- Vehicles (1:1 with profiles, optional) ---
table "vehicles" {
  schema = schema.public

  column "id" {
    type = text
    null = false
  }

  column "name" {
    type = text
    null = false
  }

  column "average_mpg" {
    type = double_precision
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "vehicles_profile_id_fkey" {
    columns     = [column.id]
    ref_columns = [table.profiles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
