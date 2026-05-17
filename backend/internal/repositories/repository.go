// Package repositories provides database access interfaces and implementations.
package repositories

import (
	"context"
	"fmt"
	"time"

	"grubbin-data/backend/internal/core"
	"grubbin-data/backend/internal/models"
)

// Repository defines the interface for data access.
type Repository interface {
	// User operations
	CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error
	CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	// Delivery operations
	CreateDelivery(ctx context.Context, delivery *models.Delivery) error
	GetDeliveryByID(ctx context.Context, userID, deliveryID string) (*models.Delivery, error)
	GetDeliveriesByUser(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.Delivery, error)
	UpdateDelivery(ctx context.Context, userID, deliveryID string, fields core.DeliveryUpdateFields) error
	UpdateDeliveries(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error
	DeleteDelivery(ctx context.Context, userID, deliveryID string) error
	DeleteDeliveries(ctx context.Context, userID string, ids []string) error
}

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	db *DB
}

// NewRepository creates a new PostgresRepository.
func NewRepository(db *DB) Repository {
	return &PostgresRepository{db: db}
}

// CreateVehicle inserts a new vehicle into the database.
func (r *PostgresRepository) CreateVehicle(ctx context.Context, vehicle *models.Vehicle) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO vehicles (id, name, average_mpg) VALUES ($1, $2, $3)",
		vehicle.ID,
		vehicle.Name,
		vehicle.AverageMPG,
	)
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}
	return nil
}

// CreateUserWithProfile creates a user, their profile, and optionally their vehicle atomically within a transaction.
// If any operation fails, the entire transaction is rolled back.
func (r *PostgresRepository) CreateUserWithProfile(ctx context.Context, user *models.User, profile *models.Profile, vehicle *models.Vehicle) error {
	return r.db.WithTx(ctx, func(tx Querier) error {
		_, err := tx.Exec(ctx,
			"INSERT INTO users (id, username, password_hash, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
			user.ID,
			user.Username,
			user.Password,
			user.Created,
			user.Updated,
		)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		_, err = tx.Exec(ctx,
			"INSERT INTO profiles (id, first_name, last_name) VALUES ($1, $2, $3)",
			profile.ID,
			profile.FirstName,
			profile.LastName,
		)
		if err != nil {
			return fmt.Errorf("create profile: %w", err)
		}

		if vehicle != nil {
			_, err = tx.Exec(ctx,
				"INSERT INTO vehicles (id, name, average_mpg) VALUES ($1, $2, $3)",
				vehicle.ID,
				vehicle.Name,
				vehicle.AverageMPG,
			)
			if err != nil {
				return fmt.Errorf("create vehicle: %w", err)
			}
		}

		return nil
	})
}

// GetUserByID retrieves a user by their ID, including their profile and vehicle.
// Password is not loaded; the field remains nil.
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	var firstName, lastName, vehicleName *string
	var vehicleMPG *float64

	err := r.db.QueryRow(ctx,
		`SELECT
			u.id, u.username, u.created_at, u.updated_at,
			p.first_name, p.last_name,
			v.name, v.average_mpg
		FROM users u
		LEFT JOIN profiles p ON u.id = p.id
		LEFT JOIN vehicles v ON p.id = v.id
		WHERE u.id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Created,
		&user.Updated,
		&firstName,
		&lastName,
		&vehicleName,
		&vehicleMPG,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	// Attach profile if present
	if firstName != nil && lastName != nil {
		user.Profile = &models.Profile{
			ID:        user.ID,
			FirstName: *firstName,
			LastName:  *lastName,
		}
		// Attach vehicle if present
		if vehicleName != nil {
			user.Profile.Vehicle = &models.Vehicle{
				ID:         user.ID,
				Name:       *vehicleName,
				AverageMPG: 0,
			}
			if vehicleMPG != nil {
				user.Profile.Vehicle.AverageMPG = *vehicleMPG
			}
		}
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username, including their profile and vehicle.
// Password hash is loaded for authentication purposes.
func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	var passwordHash string
	var firstName, lastName, vehicleName *string
	var vehicleMPG *float64

	err := r.db.QueryRow(ctx,
		`SELECT
			u.id, u.username, u.password_hash, u.created_at, u.updated_at,
			p.first_name, p.last_name,
			v.name, v.average_mpg
		FROM users u
		LEFT JOIN profiles p ON u.id = p.id
		LEFT JOIN vehicles v ON p.id = v.id
		WHERE u.username = $1`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&passwordHash,
		&user.Created,
		&user.Updated,
		&firstName,
		&lastName,
		&vehicleName,
		&vehicleMPG,
	)
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	user.Password = &passwordHash

	// Attach profile if present
	if firstName != nil && lastName != nil {
		user.Profile = &models.Profile{
			ID:        user.ID,
			FirstName: *firstName,
			LastName:  *lastName,
		}
		// Attach vehicle if present
		if vehicleName != nil {
			user.Profile.Vehicle = &models.Vehicle{
				ID:         user.ID,
				Name:       *vehicleName,
				AverageMPG: 0,
			}
			if vehicleMPG != nil {
				user.Profile.Vehicle.AverageMPG = *vehicleMPG
			}
		}
	}

	return user, nil
}

// --- Delivery Operations ---

// CreateDelivery inserts a delivery and its related records (pickup, dropoff, earnings)
// atomically within a transaction.
func (r *PostgresRepository) CreateDelivery(ctx context.Context, delivery *models.Delivery) error {
	return r.db.WithTx(ctx, func(tx Querier) error {
		// Insert pickup location
		_, err := tx.Exec(ctx,
			"INSERT INTO pickup_locations (id, name, lat, lon) VALUES ($1, $2, $3, $4)",
			delivery.PickupLocationID,
			delivery.Pickup.Name,
			delivery.Pickup.Lat,
			delivery.Pickup.Lon,
		)
		if err != nil {
			return fmt.Errorf("create pickup location: %w", err)
		}

		// Insert dropoff location
		_, err = tx.Exec(ctx,
			"INSERT INTO dropoff_locations (id, lat, lon) VALUES ($1, $2, $3)",
			delivery.DropoffLocationID,
			delivery.Dropoff.Lat,
			delivery.Dropoff.Lon,
		)
		if err != nil {
			return fmt.Errorf("create dropoff location: %w", err)
		}

		// Insert earnings
		_, err = tx.Exec(ctx,
			"INSERT INTO earnings (id, tip_cents, base_cents, bonus_cents) VALUES ($1, $2, $3, $4)",
			delivery.EarningsID,
			delivery.Earnings.Tip,
			delivery.Earnings.Base,
			delivery.Earnings.Bonus,
		)
		if err != nil {
			return fmt.Errorf("create earnings: %w", err)
		}

		// Insert delivery
		_, err = tx.Exec(ctx,
			"INSERT INTO deliveries (id, user_id, pickup_location_id, dropoff_location_id, earnings_id, start_time, end_time, note) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			delivery.ID,
			delivery.UserID,
			delivery.PickupLocationID,
			delivery.DropoffLocationID,
			delivery.EarningsID,
			delivery.Start,
			delivery.End,
			delivery.Note,
		)
		if err != nil {
			return fmt.Errorf("create delivery: %w", err)
		}

		return nil
	})
}

// GetDeliveryByID retrieves a delivery by its ID, including related pickup, dropoff, and earnings.
// Verifies the delivery belongs to the specified user.
func (r *PostgresRepository) GetDeliveryByID(ctx context.Context, userID, deliveryID string) (*models.Delivery, error) {
	delivery := &models.Delivery{}

	err := r.db.QueryRow(ctx,
		`SELECT
			d.id, d.user_id, d.pickup_location_id, d.dropoff_location_id, d.earnings_id,
			d.created_at, d.updated_at, d.start_time, d.end_time, d.note,
			p.id, p.name, p.lat, p.lon,
			dl.id, dl.lat, dl.lon,
			e.id, e.tip_cents, e.base_cents, e.bonus_cents
		FROM deliveries d
		JOIN pickup_locations p ON d.pickup_location_id = p.id
		JOIN dropoff_locations dl ON d.dropoff_location_id = dl.id
		JOIN earnings e ON d.earnings_id = e.id
		WHERE d.id = $1 AND d.user_id = $2`,
		deliveryID,
		userID,
	).Scan(
		&delivery.ID,
		&delivery.UserID,
		&delivery.PickupLocationID,
		&delivery.DropoffLocationID,
		&delivery.EarningsID,
		&delivery.Created,
		&delivery.Updated,
		&delivery.Start,
		&delivery.End,
		&delivery.Note,
		&delivery.Pickup.ID,
		&delivery.Pickup.Name,
		&delivery.Pickup.Lat,
		&delivery.Pickup.Lon,
		&delivery.Dropoff.ID,
		&delivery.Dropoff.Lat,
		&delivery.Dropoff.Lon,
		&delivery.Earnings.ID,
		&delivery.Earnings.Tip,
		&delivery.Earnings.Base,
		&delivery.Earnings.Bonus,
	)
	if err != nil {
		return nil, fmt.Errorf("get delivery by id: %w", err)
	}

	return delivery, nil
}

// GetDeliveriesByUser retrieves all deliveries for a user, optionally filtered by date range.
func (r *PostgresRepository) GetDeliveriesByUser(ctx context.Context, userID string, startDate, endDate *time.Time) ([]models.Delivery, error) {
	query := `SELECT
		d.id, d.user_id, d.pickup_location_id, d.dropoff_location_id, d.earnings_id,
		d.created_at, d.updated_at, d.start_time, d.end_time, d.note,
		p.id, p.name, p.lat, p.lon,
		dl.id, dl.lat, dl.lon,
		e.id, e.tip_cents, e.base_cents, e.bonus_cents
	FROM deliveries d
	JOIN pickup_locations p ON d.pickup_location_id = p.id
	JOIN dropoff_locations dl ON d.dropoff_location_id = dl.id
	JOIN earnings e ON d.earnings_id = e.id
	WHERE d.user_id = $1`

	args := []any{userID}
	argIndex := 2

	if startDate != nil {
		query += fmt.Sprintf(" AND d.start_time >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND d.start_time <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query += " ORDER BY d.start_time DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get deliveries by user: %w", err)
	}
	defer rows.Close()

	var deliveries []models.Delivery
	for rows.Next() {
		d := models.Delivery{}
		if scanErr := rows.Scan(
			&d.ID, &d.UserID, &d.PickupLocationID, &d.DropoffLocationID, &d.EarningsID,
			&d.Created, &d.Updated, &d.Start, &d.End, &d.Note,
			&d.Pickup.ID, &d.Pickup.Name, &d.Pickup.Lat, &d.Pickup.Lon,
			&d.Dropoff.ID, &d.Dropoff.Lat, &d.Dropoff.Lon,
			&d.Earnings.ID, &d.Earnings.Tip, &d.Earnings.Base, &d.Earnings.Bonus,
		); scanErr != nil {
			return nil, fmt.Errorf("scan delivery: %w", scanErr)
		}
		deliveries = append(deliveries, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate deliveries: %w", err)
	}

	return deliveries, nil
}

// UpdateDelivery updates a delivery and its related records based on the provided fields.
// Only non-nil fields in the DeliveryUpdateFields maps are updated.
func (r *PostgresRepository) UpdateDelivery(ctx context.Context, userID, deliveryID string, fields core.DeliveryUpdateFields) error {
	return r.db.WithTx(ctx, func(tx Querier) error {
		// Verify ownership and fetch related IDs
		var pickupID, dropoffID, earningsID string
		err := tx.QueryRow(ctx,
			"SELECT pickup_location_id, dropoff_location_id, earnings_id FROM deliveries WHERE id = $1 AND user_id = $2",
			deliveryID, userID,
		).Scan(&pickupID, &dropoffID, &earningsID)
		if err != nil {
			return fmt.Errorf("verify delivery ownership: %w", err)
		}

		// Update pickup location
		if len(fields.Pickup) > 0 {
			if updateErr := updateTable(tx, ctx, "pickup_locations", fields.Pickup, "id", pickupID); updateErr != nil {
				return fmt.Errorf("update pickup location: %w", updateErr)
			}
		}

		// Update dropoff location
		if len(fields.Dropoff) > 0 {
			if updateErr := updateTable(tx, ctx, "dropoff_locations", fields.Dropoff, "id", dropoffID); updateErr != nil {
				return fmt.Errorf("update dropoff location: %w", updateErr)
			}
		}

		// Update earnings
		if len(fields.Earnings) > 0 {
			if updateErr := updateTable(tx, ctx, "earnings", fields.Earnings, "id", earningsID); updateErr != nil {
				return fmt.Errorf("update earnings: %w", updateErr)
			}
		}

		// Update delivery
		if len(fields.Delivery) > 0 {
			fields.Delivery["updated_at"] = time.Now().UTC()
			if updateErr := updateTable(tx, ctx, "deliveries", fields.Delivery, "id", deliveryID); updateErr != nil {
				return fmt.Errorf("update delivery: %w", updateErr)
			}
		}

		return nil
	})
}

// UpdateDeliveries updates multiple deliveries in a single transaction.
func (r *PostgresRepository) UpdateDeliveries(ctx context.Context, userID string, updates []models.UpdateDeliveryRequest) error {
	return r.db.WithTx(ctx, func(tx Querier) error {
		for _, req := range updates {
			fields := core.BuildDeliveryUpdateFields(req)

			// Verify ownership and fetch related IDs
			var pickupID, dropoffID, earningsID string
			err := tx.QueryRow(ctx,
				"SELECT pickup_location_id, dropoff_location_id, earnings_id FROM deliveries WHERE id = $1 AND user_id = $2",
				req.ID, userID,
			).Scan(&pickupID, &dropoffID, &earningsID)
			if err != nil {
				return fmt.Errorf("verify delivery ownership for %s: %w", req.ID, err)
			}

			// Update related tables
			if len(fields.Pickup) > 0 {
				if updateErr := updateTable(tx, ctx, "pickup_locations", fields.Pickup, "id", pickupID); updateErr != nil {
					return fmt.Errorf("update pickup for %s: %w", req.ID, updateErr)
				}
			}
			if len(fields.Dropoff) > 0 {
				if updateErr := updateTable(tx, ctx, "dropoff_locations", fields.Dropoff, "id", dropoffID); updateErr != nil {
					return fmt.Errorf("update dropoff for %s: %w", req.ID, updateErr)
				}
			}
			if len(fields.Earnings) > 0 {
				if updateErr := updateTable(tx, ctx, "earnings", fields.Earnings, "id", earningsID); updateErr != nil {
					return fmt.Errorf("update earnings for %s: %w", req.ID, updateErr)
				}
			}
			if len(fields.Delivery) > 0 {
				fields.Delivery["updated_at"] = time.Now().UTC()
				if updateErr := updateTable(tx, ctx, "deliveries", fields.Delivery, "id", req.ID); updateErr != nil {
					return fmt.Errorf("update delivery %s: %w", req.ID, updateErr)
				}
			}
		}
		return nil
	})
}

// DeleteDelivery deletes a delivery and its related records (earnings, pickup, dropoff).
func (r *PostgresRepository) DeleteDelivery(ctx context.Context, userID, deliveryID string) error {
	return r.db.WithTx(ctx, func(tx Querier) error {
		return deleteDeliveryRecords(tx, ctx, userID, deliveryID)
	})
}

// DeleteDeliveries deletes multiple deliveries and their related records.
func (r *PostgresRepository) DeleteDeliveries(ctx context.Context, userID string, ids []string) error {
	return r.db.WithTx(ctx, func(tx Querier) error {
		for _, id := range ids {
			if err := deleteDeliveryRecords(tx, ctx, userID, id); err != nil {
				return err
			}
		}
		return nil
	})
}

// deleteDeliveryRecords deletes a delivery and its related records within a transaction.
// Order matters: delete the delivery first (removes FK references), then the orphaned records.
func deleteDeliveryRecords(tx Querier, ctx context.Context, userID, deliveryID string) error {
	// Fetch related IDs
	var pickupID, dropoffID, earningsID string
	err := tx.QueryRow(ctx,
		"SELECT pickup_location_id, dropoff_location_id, earnings_id FROM deliveries WHERE id = $1 AND user_id = $2",
		deliveryID, userID,
	).Scan(&pickupID, &dropoffID, &earningsID)
	if err != nil {
		return fmt.Errorf("verify delivery ownership: %w", err)
	}

	// Delete delivery first (removes FK references to related tables)
	if _, err := tx.Exec(ctx, "DELETE FROM deliveries WHERE id = $1", deliveryID); err != nil {
		return fmt.Errorf("delete delivery: %w", err)
	}

	// Delete orphaned related records
	if _, err := tx.Exec(ctx, "DELETE FROM earnings WHERE id = $1", earningsID); err != nil {
		return fmt.Errorf("delete earnings: %w", err)
	}
	if _, err := tx.Exec(ctx, "DELETE FROM pickup_locations WHERE id = $1", pickupID); err != nil {
		return fmt.Errorf("delete pickup location: %w", err)
	}
	if _, err := tx.Exec(ctx, "DELETE FROM dropoff_locations WHERE id = $1", dropoffID); err != nil {
		return fmt.Errorf("delete dropoff location: %w", err)
	}

	return nil
}

// updateTable builds and executes a dynamic UPDATE query for the given table and fields.
func updateTable(tx Querier, ctx context.Context, table string, fields map[string]any, idColumn, idValue string) error {
	if len(fields) == 0 {
		return nil
	}

	setClauses := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields)+1)
	argIndex := 1

	for column, value := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", table, joinStrings(setClauses, ", "), idColumn, argIndex)
	args = append(args, idValue)

	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update %s: %w", table, err)
	}

	return nil
}

// joinStrings joins a slice of strings with the given separator.
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
