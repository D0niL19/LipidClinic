-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS relatives (
    id SERIAL PRIMARY KEY,
    patient_id INT REFERENCES patients(id) ON DELETE CASCADE,  -- пациент
    relative_id INT REFERENCES patients(id) ON DELETE CASCADE,  -- родственник
    relationship_type VARCHAR(20),  -- тип связи (например, отец, мать, брат, сестра и т.д.)
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS relatives;
-- +goose StatementEnd
