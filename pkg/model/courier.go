package model

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Point struct {
	Lng float64
	Lat float64
}

func (p Point) GormDataType() string {
	return "point"
}

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []interface{}{fmt.Sprintf("POINT(%f %f)", p.Lng, p.Lat)},
	}
}

func (p *Point) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	// Format: POINT(longitude latitude)
	format := fmt.Sprintf("POINT(%f %f)", p.Lng, p.Lat)
	return format, nil
}

func (p *Point) Scan(value interface{}) error {
	if value == nil {
		*p = Point{}
		return nil
	}

	s, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("gagal scan Point, tipe tidak valid: %T", value)
	}

	fmt.Println(string(s))

	return nil

}

// Courier represents the couriers table
type Courier struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    *uint          `json:"user_id" gorm:"index"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	Latitude  *float64       `json:"latitude" gorm:"type:double"`
	Longitude *float64       `json:"longitude" gorm:"type:double"`
	Location  Point          `json:"location" gorm:"type:point"`
	CreatedAt *time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt *time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:CourierID"`
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the Courier model
func (Courier) TableName() string {
	return "couriers"
}
