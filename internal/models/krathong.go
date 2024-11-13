package models

import (
	"time"
)

// Krathong represents the structure of the Krathong model
type Krathong struct {
	KrathongID    uint      `gorm:"primaryKey" json:"krathong_id"`
	KrathongType  int       `gorm:"type:int;not null" json:"krathong_type"`
	EmpName       string    `gorm:"type:varchar(50);not null" json:"emp_name"`
	EmpDepartment string    `gorm:"type:varchar(50);not null" json:"emp_department"`
	EmpWish       string    `gorm:"type:text" json:"emp_wish"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}
