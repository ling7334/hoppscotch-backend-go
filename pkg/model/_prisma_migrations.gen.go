// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNamePrismaMigration = "_prisma_migrations"

// PrismaMigration mapped from table <_prisma_migrations>
type PrismaMigration struct {
	ID                string     `gorm:"column:id;type:character varying(36);primaryKey" json:"id"`
	Checksum          string     `gorm:"column:checksum;type:character varying(64);not null" json:"checksum"`
	FinishedAt        *time.Time `gorm:"column:finished_at;type:timestamp with time zone" json:"finished_at"`
	MigrationName     string     `gorm:"column:migration_name;type:character varying(255);not null" json:"migration_name"`
	Logs              *string    `gorm:"column:logs;type:text" json:"logs"`
	RolledBackAt      *time.Time `gorm:"column:rolled_back_at;type:timestamp with time zone" json:"rolled_back_at"`
	StartedAt         time.Time  `gorm:"column:started_at;type:timestamp with time zone;not null;default:now()" json:"started_at"`
	AppliedStepsCount int32      `gorm:"column:applied_steps_count;type:integer;not null" json:"applied_steps_count"`
}

// TableName PrismaMigration's table name
func (*PrismaMigration) TableName() string {
	return TableNamePrismaMigration
}
