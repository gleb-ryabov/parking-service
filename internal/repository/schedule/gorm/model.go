package gorm

import (
	"time"
	"uuid"

	"github.com/lib/pq"

	"parking/internal/domain"
)

type model struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ParkingID uuid.UUID `gorm:"type:uuid;not null"`
	Active    bool      `gorm:"not null"`

	ValidFrom time.Time  `gorm:"type:date;not null"`
	ValidTo   *time.Time `gorm:"type:date"`

	StartTime string `gorm:"type:time;not null"`
	EndTime   string `gorm:"type:time;not null"`

	WeekDays pq.Int64Array `gorm:"type:smallint[];not null"`
	Cost     int64         `gorm:"not null"`
}

// TableName returns table name in database.
func (model) TableName() string {
	return "schedules"
}

func toModel(s domain.Schedule) model {
	weekDays := make(pq.Int64Array, len(s.WeekDays))
	for i, day := range s.WeekDays {
		weekDays[i] = int64(day)
	}

	return model{
		ID:        s.ID,
		ParkingID: s.ParkingID,
		Active:    s.Active,
		ValidFrom: s.ValidFrom,
		ValidTo:   s.ValidTo,
		StartTime: s.StartTime.Format("15:04:05"),
		EndTime:   s.EndTime.Format("15:04:05"),
		WeekDays:  weekDays,
		Cost:      s.Cost,
	}
}

func toDomain(m model) (domain.Schedule, error) {
	startTime, err := time.Parse("15:04:05", m.StartTime)
	if err != nil {
		return domain.Schedule{}, err
	}

	endTime, err := time.Parse("15:04:05", m.EndTime)
	if err != nil {
		return domain.Schedule{}, err
	}

	weekDays := make([]int, len(m.WeekDays))
	for i, day := range m.WeekDays {
		weekDays[i] = int(day)
	}

	return domain.Schedule{
		ID:        m.ID,
		ParkingID: m.ParkingID,
		Active:    m.Active,
		ValidFrom: m.ValidFrom,
		ValidTo:   m.ValidTo,
		StartTime: startTime,
		EndTime:   endTime,
		WeekDays:  weekDays,
		Cost:      m.Cost,
	}, nil
}
