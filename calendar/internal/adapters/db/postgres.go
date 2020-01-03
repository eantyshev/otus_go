package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/eantyshev/otus_go/calendar/internal/entity"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"time"
)

// implements entity.Repository + sync.Locker
type PgRepo struct {
	db *sql.DB
}

func NewPgRepo(dsn string) (*PgRepo, error) {
	db, err := sql.Open("pgx", dsn) // *sql.DB
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db ping: %s", err)
	}
	return &PgRepo{db: db}, nil
}

func (pgr *PgRepo) GetById(ctx context.Context, uid *uuid.UUID) (ap *appointment.Appointment, err error) {
	query := `
		SELECT summary, description, owner, time_start, time_end
		FROM appointment WHERE uuid = $1
	`
	row := pgr.db.QueryRowContext(ctx, query, uid)
	ap = &appointment.Appointment{Uuid: *uid}
	err = row.Scan(&ap.Summary, &ap.Description, &ap.Owner, &ap.TimeStart, &ap.TimeEnd)
	if err == sql.ErrNoRows {
		return nil, appointment.ErrIdNotFound
	} else if err != nil {
		return nil, err
	}
	return ap, nil
}

func (pgr *PgRepo) Create(ctx context.Context, ap *appointment.Appointment) (*uuid.UUID, error) {
	query := `
		INSERT INTO appointment(summary, description, owner, time_start, time_end)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uuid
	`
	var uidStr string
	err := pgr.db.QueryRowContext(ctx, query,
		ap.Summary,
		ap.Description,
		ap.Owner,
		ap.TimeStart.Format(time.RFC3339),
		ap.TimeEnd.Format(time.RFC3339),
	).Scan(&uidStr)
	if err != nil {
		return nil, err
	}
	uid, err := uuid.Parse(uidStr)
	if err != nil {
		return nil, err
	}
	return &uid, nil
}

func (pgr *PgRepo) Update(ctx context.Context, ap *appointment.Appointment) error {
	query := `
		UPDATE appointment
		SET summary = $1, description = $2, owner = $3, time_start = $4, time_end = $5
		WHERE uuid = $6
	`
	result, err := pgr.db.ExecContext(ctx, query,
		ap.Summary,
		ap.Description,
		ap.Owner,
		ap.TimeStart.Format(time.RFC3339),
		ap.TimeEnd.Format(time.RFC3339),
		ap.Uuid.String(),
	)
	if err != nil {
		return err
	}
	if rowCnt, err := result.RowsAffected(); err != nil {
		return err
	} else if rowCnt == 0 {
		return appointment.ErrIdNotFound
	}
	return nil
}

func (pgr *PgRepo) Delete(ctx context.Context, uid *uuid.UUID) error {
	query := `DELETE FROM appointment WHERE uuid = $1`
	result, err := pgr.db.ExecContext(ctx, query, uid)
	if err != nil {
		return err
	}
	rowCnt, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return appointment.ErrIdNotFound
	}
	return nil
}

func (pgr *PgRepo) ListOwnerPeriod(
	ctx context.Context,
	owner string,
	timeFrom time.Time,
	timeTo time.Time,
) (aps []*appointment.Appointment, err error) {
	query := `SELECT uuid, summary, description, owner, time_start, time_end
				FROM appointment
				WHERE owner = $1 AND time_start >= $2 AND time_end <= $3`
	rows, err := pgr.db.QueryContext(ctx, query, owner, timeFrom, timeTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ap := appointment.Appointment{}
		err = rows.Scan(&ap.Uuid, &ap.Summary, &ap.Description, &ap.Owner, &ap.TimeStart, &ap.TimeEnd)
		if err != nil {
			break
		}
		aps = append(aps, &ap)
	}
	err = rows.Err()
	return aps, err
}

func (pgr *PgRepo) FetchPeriod(
	ctx context.Context,
	timeFrom time.Time,
	timeTo time.Time,
) (aps []*appointment.Appointment, err error) {
	query := `SELECT uuid, summary, description, owner, time_start, time_end
				FROM appointment
				WHERE time_start >= $1 AND time_start < $2
				ORDER BY time_start
				`
	rows, err := pgr.db.QueryContext(ctx, query, timeFrom.Format(time.RFC3339), timeTo.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		ap := appointment.Appointment{}
		err = rows.Scan(&ap.Uuid, &ap.Summary, &ap.Description, &ap.Owner, &ap.TimeStart, &ap.TimeEnd)
		if err != nil {
			break
		}
		aps = append(aps, &ap)
	}
	err = rows.Err()
	return aps, err
}
