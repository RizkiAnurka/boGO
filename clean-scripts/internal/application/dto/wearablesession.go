package dto

import (
	"clean-scripts/internal/domain/model"
)

// WearableSession representing wearablesession dto
type WearableSession struct {
	ID         int64  `json:"id,omitempty"`
	RefID      string `json:"ref_id,omitempty"`
	ClientType int64  `json:"client_type,omitempty"`
	ClientID   string `json:"client_id,omitempty"`
	KeyTime    int64  `json:"key_time,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	MermaidID  string `json:"mermaid_id,omitempty"`
}

// WearableSessions representing collection of WearableSession
type WearableSessions []WearableSession

// Marshal converts DTO to domain model
func (d *WearableSession) Marshal() (model.WearableSession, error) {
	domainModel := model.WearableSession{
		MetaField:  model.MetaField{ID: d.ID},
		RefID:      d.RefID,
		ClientType: d.ClientType,
		ClientID:   d.ClientID,
		KeyTime:    d.KeyTime,
		DeviceName: d.DeviceName,
		UserID:     d.UserID,
		MermaidID:  d.MermaidID,
	}

	return domainModel, nil
}

// Unmarshal converts domain model to DTO
func (d *WearableSession) Unmarshal(domainModel *model.WearableSession) {
	d.ID = domainModel.MetaField.ID
	d.RefID = domainModel.RefID
	d.ClientType = domainModel.ClientType
	d.ClientID = domainModel.ClientID
	d.KeyTime = domainModel.KeyTime
	d.DeviceName = domainModel.DeviceName
	d.UserID = domainModel.UserID
	d.MermaidID = domainModel.MermaidID
}

// Unmarshal converts slice of domain models to DTOs
func (d *WearableSessions) Unmarshal(domainModels []model.WearableSession) {
	for _, domainModel := range domainModels {
		var dto WearableSession
		dto.Unmarshal(&domainModel)
		*d = append(*d, dto)
	}
}
