package serializers

import (
	"gitlab.com/almanac-app/models"
)

type PublicReportResponse struct {
	ReportContent models.ReportContent `json:"reportContent"`
	Date          string               `json:"date"`
}

type ReportSerializer struct {
	models.DailyReport
}

func (ser *ReportSerializer) PublicResponse() PublicReportResponse {
	response := PublicReportResponse{
		ReportContent: ser.ReportContent,
		Date:          ser.Date,
	}
	return response
}
