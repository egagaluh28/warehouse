package repositories

import (
	"database/sql"
	"fmt"
	"warehouse-api/models"
)

type BarangRepository interface {
	Create(barang *models.Barang) error
	Update(barang *models.Barang) error
	Delete(id int) error
	GetByID(id int) (*models.BarangWithStok, error)
	GetAll(search string, limit, offset int, sortBy, order string) ([]models.Barang, int, error) // Returns data, total count, error
	GetAllWithStok(search string, limit, offset int, sortBy, order string) ([]models.BarangWithStok, int, error)
    Exists(id int) (bool, error)
}

type barangRepository struct {
	db *sql.DB
}

func NewBarangRepository(db *sql.DB) BarangRepository {
	return &barangRepository{db}
}

func (r *barangRepository) Create(barang *models.Barang) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var nextID int
	err = tx.QueryRow("SELECT nextval(pg_get_serial_sequence('master_barang','id'))").Scan(&nextID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	kode := fmt.Sprintf("BRG-%03d", nextID)

	query := `INSERT INTO master_barang (id, kode_barang, nama_barang, deskripsi, satuan, harga_beli, harga_jual)
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(query, nextID, kode, barang.NamaBarang, barang.Deskripsi, barang.Satuan, barang.HargaBeli, barang.HargaJual)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	barang.ID = nextID
	barang.KodeBarang = kode
	return nil
}

func (r *barangRepository) Update(barang *models.Barang) error {
	query := `UPDATE master_barang SET kode_barang=$1, nama_barang=$2, deskripsi=$3, satuan=$4, harga_beli=$5, harga_jual=$6 WHERE id=$7`
	_, err := r.db.Exec(query, barang.KodeBarang, barang.NamaBarang, barang.Deskripsi, barang.Satuan, barang.HargaBeli, barang.HargaJual, barang.ID)
	return err
}

func (r *barangRepository) Delete(id int) error {
	query := `DELETE FROM master_barang WHERE id=$1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *barangRepository) GetByID(id int) (*models.BarangWithStok, error) {
	query := `
        SELECT b.id, b.kode_barang, b.nama_barang, b.deskripsi, b.satuan, b.harga_beli, b.harga_jual, COALESCE(s.stok_akhir, 0)
        FROM master_barang b
        LEFT JOIN mstok s ON b.id = s.barang_id
        WHERE b.id = $1`
	var barang models.BarangWithStok
	err := r.db.QueryRow(query, id).Scan(
		&barang.ID, &barang.KodeBarang, &barang.NamaBarang, &barang.Deskripsi, &barang.Satuan, &barang.HargaBeli, &barang.HargaJual, &barang.Stok,
	)
	if err != nil {
		return nil, err
	}
	return &barang, nil
}

func (r *barangRepository) GetAll(search string, limit, offset int, sortBy, order string) ([]models.Barang, int, error) {
	var whereClause string
	var args []interface{}
	idx := 1

	if search != "" {
		whereClause = "WHERE kode_barang ILIKE $1 OR nama_barang ILIKE $1"
		args = append(args, "%"+search+"%")
		idx++
	}

    // Default Sorting
    orderByClause := "ORDER BY id ASC"
    if sortBy != "" {
        // Whitelist allowed columns to prevent SQL Injection
        allowedSorts := map[string]string{
            "harga_beli": "harga_beli",
            "harga_jual": "harga_jual",
            "nama":       "nama_barang",
            "id":         "id",
        }
        
        if col, ok := allowedSorts[sortBy]; ok {
            ord := "ASC"
            if order == "desc" || order == "DESC" {
                ord = "DESC"
            }
            orderByClause = fmt.Sprintf("ORDER BY %s %s", col, ord)
        }
    }

	// Get Total Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM master_barang %s", whereClause)
	var total int
	// Re-construct args for count query (without limit/offset)
    countArgs := make([]interface{}, len(args))
    copy(countArgs, args)
    
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get Data
	query := fmt.Sprintf("SELECT id, kode_barang, nama_barang, deskripsi, satuan, harga_beli, harga_jual FROM master_barang %s %s LIMIT $%d OFFSET $%d", whereClause, orderByClause, idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var barangs []models.Barang
	for rows.Next() {
		var b models.Barang
		if err := rows.Scan(&b.ID, &b.KodeBarang, &b.NamaBarang, &b.Deskripsi, &b.Satuan, &b.HargaBeli, &b.HargaJual); err != nil {
			return nil, 0, err
		}
		barangs = append(barangs, b)
	}

	return barangs, total, nil
}

func (r *barangRepository) GetAllWithStok(search string, limit, offset int, sortBy, order string) ([]models.BarangWithStok, int, error) {
	var whereClause string
	var args []interface{}
	idx := 1

	if search != "" {
		whereClause = "WHERE b.kode_barang ILIKE $1 OR b.nama_barang ILIKE $1"
		args = append(args, "%"+search+"%")
		idx++
	}

	orderByClause := "ORDER BY b.id ASC"
	if sortBy != "" {
		allowedSorts := map[string]string{
			"harga_beli": "b.harga_beli",
			"harga_jual": "b.harga_jual",
			"nama":       "b.nama_barang",
			"id":         "b.id",
			"stok":       "COALESCE(s.stok_akhir, 0)",
		}

		if col, ok := allowedSorts[sortBy]; ok {
			ord := "ASC"
			if order == "desc" || order == "DESC" {
				ord = "DESC"
			}
			orderByClause = fmt.Sprintf("ORDER BY %s %s", col, ord)
		}
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM master_barang b %s", whereClause)
	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)

	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT b.id, b.kode_barang, b.nama_barang, b.deskripsi, b.satuan, b.harga_beli, b.harga_jual, COALESCE(s.stok_akhir, 0)
		FROM master_barang b
		LEFT JOIN mstok s ON b.id = s.barang_id
		%s
		%s
		LIMIT $%d OFFSET $%d`, whereClause, orderByClause, idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var barangs []models.BarangWithStok
	for rows.Next() {
		var b models.BarangWithStok
		if err := rows.Scan(&b.ID, &b.KodeBarang, &b.NamaBarang, &b.Deskripsi, &b.Satuan, &b.HargaBeli, &b.HargaJual, &b.Stok); err != nil {
			return nil, 0, err
		}
		barangs = append(barangs, b)
	}

	return barangs, total, nil
}

func (r *barangRepository) Exists(id int) (bool, error) {
    var exists bool
    query := "SELECT EXISTS(SELECT 1 FROM master_barang WHERE id=$1)"
    err := r.db.QueryRow(query, id).Scan(&exists)
    return exists, err
}
