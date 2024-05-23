package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/enchik0reo/sup-back/internal/models"
)

type ReservationStoage struct {
	db *sql.DB
}

func NewReservationStorage(db *sql.DB) *ReservationStoage {
	return &ReservationStoage{db: db}
}

func (s *ReservationStoage) GetReserved(ctx context.Context, from, to time.Time) ([]models.Sup, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT r.day, r.fk_sup_id, s.model_name, s.price
	FROM reserved r
	INNER JOIN supboards s ON r.fk_sup_id = s.sup_id
	WHERE r.day BETWEEN $1 AND $2
	ORDER BY r.day`)
	if err != nil {
		return nil, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, correctTime(from), correctTime(to))
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

func (s *ReservationStoage) GetApproveList(ctx context.Context) ([]models.Approve, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT client_phone, client_name, price, order_info
	FROM approve
	WHERE status = 1`)
	if err != nil {
		return nil, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get approve list: %w", err)
	}
	defer rows.Close()

	approves := []models.Approve{}

	for rows.Next() {
		approve := models.Approve{}
		info := ""

		if err := rows.Scan(&approve.ClientNumber, &approve.ClientName, &approve.FullPrice, &info); err != nil {
			return nil, fmt.Errorf("can't scan row: %w", err)
		}

		supInfo, err := infoFromJSON(info)
		if err != nil {
			return nil, fmt.Errorf("can't make sup info: %w", err)
		}

		approve.SupsInfo = supInfo

		approves = append(approves, approve)
	}

	return approves, nil
}

func (c *ReservationStoage) CreateApprove(ctx context.Context, approve models.Approve) (int64, error) {
	stmt, err := c.db.PrepareContext(ctx, `INSERT INTO approve (client_phone, client_name, price, order_info)
	VALUES ($1, $2, $3, $4) RETURNING approve_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	info, err := infoToJSON(approve.SupsInfo)
	if err != nil {
		return 0, fmt.Errorf("can't make info: %w", err)
	}

	row := stmt.QueryRowContext(ctx, approve.ClientNumber, approve.ClientName, approve.FullPrice, info)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't create approve: %w", err)
	}

	var id int64

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get last insert id: %w", err)
	}

	return id, nil
}

func (c *ReservationStoage) ConfirmApprove(ctx context.Context, id int64, phone string) (int64, error) {
	stmt, err := c.db.PrepareContext(ctx, `UPDATE approve SET status = 2
	WHERE approve_id = $1 AND client_phone = $2 RETURNING approve_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id, phone)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't confirm approve: %w", err)
	}

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get confirmed id: %w", err)
	}

	return id, nil
}

func (c *ReservationStoage) CancelApprove(ctx context.Context, id int64, phone string) (int64, error) {
	stmt, err := c.db.PrepareContext(ctx, `UPDATE approve SET status = 0
	WHERE approve_id = $1 AND client_phone = $2 RETURNING approve_id`)
	if err != nil {
		return 0, fmt.Errorf("can't prepare statement: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id, phone)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("can't cancel approve: %w", err)
	}

	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("can't get canceled id: %w", err)
	}

	return id, nil
}

func correctTime(timestamp time.Time) string {
	return timestamp.Format(time.DateOnly)
}

func supInfo(reserved []models.Reserved) []models.Sup {
	sups := make([]models.Sup, 3)

	for _, res := range reserved {
		sups[res.ModelID-1].ID = res.ModelID
		sups[res.ModelID-1].Name = res.ModelName
		sups[res.ModelID-1].Price = res.ModelPrice
		sups[res.ModelID-1].ReservedDays = append(sups[res.ModelID].ReservedDays, res.Day)
	}

	return sups
}

func infoToJSON(supsInfo []models.ApproveSup) (string, error) {
	data, err := json.Marshal(supsInfo)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func infoFromJSON(data string) ([]models.ApproveSup, error) {
	res := []models.ApproveSup{}

	if err := json.Unmarshal([]byte(data), &res); err != nil {
		return nil, err
	}

	return res, nil
}
