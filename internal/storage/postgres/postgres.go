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

var reverseRelationships = map[string]string{
	"grandmother": "grandson",
	"grandson":    "grandmother",
	"mother":      "child",
	"child":       "parent",
	"father":      "child",
	"daughter":    "parent",
	"son":         "parent",
	"brother":     "sibling",
	"sister":      "sibling",
}

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

func (s *Storage) AddTempUser(tempUser *models.TempUser) error {
	const op = "storage.postgres.NewAddTempUser"

	q := `INSERT INTO temp_users (email, hashed_password, token, created_at) VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(q, tempUser.Email, tempUser.HashedPassword, tempUser.Token, time.Now().UTC())
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) TempUserByEmail(email string) (*models.TempUser, error) {
	const op = "storage.postgres.TempUserByEmail"

	q := `SELECT * FROM temp_users WHERE email = $1`

	var user models.TempUser
	err := s.db.QueryRow(q, email).Scan(
		&user.Id,
		&user.Email,
		&user.HashedPassword,
		&user.Token,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil

}

func (s *Storage) UpdateTempUser(tempUser *models.TempUser) error {
	const op = "storage.postgres.UpdateTempUser"
	q := `UPDATE temp_users SET token=$1, created_at=$2 WHERE email=$3`

	res, err := s.db.Exec(q, tempUser.Token, time.Now().UTC(), tempUser.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}

	return nil
}

func (s *Storage) DeleteTempUser(id int64) error {
	const op = "storage.postgres.DeleteTempUser"

	q := `DELETE FROM temp_users WHERE id = $1`
	_, err := s.db.Exec(q, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) AddUser(user *models.User) error {
	const op = "storage.postgres.AddUser"

	q := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3)`

	_, err := s.db.Exec(q, user.Email, user.HashedPassword, user.Role)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UserByEmail(email string) (*models.User, error) {
	const op = "storage.postgres.NewUserByEmail"

	q := `SELECT * FROM users WHERE email = $1`

	var user models.User
	err := s.db.QueryRow(q, email).Scan(
		&user.Id,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) UpdatePasswordUser(id int64, password string) error {
	const op = "storage.postgres.ChangePasswordUser"

	q := `UPDATE users SET password = $1 WHERE id = $2`

	_, err := s.db.Exec(q, password, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) AddPatient(patient *models.Patient) error {
	const op = "storage.postgres.AddPatient"

	q := `INSERT INTO patients (
		name, 
		surname, 
		birth_date,
        email,
		date_of_baseline_visit, 
		age_visit_baseline, 
		hypertension_baseline, 
		diabetes_baseline, 
		smoking_status_baseline, 
		cvd_baseline, 
		cad_baseline, 
		mi_baseline, 
		cad_revascularization_baseline, 
		stroke_baseline, 
		stroke_premature_baseline, 
		liver_steatosis_baseline, 
		xanthelasma, 
		weight_baseline, 
		height_baseline, 
		thyroid_disease, 
		menarche_reached_baseline, 
		age_at_menarche_baseline, 
		menopause_reached_baseline
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23
	)`

	_, err := s.db.Exec(q,
		patient.Name,
		patient.Surname,
		patient.BirthDate,
		patient.Email,
		patient.DateOfBaselineVisit,
		patient.AgeVisitBaseline,
		patient.HypertensionBaseline,
		patient.DiabetesBaseline,
		patient.SmokingStatusBaseline,
		patient.CVDBaseline,
		patient.CADBaseline,
		patient.MIBaseline,
		patient.CADRevascularization,
		patient.StrokeBaseline,
		patient.StrokePremature,
		patient.LiverSteatosisBaseline,
		patient.Xanthelasma,
		patient.WeightBaseline,
		patient.HeightBaseline,
		patient.ThyroidDisease,
		patient.MenarcheBaseline,
		patient.AgeMenarche,
		patient.MenopauseBaseline,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrPatientExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) PatientByEmail(email string) (*models.Patient, error) {
	const op = "storage.postgres.PatientByEmail"

	q := `SELECT * FROM patients WHERE email = $1`

	var patient models.Patient

	err := s.db.QueryRow(q, email).Scan(
		&patient.ID,
		&patient.Name,
		&patient.Surname,
		&patient.BirthDate,
		&patient.Email,
		&patient.DateOfBaselineVisit,
		&patient.AgeVisitBaseline,
		&patient.HypertensionBaseline,
		&patient.DiabetesBaseline,
		&patient.SmokingStatusBaseline,
		&patient.CVDBaseline,
		&patient.CADBaseline,
		&patient.MIBaseline,
		&patient.CADRevascularization,
		&patient.StrokeBaseline,
		&patient.StrokePremature,
		&patient.LiverSteatosisBaseline,
		&patient.Xanthelasma,
		&patient.WeightBaseline,
		&patient.HeightBaseline,
		&patient.ThyroidDisease,
		&patient.MenarcheBaseline,
		&patient.AgeMenarche,
		&patient.MenopauseBaseline,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrPatientNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &patient, nil
}

func (s *Storage) PatientById(id int64) (*models.Patient, error) {
	const op = "storage.postgres.PatientByEmail"

	q := `SELECT * FROM patients WHERE id = $1`

	var patient models.Patient

	err := s.db.QueryRow(q, id).Scan(
		&patient.ID,
		&patient.Name,
		&patient.Surname,
		&patient.BirthDate,
		&patient.Email,
		&patient.DateOfBaselineVisit,
		&patient.AgeVisitBaseline,
		&patient.HypertensionBaseline,
		&patient.DiabetesBaseline,
		&patient.SmokingStatusBaseline,
		&patient.CVDBaseline,
		&patient.CADBaseline,
		&patient.MIBaseline,
		&patient.CADRevascularization,
		&patient.StrokeBaseline,
		&patient.StrokePremature,
		&patient.LiverSteatosisBaseline,
		&patient.Xanthelasma,
		&patient.WeightBaseline,
		&patient.HeightBaseline,
		&patient.ThyroidDisease,
		&patient.MenarcheBaseline,
		&patient.AgeMenarche,
		&patient.MenopauseBaseline,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrPatientNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &patient, nil
}

func (s *Storage) PatientsAll() ([]*models.Patient, error) {
	const op = "storage.postgres.AllPatients"

	q := `SELECT * FROM patients`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var patients []*models.Patient

	for rows.Next() {
		var patient models.Patient
		err := rows.Scan(
			&patient.ID,
			&patient.Name,
			&patient.Surname,
			&patient.BirthDate,
			&patient.Email,
			&patient.DateOfBaselineVisit,
			&patient.AgeVisitBaseline,
			&patient.HypertensionBaseline,
			&patient.DiabetesBaseline,
			&patient.SmokingStatusBaseline,
			&patient.CVDBaseline,
			&patient.CADBaseline,
			&patient.MIBaseline,
			&patient.CADRevascularization,
			&patient.StrokeBaseline,
			&patient.StrokePremature,
			&patient.LiverSteatosisBaseline,
			&patient.Xanthelasma,
			&patient.WeightBaseline,
			&patient.HeightBaseline,
			&patient.ThyroidDisease,
			&patient.MenarcheBaseline,
			&patient.AgeMenarche,
			&patient.MenopauseBaseline,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		patients = append(patients, &patient)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return patients, nil
}

func (s *Storage) DeletePatient(id int64) error {
	const op = "storage.postgres.DeletePatient"

	q := `DELETE FROM patients WHERE id = $1`
	_, err := s.db.Exec(q, id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) AddRelation(relation *models.Relation) error {
	const op = "storage.postgres.AddRelation"

	q := `INSERT INTO relatives (patient_id, relative_id, relationship_type, created_at) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(q, relation.PatientID, relation.RelativeID, relation.RelationshipType, time.Now().UTC())
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrRelationExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) Relation(id int64) (*models.Relation, error) {
	const op = "storage.postgres.Relation"

	q := `SELECT * FROM relatives WHERE id = $1`

	var relation models.Relation
	err := s.db.QueryRow(q, id).Scan(
		&relation.ID,
		&relation.PatientID,
		&relation.RelativeID,
		&relation.RelationshipType,
		&relation.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrRelationNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &relation, nil
}

func (s *Storage) AllRelations(patientID int64) ([]models.Relation, error) {
	const op = "storage.postgres.AllRelations"

	q := `SELECT id, patient_id, relative_id, relationship_type
		  FROM relatives
		  WHERE patient_id = $1 OR relative_id = $1`

	rows, err := s.db.Query(q, patientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var relations []models.Relation
	for rows.Next() {
		var relation models.Relation
		if err := rows.Scan(&relation.ID, &relation.PatientID, &relation.RelativeID, &relation.RelationshipType); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		// Проверяем, если patientID совпадает с RelativeID, меняем их местами и тип связи
		if relation.RelativeID == patientID {
			relation.PatientID, relation.RelativeID = relation.RelativeID, relation.PatientID
			if reverseType, exists := reverseRelationships[relation.RelationshipType]; exists {
				relation.RelationshipType = reverseType
			}
		}

		relations = append(relations, relation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return relations, nil
}

func (s *Storage) DeleteRelation(patId, relId int64) error {
	const op = "storage.postgres.DeletePatient"

	q := `DELETE FROM relatives WHERE patient_id = $1 AND relative_id = $2`
	_, err := s.db.Exec(q, patId, relId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrRelationNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
