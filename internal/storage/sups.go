package storage

import (
	"context"
	"fmt"

	"github.com/enchik0reo/sup-back/internal/models"
)

func (s *RentStoage) GetPrices(ctx context.Context) ([]models.SupInfo, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT sup_id, model_name, price FROM sups ORDER BY sup_id`)
	if err != nil {
		return nil, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get price list: %w", err)
	}
	defer rows.Close()

	sups := []models.SupInfo{}

	for rows.Next() {
		sup := models.SupInfo{}

		if err := rows.Scan(&sup.ID, &sup.Name, &sup.Price); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}

		sups = append(sups, sup)
	}

	return sups, nil
}

func (s *RentStoage) EditPrice(ctx context.Context, id, newPrice int64) (int64, error) {
	stmt, err := s.db.PrepareContext(ctx, `UPDATE sups SET price = $1
	WHERE sup_id = $2 RETURNING sup_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, newPrice, id)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't edit price: %w", err)
	}

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get updated id: %w", err)
	}

	return id, nil
}

func (s *RentStoage) NewSup(ctx context.Context, name string, price int64) (int64, error) {
	stmt, err := s.db.PrepareContext(ctx, `INSERT INTO sups (model_name, price)
	VALUES ($1, $2) RETURNING sup_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, name, price)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't add sup: %w", err)
	}

	var id int64

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get inserted id: %w", err)
	}

	return id, nil
}

func (s *RentStoage) DeleteSup(ctx context.Context, supID int64) (int64, error) {
	stmt, err := s.db.PrepareContext(ctx, `DELETE FROM sups WHERE sup_id = $1 RETURNING sup_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, supID)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't delete sup: %w", err)
	}

	var id int64

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get last deleted id: %w", err)
	}

	return id, nil
}
