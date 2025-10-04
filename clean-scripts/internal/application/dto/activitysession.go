package dto

import (
	"clean-scripts/internal/domain/model"
)

// ActivitySession representing activitysession dto
type ActivitySession struct {
	ID           int64   `json:"id,omitempty"`
	UserID       int64   `json:"user_id,omitempty"`
	SessionID    int64   `json:"session_id,omitempty"`
	Timestamp    int64   `json:"timestamp,omitempty"`
	ActivityFlag int64   `json:"activity_flag,omitempty"`
	Battery      int64   `json:"battery,omitempty"`
	StepCount    int64   `json:"step_count,omitempty"`
	Distance     float64 `json:"distance,omitempty"`
	ActiveMins   int64   `json:"active_mins,omitempty"`
	Calories     float64 `json:"calories,omitempty"`
	Interval     int64   `json:"interval,omitempty"`
	Tatpim       float64 `json:"tatpim,omitempty"`
	Bmr          float64 `json:"bmr,omitempty"`
}

// ActivitySessions representing collection of ActivitySession
type ActivitySessions []ActivitySession

// Marshal converts DTO to domain model
func (d *ActivitySession) Marshal() (model.ActivitySession, error) {
	domainModel := model.ActivitySession{
		MetaField:    model.MetaField{ID: d.ID},
		UserID:       d.UserID,
		SessionID:    d.SessionID,
		Timestamp:    d.Timestamp,
		ActivityFlag: d.ActivityFlag,
		Battery:      d.Battery,
		StepCount:    d.StepCount,
		Distance:     d.Distance,
		ActiveMins:   d.ActiveMins,
		Calories:     d.Calories,
		Interval:     d.Interval,
		Tatpim:       d.Tatpim,
		Bmr:          d.Bmr,
	}

	return domainModel, nil
}

// Unmarshal converts domain model to DTO
func (d *ActivitySession) Unmarshal(domainModel *model.ActivitySession) {
	d.ID = domainModel.MetaField.ID
	d.UserID = domainModel.UserID
	d.SessionID = domainModel.SessionID
	d.Timestamp = domainModel.Timestamp
	d.ActivityFlag = domainModel.ActivityFlag
	d.Battery = domainModel.Battery
	d.StepCount = domainModel.StepCount
	d.Distance = domainModel.Distance
	d.ActiveMins = domainModel.ActiveMins
	d.Calories = domainModel.Calories
	d.Interval = domainModel.Interval
	d.Tatpim = domainModel.Tatpim
	d.Bmr = domainModel.Bmr
}

// Unmarshal converts slice of domain models to DTOs
func (d *ActivitySessions) Unmarshal(domainModels []model.ActivitySession) {
	for _, domainModel := range domainModels {
		var dto ActivitySession
		dto.Unmarshal(&domainModel)
		*d = append(*d, dto)
	}
}
