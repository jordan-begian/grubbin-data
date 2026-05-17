// Atlas project configuration
// Usage: atlas schema apply --env local
//        atlas migrate diff initial --env local

// --- Database Connection Variables ---
// These read from environment variables with sensible defaults for local development.

variable "db_host" {
  type    = string
  default = getenv("DB_HOST")
}

variable "db_port" {
  type    = string
  default = getenv("DB_PORT")
}

variable "db_name" {
  type    = string
  default = getenv("DB_NAME")
}

variable "db_user" {
  type    = string
  default = getenv("DB_USER")
}

variable "db_password" {
  type    = string
  default = getenv("DB_PASSWORD")
}

variable "db_ssl_mode" {
  type    = string
  default = getenv("DB_SSL_MODE")
}

// --- Computed Database URL ---
// Constructs the PostgreSQL connection string from individual env vars.
// Falls back to a single DATABASE_URL if the component vars are not set.

variable "database_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

locals {
  // Prefer DATABASE_URL if set; otherwise construct from components
  db_url = var.database_url != "" ? var.database_url : "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=${var.db_ssl_mode}"
}

// --- Environments ---

env "local" {
  src = "file://schemas/"
  url = local.db_url
  dev = "docker://postgres/17/dev?search_path=public"
}

env "docker" {
  src = "file://schemas/"
  url = local.db_url
  // Uses a real database (atlas_dev) instead of Docker-in-Docker for containerized environments
  dev = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/atlas_dev?sslmode=${var.db_ssl_mode}&search_path=public"
}

env "prod" {
  src = "file://schemas/"
  url = local.db_url
  dev = "docker://postgres/17/dev?search_path=public"
}
