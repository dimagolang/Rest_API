package repository

import (
	"Rest_API/internal/models"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

// создать структуру FlightsRepo у которой будет подключение к БД
// инициализировать структуру NewFlightsRepo
// добавить метод получения всех рейсов
// в мейне ее создать и передать в сервис flight

type FlightsRepo struct {
	db *pgx.Conn // или *pgxpool.Pool, если ты используешь пул
}

func NewFlightsRepo(db *pgx.Conn) *FlightsRepo {
	return &FlightsRepo{db: db}
}

func (r *FlightsRepo) GetAllFlightsFromDB(ctx context.Context) ([]models.Flight, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM public.flights")
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch flights from database")
	}
	defer rows.Close()

	var flights []models.Flight
	for rows.Next() {
		var flight models.Flight
		if err := rows.Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo, &flight.DeleteAt); err != nil {
			return nil, errors.Wrap(err, "error scanning flights")
		}
		flights = append(flights, flight)
	}

	if len(flights) == 0 {
		return nil, pgx.ErrNoRows // ✅ Return pgx.ErrNoRows if no active flights exist
	}

	return flights, nil
}

func (r *FlightsRepo) InsertFlightToDB(ctx context.Context, flight *models.Flight) error {
	err := r.db.QueryRow(ctx,
		"INSERT INTO public.flights (destination_from, destination_to, delete_at) VALUES ($1, $2, $3) RETURNING id",
		flight.DestinationFrom, flight.DestinationTo, flight.DeleteAt,
	).Scan(&flight.FlightID)

	if err != nil {
		return errors.Wrap(err, "failed to insert flight")
	}
	return nil
}

func (r *FlightsRepo) GetFlightByIDFromDB(ctx context.Context, id int) (*models.Flight, error) {
	var flight models.Flight
	err := r.db.QueryRow(ctx,
		"SELECT id, destination_from, destination_to, delete_at FROM public.flights WHERE id = $1 AND delete_at = 0",
		id,
	).Scan(&flight.FlightID, &flight.DestinationFrom, &flight.DestinationTo, &flight.DeleteAt)

	if err != nil {
		return nil, errors.Wrap(err, "flight not found")
	}
	return &flight, nil
}

func (r *FlightsRepo) UpdateFlightInDB(ctx context.Context, flight *models.Flight) error {
	commandTag, err := r.db.Exec(ctx,
		"UPDATE public.flights SET destination_from = $1, destination_to = $2 WHERE id = $3",
		flight.DestinationFrom, flight.DestinationTo, flight.FlightID)

	if err != nil {
		return errors.Wrap(err, "failed to update flight")
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows // ничего не обновлено (возможно, soft deleted)
	}
	return nil
}

func (r *FlightsRepo) DeleteFlightFromDB(ctx context.Context, id int) error {
	commandTag, err := r.db.Exec(ctx,
		"UPDATE public.flights SET delete_at = EXTRACT(EPOCH FROM NOW())::BIGINT WHERE id = $1 AND delete_at = 0",
		id)

	if err != nil {
		return errors.Wrap(err, "failed to delete flight")
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows // ничего не удалено (возможно, уже удалено)
	}
	return nil
}
