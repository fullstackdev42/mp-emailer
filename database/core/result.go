package core

import "gorm.io/gorm"

// GormResult wraps GORM's *gorm.Statement to implement our Result interface
type GormResult struct {
	statement *gorm.Statement
}

// NewResult creates a new Result instance
func NewResult(statement *gorm.Statement) Result {
	return &GormResult{statement: statement}
}

// Scan implements Result.Scan
func (r *GormResult) Scan(dest interface{}) Result {
	r.statement.Scan(dest)
	return r
}

// Error implements Result.Error
func (r *GormResult) Error() error {
	return r.statement.Error
}
