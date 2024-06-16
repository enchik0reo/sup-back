package storage

import (
	"context"
	"fmt"
	"sort"

	"github.com/enchik0reo/sup-back/internal/models"
)

func (s *RentStoage) GetReserved(ctx context.Context, from, to string) ([]models.Sup, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT r.day, s.sup_id, s.model_name, s.price
	FROM sups s
	LEFT JOIN (SELECT fk_sup_id, day FROM reserved WHERE day BETWEEN $1 AND $2) r ON s.sup_id = r.fk_sup_id
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

func (s *RentStoage) CreateReservedList(ctx context.Context, reserveList []models.ApproveReserv) error {
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

	temp := make(map[int64]models.Sup)

	for _, res := range reserved {
		s := temp[res.ModelID]

		s.ID = res.ModelID
		s.Name = res.ModelName
		s.Price = res.ModelPrice

		if res.Day != nil {
			s.ReservedDays = append(s.ReservedDays, *res.Day)
		}

		temp[res.ModelID] = s
	}

	sups := make([]models.Sup, 0, len(temp))

	for _, sup := range temp {
		sups = append(sups, sup)
	}

	sort.SliceStable(sups, func(i, j int) bool { return sups[i].ID < sups[j].ID })

	return sups
}
