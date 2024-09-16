package postgres

import (
	"LipidClinic/internal/config"
	"LipidClinic/internal/models"
	"LipidClinic/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Sslmode)

	//db, err := sql.Open("postgres", connStr)
	var err error
	var db *sql.DB
	for attempts := 1; attempts <= cfg.MaxAttempts; attempts++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			// Проверяем успешность подключения
			err = db.Ping()
			if err == nil {
				break // Выход из цикла, если подключение успешно
			}
		}

		log.Printf("%s: failed to connect to database, attempt %d/%d: %v", op, attempts, cfg.DB.MaxAttempts, err)

		// Задержка перед следующей попыткой
		time.Sleep(cfg.Delay)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = goose.Up(db, cfg.DB.MigrationsPath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddUser(user *models.User) error {
	const op = "storage.postgres.AddUser"

	q := `INSERT INTO users (first_name, last_name, date_of_birth, gender, blood_type) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.Exec(q, user.FirstName, user.LastName, user.DateOfBirth, user.Gender, user.BloodType)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteUser(id int) error {
	const op = "storage.postgres.DeleteUser"

	q := `DELETE FROM users WHERE id = $1`
	result, err := s.db.Exec(q, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверка количества удалённых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to check affected rows: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}

	return nil
}

func (s *Storage) GetUser(id int) (*models.User, error) {
	const op = "storage.postgres.GetUser"

	q := `SELECT * FROM users WHERE id = $1`
	var user models.User
	err := s.db.QueryRow(q, id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Gender,
		&user.BloodType,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}
