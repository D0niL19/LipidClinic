package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")

	ErrRelationExists   = errors.New("relation already exists")
	ErrRelationNotFound = errors.New("relation not found")

	ErrPatientExists   = errors.New("relation already exists")
	ErrPatientNotFound = errors.New("patient not found")
)
