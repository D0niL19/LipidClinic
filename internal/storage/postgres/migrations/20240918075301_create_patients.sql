-- +goose Up
-- +goose StatementBegin
CREATE TABLE patients (
    id SERIAL PRIMARY KEY,
    name varchar(50),
    surname varchar(50),
    birth_date DATE,
    email VARCHAR(50) UNIQUE ,
    date_of_baseline_visit DATE,
    age_visit_baseline INT,
    hypertension_baseline BOOLEAN,
    diabetes_baseline BOOLEAN,
    smoking_status_baseline VARCHAR(50),
    cvd_baseline BOOLEAN,
    cad_baseline BOOLEAN,
    mi_baseline BOOLEAN,
    cad_revascularization_baseline BOOLEAN,
    stroke_baseline BOOLEAN,
    stroke_premature_baseline BOOLEAN,
    liver_steatosis_baseline BOOLEAN,
    xanthelasma BOOLEAN,
    weight_baseline DECIMAL(5,2),
    height_baseline DECIMAL(5,2),
    thyroid_disease BOOLEAN,
    menarche_reached_baseline BOOLEAN,
    age_at_menarche_baseline INT,
    menopause_reached_baseline BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS patients;
-- +goose StatementEnd
