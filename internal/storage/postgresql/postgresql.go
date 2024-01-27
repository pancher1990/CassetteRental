package postgresql

import (
	"CassetteRental/internal/config"
	"CassetteRental/internal/lib/random"
	"CassetteRental/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(cfgS config.StorageDb) (*Storage, error) {
	const op = "storage.postgresql.New"

	connStr := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		cfgS.Host, cfgS.Port, cfgS.User, cfgS.Password, cfgS.DbName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func rollback(tx *sql.Tx, err error) error {
	if errTx := tx.Rollback(); errTx != nil {
		return fmt.Errorf("rollback failed: %v, original error: %w", errTx, err)
	}
	return err
}

func (s *Storage) getTransaction(ctx context.Context) (context.Context, *sql.Tx, error) {
	const op = "storage.postgresql.getTransaction"

	var tr *sql.Tx
	var err error

	if ctx.Value("tr") == nil {
		tr, err = s.db.Begin()
		if err != nil {
			return ctx, nil, fmt.Errorf("%s: %w", op, err)
		}
		ctx = context.WithValue(ctx, "tr", tr)
	} else {
		tr = ctx.Value("tr").(*sql.Tx)
	}

	return ctx, tr, nil
}
func (s *Storage) commitIfNeed(ctx context.Context, tr *sql.Tx) error {
	const op = "storage.postgresql.commitIfNeed"

	if ctx.Value("returnTransaction") == nil {
		if err := tr.Commit(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
func (s *Storage) AddNewCustomer(name string, isActive bool, balance int) (string, error) {
	const op = "storage.postgresql.AddNewCustomer"

	id := random.GenerateGUID()
	stmt, err := s.db.Prepare("INSERT INTO customer (id, name, is_active, balance) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id, name, isActive, balance)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) SetCustomerBalance(ctx context.Context, id string, balance int) (context.Context, error) {
	const op = "storage.postgresql.SetCustomerBalance"

	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tr.Prepare("UPDATE customer SET balance = $1 WHERE id = $2")
	if err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	_, err = stmt.Exec(balance, id)
	if err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))

	}

	return ctx, nil
}
func (s *Storage) GetCustomerBalance(ctx context.Context, id string) (context.Context, int, error) {
	const op = "storage.postgresql.GetCustomerBalance"

	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tr.Prepare("SELECT (balance) FROM customer WHERE id = $1")
	if err != nil {
		return ctx, 0, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	var balance int
	if err = stmt.QueryRow(id).Scan(&balance); errors.Is(err, sql.ErrNoRows) {
		return ctx, 0, rollback(tr, storage.ErrCustomerNotFound)
	}
	if err != nil {
		return ctx, 0, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, 0, rollback(tr, fmt.Errorf("%s: %w", op, err))

	}

	return ctx, balance, nil
}

func (s *Storage) AddNewFilm(name string, dayPrice int) (string, error) {
	const op = "storage.postgresql.AddNewFilm"

	id := random.GenerateGUID()
	stmt, err := s.db.Prepare("INSERT INTO film (id, title, day_price) VALUES ($1, $2, $3)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id, name, dayPrice)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetFilm(ctx context.Context, title string) (context.Context, string, int, error) {
	const op = "storage.postgresql.GetFilm"

	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, "", 0, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tr.Prepare("SELECT id, day_price FROM film WHERE title = $1")
	if err != nil {
		return ctx, "", 0, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	var id string
	var dayCost int
	if err = stmt.QueryRow(title).Scan(&id, &dayCost); errors.Is(err, sql.ErrNoRows) {
		return ctx, "", 0, rollback(tr, storage.ErrFilmNotFound)
	}
	if err != nil {
		return ctx, "", 0, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, "", 0, rollback(tr, fmt.Errorf("%s: %w", op, err))

	}

	return ctx, id, dayCost, nil
}

func (s *Storage) AddNewCassette(filmId string) (string, error) {
	const op = "storage.postgresql.AddNewCassette"

	id := random.GenerateGUID()
	stmt, err := s.db.Prepare("INSERT INTO cassette (id, film_id, available) VALUES ($1, $2, $3)")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(id, filmId, true)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetCassetteStatus(id string) (bool, error) {
	const op = "storage.postgresql.GetCassetteStatus"

	stmt, err := s.db.Prepare("SELECT (available)  FROM cassette WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	var available bool
	// TODO исправить работу с ошибками свернуть везде в один if по аналогии с тем как ниже
	if err = stmt.QueryRow(id).Scan(&available); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return available, nil
}

func (s *Storage) FindAvailableCassette(ctx context.Context, filmId string) (context.Context, string, error) {
	const op = "storage.postgresql.FindAvailableCassette"

	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, "", fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tr.Prepare("SELECT (id)  FROM cassette WHERE film_id = $1 and available = true LIMIT 1")
	if err != nil {
		return ctx, "", fmt.Errorf("%s: %w", op, err)
	}

	var id string
	if err = stmt.QueryRow(filmId).Scan(&id); errors.Is(err, sql.ErrNoRows) {
		return ctx, "", storage.ErrCassetteNotFound
	}
	if err != nil {
		return ctx, "", fmt.Errorf("%s: %w", op, err)

	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, "", rollback(tr, fmt.Errorf("%s: %w", op, err))

	}

	return ctx, id, nil
}

func (s *Storage) SetCassetteStatus(ctx context.Context, id string, available bool) (context.Context, error) {
	const op = "storage.postgresql.SetCassetteStatus"

	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tr.Prepare("UPDATE cassette SET available = $1 WHERE id = $2")
	if err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	_, err = stmt.Exec(available, id)
	if err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))

	}

	return ctx, nil
}

func (s *Storage) CreateOrder(ctx context.Context, customerId string) (context.Context, string, error) {
	const op = "storage.postgresql.CreateOrder"
	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, "", fmt.Errorf("%s: %w", op, err)
	}

	id := random.GenerateGUID()
	currentTime := time.Now()
	stmt, err := tr.Prepare("INSERT INTO \"order\" (id, customer_id, date) VALUES ($1, $2, $3)")
	if err != nil {
		return ctx, "", rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	_, err = stmt.Exec(id, customerId, currentTime)
	if err != nil {
		return ctx, "", rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, "", rollback(tr, fmt.Errorf("%s: %w", op, err))

	}

	return ctx, id, nil
}

func (s *Storage) CreateCassetteInOrder(ctx context.Context, cassetteId string, orderId string, rentCost int) (context.Context, error) {
	const op = "storage.postgresql.CreateCassetteInOrder"

	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := tr.Prepare(
		"INSERT INTO cassette_in_order (cassette_id, order_id, rent_cost) VALUES ($1, $2, $3)")
	if err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	_, err = stmt.Exec(cassetteId, orderId, rentCost)
	if err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	return ctx, nil
}

func (s *Storage) CreateRent(ctx context.Context, customerId string, cassetteId string, rentDays int) (context.Context, string, error) {
	const op = "storage.postgresql.CreateRent"

	ctx, tr, err := s.getTransaction(ctx)
	if err != nil {
		return ctx, "", fmt.Errorf("%s: %w", op, err)
	}

	id := random.GenerateGUID()
	stmt, err := tr.Prepare(
		"INSERT INTO rent (id, customer_id, cassette_id, start_datetime, end_datetime, return_sign) " +
			"VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return ctx, "", rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	startDatetime := time.Now().Truncate(24 * time.Hour)
	endDatetime := startDatetime.AddDate(0, 0, rentDays)
	_, err = stmt.Exec(id, customerId, cassetteId, startDatetime, endDatetime, false)
	if err != nil {
		return ctx, "", rollback(tr, fmt.Errorf("%s: %w", op, err))

	}

	if err = s.commitIfNeed(ctx, tr); err != nil {
		return ctx, "", rollback(tr, fmt.Errorf("%s: %w", op, err))
	}

	return ctx, id, nil
}
