package availability

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/types"
)

type CreateAvailabilityRequest struct {
	Weekday   string `json:"weekday" binding:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
	StartTime string `json:"start_time" binding:"required,len=5"` // Format: HH:MM
	EndTime   string `json:"end_time" binding:"required,len=5"`   // Format: HH:MM
}

type AvailabilityResponse struct {
	UUID      types.BinaryUUID `json:"uuid"`
	Weekday   string           `json:"weekday"`
	StartTime string           `json:"start_time"`
	EndTime   string           `json:"end_time"`
}

func ToAvailabilityResponse(availability *models.Availability) *AvailabilityResponse {
	return &AvailabilityResponse{
		UUID:      availability.UUID,
		Weekday:   availability.Weekday,
		StartTime: availability.StartTime,
		EndTime:   availability.EndTime,
	}
}

func ToAvailabilityResponseList(availabilities []*models.Availability) []*AvailabilityResponse {
	var response []*AvailabilityResponse
	for _, availability := range availabilities {
		response = append(response, ToAvailabilityResponse(availability))
	}
	return response
}
