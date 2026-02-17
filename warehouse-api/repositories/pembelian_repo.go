package repositories

import (
	"database/sql"
	"warehouse-api/models"
)

type PembelianRepository interface {
	Create(tx *sql.Tx, header *models.BeliHeader, details []models.BeliDetail) error
	GetAll(startDate, endDate string) ([]models.BeliHeader, error)
	GetByID(id int) (*models.BeliHeader, error)
}

type pembelianRepository struct {
	db *sql.DB
}

func NewPembelianRepository(db *sql.DB) PembelianRepository {
	return &pembelianRepository{db}
}

func (r *pembelianRepository) Create(tx *sql.Tx, header *models.BeliHeader, details []models.BeliDetail) error {
	// Insert Header
	queryHeader := `INSERT INTO beli_header (no_faktur, supplier, total, user_id, status) 
                    VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := tx.QueryRow(queryHeader, header.NoFaktur, header.Supplier, header.Total, header.UserID, header.Status).Scan(&header.ID, &header.CreatedAt)
	if err != nil {
		return err
	}

	// Insert Details
	queryDetail := `INSERT INTO beli_detail (beli_header_id, barang_id, qty, harga, subtotal) 
                    VALUES ($1, $2, $3, $4, $5)`
	for _, d := range details {
		_, err := tx.Exec(queryDetail, header.ID, d.BarangID, d.Qty, d.Harga, d.Subtotal)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *pembelianRepository) GetAll(startDate, endDate string) ([]models.BeliHeader, error) {
	query := `SELECT h.id, h.no_faktur, h.supplier, h.total, h.user_id, h.status, h.created_at, u.username 
              FROM beli_header h
              JOIN users u ON h.user_id = u.id`
	
	var args []interface{}
	if startDate != "" && endDate != "" {
		query += " WHERE h.created_at BETWEEN $1 AND $2"
		args = append(args, startDate, endDate)
	}
	query += " ORDER BY h.created_at DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var headers []models.BeliHeader
	for rows.Next() {
		var h models.BeliHeader
		h.User = &models.User{}
		if err := rows.Scan(&h.ID, &h.NoFaktur, &h.Supplier, &h.Total, &h.UserID, &h.Status, &h.CreatedAt, &h.User.Username); err != nil {
			return nil, err
		}
		headers = append(headers, h)
	}
	return headers, nil
}

func (r *pembelianRepository) GetByID(id int) (*models.BeliHeader, error) {
	queryHeader := `SELECT h.id, h.no_faktur, h.supplier, h.total, h.user_id, h.status, h.created_at, u.username 
                    FROM beli_header h
                    JOIN users u ON h.user_id = u.id
                    WHERE h.id = $1`
	var h models.BeliHeader
	h.User = &models.User{}
	err := r.db.QueryRow(queryHeader, id).Scan(&h.ID, &h.NoFaktur, &h.Supplier, &h.Total, &h.UserID, &h.Status, &h.CreatedAt, &h.User.Username)
	if err != nil {
		return nil, err
	}

	queryDetails := `SELECT d.id, d.barang_id, d.qty, d.harga, d.subtotal, b.kode_barang, b.nama_barang, b.satuan
                     FROM beli_detail d
                     JOIN master_barang b ON d.barang_id = b.id
                     WHERE d.beli_header_id = $1`
	rows, err := r.db.Query(queryDetails, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.BeliDetail
		d.Barang = &models.Barang{}
		if err := rows.Scan(&d.ID, &d.BarangID, &d.Qty, &d.Harga, &d.Subtotal, &d.Barang.KodeBarang, &d.Barang.NamaBarang, &d.Barang.Satuan); err != nil {
			return nil, err
		}
		h.Details = append(h.Details, d)
	}

	return &h, nil
}
