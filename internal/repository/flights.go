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

func (r *FlightsRepo) GetAll(ctx context.Context) ([]models.Flight, error) {
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
