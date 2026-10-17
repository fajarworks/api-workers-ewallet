package model

type DashboardSummary struct {
	Balance int64 `db:"balance"`
	Income  int64 `db:"income"`
	Expense int64 `db:"expense"`
}
