package models

import "time"

type Relation struct {
	ID               int64     `json:"id"`
	PatientID        int64     `json:"patient_id" binding:"required"`
	RelativeID       int64     `json:"relative_id" binding:"required"`
	RelationshipType string    `json:"relationship_type" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
}
