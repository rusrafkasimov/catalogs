package models

import "time"

type OperationType string
type OperationMethod string

type Operation struct {
	Type      OperationType   `json:"type"`
	Method    OperationMethod `json:"method"`
	Catalog   *Catalog        `json:"catalog,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

const (
	OperationTypeCatalogs OperationType = "catalogs"

	OperationMethodUpsert OperationMethod = "upsert"
	OperationMethodDelete OperationMethod = "delete"
)
