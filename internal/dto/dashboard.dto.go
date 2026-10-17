package dto

type DashboardSummaryResponse struct {
	Balance int64 `json:"balance"`
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
}
