package opsgenieclient

type Alert struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type AlertSummary map[string][]Alert

type CountSummary map[string]AlertCountSummary

type AlertCountSummary []AlertCountSummaryItem

type AlertCountSummaryItem struct {
	CurrentCount  Count
	PreviousCount Count
	Change        Change
	Display       string
}

type Change struct {
	Diff       Count
	Percentage Count
}

type Count struct {
	BusinessHours    int
	NonBusinessHours int
	Total            int
}

type Period struct {
	NumDays int
	Display string
}
