package repositories

import (
	"database/sql"
	"warehouse-api/models"
)

type PenjualanRepository interface {
	Create(tx *sql.Tx, header *models.JualHeader, details []models.JualDetail) error
	GetAll(startDate, endDate string) ([]models.JualHeader, error)
	GetByID(id int) (*models.JualHeader, error)
}

type penjualanRepository struct {
	db *sql.DB
}

func NewPenjualanRepository(db *sql.DB) PenjualanRepository {
	return &penjualanRepository{db}
}

func (r *penjualanRepository) Create(tx *sql.Tx, header *models.JualHeader, details []models.JualDetail) error {
    // Insert Header
    queryHeader := `INSERT INTO jual_header (no_faktur, customer, total, user_id, status) 
                    VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
    err := tx.QueryRow(queryHeader, header.NoFaktur, header.Customer, header.Total, header.UserID, header.Status).Scan(&header.ID, &header.CreatedAt)
    if err != nil {
        return err
    }

    // Insert Details
    queryDetail := `INSERT INTO jual_detail (jual_header_id, barang_id, qty, harga, subtotal) 
                    VALUES ($1, $2, $3, $4, $5)`
    for _, d := range details {
        _, err := tx.Exec(queryDetail, header.ID, d.BarangID, d.Qty, d.Harga, d.Subtotal)
        if err != nil {
            return err
        }
    }

    return nil
}

func (r *penjualanRepository) GetAll(startDate, endDate string) ([]models.JualHeader, error) {
    query := `SELECT h.id, h.no_faktur, h.customer, h.total, h.user_id, h.status, h.created_at, u.username 
              FROM jual_header h
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

    var headers []models.JualHeader
    for rows.Next() {
        var h models.JualHeader
        h.User = &models.User{}
        if err := rows.Scan(&h.ID, &h.NoFaktur, &h.Customer, &h.Total, &h.UserID, &h.Status, &h.CreatedAt, &h.User.Username); err != nil {
            return nil, err
        }
        headers = append(headers, h)
    }
    return headers, nil
}

func (r *penjualanRepository) GetByID(id int) (*models.JualHeader, error) {
    queryHeader := `SELECT h.id, h.no_faktur, h.customer, h.total, h.user_id, h.status, h.created_at, u.username 
                    FROM jual_header h
                    JOIN users u ON h.user_id = u.id
                    WHERE h.id = $1`
    var h models.JualHeader
    h.User = &models.User{}
    err := r.db.QueryRow(queryHeader, id).Scan(&h.ID, &h.NoFaktur, &h.Customer, &h.Total, &h.UserID, &h.Status, &h.CreatedAt, &h.User.Username)
    if err != nil {
        return nil, err
    }

    queryDetails := `SELECT d.id, d.barang_id, d.qty, d.harga, d.subtotal, b.kode_barang, b.nama_barang, b.satuan
                     FROM jual_detail d
                     JOIN master_barang b ON d.barang_id = b.id
                     WHERE d.jual_header_id = $1`
    rows, err := r.db.Query(queryDetails, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var d models.JualDetail
        d.Barang = &models.Barang{}
        if err := rows.Scan(&d.ID, &d.BarangID, &d.Qty, &d.Harga, &d.Subtotal, &d.Barang.KodeBarang, &d.Barang.NamaBarang, &d.Barang.Satuan); err != nil {
            return nil, err
        }
        h.Details = append(h.Details, d)
    }

    return &h, nil
}
