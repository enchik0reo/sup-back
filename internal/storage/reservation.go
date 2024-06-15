package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/enchik0reo/sup-back/internal/models"
)

func (s *RentStoage) GetReserved(ctx context.Context, from, to string) ([]models.Sup, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT day, fk_sup_id
	FROM reserved
	WHERE day BETWEEN $1 AND $2
	ORDER BY day`)
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

		if err := rows.Scan(&reserv.Day, &reserv.ModelID); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}

		reserves = append(reserves, reserv)
	}

	supsList, err := s.GetPrices(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get sups: %w", err)
	}

	sups := supInfo(reserves, supsList)

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

func supInfo(reserved []models.Reserved, supsList []models.SupInfo) []models.Sup {
	temp := make(map[int64][]time.Time, len(supsList))

	for _, res := range reserved {
		rd := temp[res.ModelID]
		rd = append(rd, res.Day)
		temp[res.ModelID] = rd
	}

	sups := make([]models.Sup, len(supsList))

	for i, sup := range supsList {
		sups[i].ID = sup.ID
		sups[i].Name = sup.Name
		sups[i].Price = sup.Price
		sups[i].ReservedDays = temp[sup.ID]
	}

	return sups
}
