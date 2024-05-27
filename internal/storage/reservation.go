package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/enchik0reo/sup-back/internal/models"
)

type Tokener interface {
	Create(phone string) (string, error)
	Parse(tokenString string) (string, error)
}

type RentStoage struct {
	token Tokener

	db *sql.DB
}

func NewRentStorage(db *sql.DB, t Tokener) *RentStoage {
	return &RentStoage{db: db, token: t}
}

func (s *RentStoage) GetReserved(ctx context.Context, from, to string) ([]models.Sup, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT r.day, r.fk_sup_id, s.model_name, s.price
	FROM reserved r
	INNER JOIN sups s ON r.fk_sup_id = s.sup_id
	WHERE r.day BETWEEN $1 AND $2
	ORDER BY r.day`)
	if err != nil {
		return nil, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("can't get reserv list: %w", err)
	}
	defer rows.Close()

	reserves := []models.Reserved{}

	for rows.Next() {
		reserv := models.Reserved{}

		if err := rows.Scan(&reserv.Day, &reserv.ModelID, &reserv.ModelName, &reserv.ModelPrice); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}

		reserves = append(reserves, reserv)
	}

	sups := supInfo(reserves)

	return sups, nil
}

func (s *RentStoage) CreateReserved(ctx context.Context, reserve models.Reserved) (int64, error) {
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO reserved (day, fk_sup_id, fk_approve_id)
	VALUES ($1, $2, $3) RETURNING reserv_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, reserve.Day, reserve.ModelID, reserve.ApproveID)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't create reserve: %w", err)
	}

	var id int64

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get last insert id: %w", err)
	}

	return id, nil
}

func (s *RentStoage) CreateReservedList(ctx context.Context, reserveList []models.Reserved) error {
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO reserved (day, fk_sup_id, fk_approve_id)
	VALUES ($1, $2, $3) RETURNING reserv_id`)
	if err != nil {
		return fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, reserve := range reserveList {
		row := stmt.QueryRowContext(ctx, reserve.Day, reserve.ModelID, reserve.ApproveID)

		if err := row.Err(); err != nil {
			return fmt.Errorf("can't create reserve: %w", err)
		}

		var id int64

		if err := row.Scan(&id); err != nil {
			return fmt.Errorf("can't get last insert id: %w", err)
		}
	}

	return nil
}

func (s *RentStoage) DeleteReserved(ctx context.Context, approveID int64) (int64, error) {
	stmt, err := s.db.PrepareContext(ctx, `DELETE FROM reserved WHERE fk_approve_id = $1 RETURNING reserv_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, approveID)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't delete reserve: %w", err)
	}

	var id int64

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get last deleted id: %w", err)
	}

	return id, nil
}

func supInfo(reserved []models.Reserved) []models.Sup {
	if len(reserved) == 0 {
		return nil
	}

	sups := make([]models.Sup, 3)

	for _, res := range reserved {
		sups[res.ModelID-1].ID = res.ModelID
		sups[res.ModelID-1].Name = res.ModelName
		sups[res.ModelID-1].Price = res.ModelPrice
		sups[res.ModelID-1].ReservedDays = append(sups[res.ModelID].ReservedDays, res.Day)
	}

	return sups
}
