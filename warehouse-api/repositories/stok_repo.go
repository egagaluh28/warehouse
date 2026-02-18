package repositories

import (
	"database/sql"
	"warehouse-api/models"
)

type StokRepository interface {
	GetAll() ([]models.Stok, error)
	GetByBarangID(barangID int) (*models.Stok, error)
	GetByBarangIDWithTx(tx *sql.Tx, barangID int) (*models.Stok, error)
    CreateOrUpdate(tx *sql.Tx, barangID, qtyChange int) error
	GetHistory(barangID int) ([]models.HistoryStok, error)
    CreateHistory(tx *sql.Tx, history *models.HistoryStok) error
}

type stokRepository struct {
	db *sql.DB
}

func NewStokRepository(db *sql.DB) StokRepository {
	return &stokRepository{db}
}

func (r *stokRepository) GetAll() ([]models.Stok, error) {
	query := `
        SELECT s.id, s.barang_id, s.stok_akhir, s.updated_at,
               b.kode_barang, b.nama_barang
        FROM mstok s
        JOIN master_barang b ON s.barang_id = b.id`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var stoks []models.Stok
    for rows.Next() {
        var s models.Stok
        s.Barang = &models.Barang{}
        if err := rows.Scan(&s.ID, &s.BarangID, &s.StokAkhir, &s.UpdatedAt, &s.Barang.KodeBarang, &s.Barang.NamaBarang); err != nil {
            return nil, err
        }
        stoks = append(stoks, s)
    }
    return stoks, nil
}

func (r *stokRepository) GetByBarangID(barangID int) (*models.Stok, error) {
    query := `
        SELECT s.id, s.barang_id, s.stok_akhir, s.updated_at,
               b.kode_barang, b.nama_barang, b.satuan, b.harga_jual
        FROM mstok s
        JOIN master_barang b ON s.barang_id = b.id
        WHERE s.barang_id = $1`
    
    var s models.Stok
    s.Barang = &models.Barang{}
    err := r.db.QueryRow(query, barangID).Scan(
        &s.ID, &s.BarangID, &s.StokAkhir, &s.UpdatedAt, 
        &s.Barang.KodeBarang, &s.Barang.NamaBarang, &s.Barang.Satuan, &s.Barang.HargaJual,
    )
    if err != nil {
        return nil, err
    }
    return &s, nil
}

// GetByBarangIDWithTx queries stock within a transaction (important for consistent audit trail)
func (r *stokRepository) GetByBarangIDWithTx(tx *sql.Tx, barangID int) (*models.Stok, error) {
    query := `
        SELECT s.id, s.barang_id, s.stok_akhir, s.updated_at,
               b.kode_barang, b.nama_barang, b.satuan, b.harga_jual
        FROM mstok s
        JOIN master_barang b ON s.barang_id = b.id
        WHERE s.barang_id = $1`
    
    var s models.Stok
    s.Barang = &models.Barang{}
    err := tx.QueryRow(query, barangID).Scan(
        &s.ID, &s.BarangID, &s.StokAkhir, &s.UpdatedAt, 
        &s.Barang.KodeBarang, &s.Barang.NamaBarang, &s.Barang.Satuan, &s.Barang.HargaJual,
    )
    if err != nil {
        return nil, err
    }
    return &s, nil
}

// CreateOrUpdate handles stock update logic. If passed a tx, it uses it.
func (r *stokRepository) CreateOrUpdate(tx *sql.Tx, barangID, qtyChange int) error {
    // Check if stock exists
    var exists bool
    checkQuery := "SELECT EXISTS(SELECT 1 FROM mstok WHERE barang_id=$1)"
    
    var err error
    if tx != nil {
       err = tx.QueryRow(checkQuery, barangID).Scan(&exists)
    } else {
       err = r.db.QueryRow(checkQuery, barangID).Scan(&exists)
    }
    if err != nil {
        return err
    }

    if exists {
        updateQuery := "UPDATE mstok SET stok_akhir = stok_akhir + $1, updated_at = CURRENT_TIMESTAMP WHERE barang_id = $2"
        if tx != nil {
            _, err = tx.Exec(updateQuery, qtyChange, barangID)
        } else {
            _, err = r.db.Exec(updateQuery, qtyChange, barangID)
        }
    } else {
        insertQuery := "INSERT INTO mstok (barang_id, stok_akhir) VALUES ($1, $2)"
        if tx != nil {
             _, err = tx.Exec(insertQuery, barangID, qtyChange)
        } else {
             _, err = r.db.Exec(insertQuery, barangID, qtyChange)
        }
    }
    return err
}

func (r *stokRepository) GetHistory(barangID int) ([]models.HistoryStok, error) {
    query := `
        SELECT h.id, h.barang_id, h.user_id, h.jenis_transaksi, h.jumlah, h.stok_sebelum, h.stok_sesudah, h.keterangan, h.created_at,
               b.nama_barang, u.username
        FROM history_stok h
        JOIN master_barang b ON h.barang_id = b.id
        JOIN users u ON h.user_id = u.id
    `
    var args []interface{}
    if barangID != 0 {
        query += " WHERE h.barang_id = $1"
        args = append(args, barangID)
    }
    query += " ORDER BY h.created_at DESC"

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var history []models.HistoryStok
    for rows.Next() {
        var h models.HistoryStok
        h.Barang = &models.Barang{}
        h.User = &models.User{}
        if err := rows.Scan(&h.ID, &h.BarangID, &h.UserID, &h.JenisTransaksi, &h.Jumlah, &h.StokSebelum, &h.StokSesudah, &h.Keterangan, &h.CreatedAt, &h.Barang.NamaBarang, &h.User.Username); err != nil {
            return nil, err
        }
        history = append(history, h)
    }
    return history, nil
}

func (r *stokRepository) CreateHistory(tx *sql.Tx, h *models.HistoryStok) error {
    query := `INSERT INTO history_stok (barang_id, user_id, jenis_transaksi, jumlah, stok_sebelum, stok_sesudah, keterangan) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
    
    var err error
    if tx != nil {
        _, err = tx.Exec(query, h.BarangID, h.UserID, h.JenisTransaksi, h.Jumlah, h.StokSebelum, h.StokSesudah, h.Keterangan)
    } else {
        _, err = r.db.Exec(query, h.BarangID, h.UserID, h.JenisTransaksi, h.Jumlah, h.StokSebelum, h.StokSesudah, h.Keterangan)
    }
    return err
}
