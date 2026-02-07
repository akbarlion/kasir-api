package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetDailyReport() (*models.Report, error) {
	today := time.Now().Format("2006-01-02")
	return r.GetReportByDateRange(today, today)
}

func (r *ReportRepository) GetReportByDateRange(startDate, endDate string) (*models.Report, error) {
	var report models.Report

	// Get total revenue and transaction count
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaction)
	if err != nil {
		return nil, err
	}

	// Get best selling product
	err = r.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0)
		FROM transaction_details td
		JOIN product p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`, startDate, endDate).Scan(&report.ProdukTerlaris.Nama, &report.ProdukTerlaris.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &report, nil
}
