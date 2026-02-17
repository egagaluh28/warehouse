package integration

import (
	"fmt"
	"testing"

	"warehouse-api/models"
	"warehouse-api/repositories"

	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryIntegration(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	testDB.Exec("TRUNCATE users CASCADE")

	repo := repositories.NewUserRepository(testDB)

	user := &models.User{
		Username: "integrationuser",
		Password: "hashedpassword",
		Email:    "integration@test.com",
		FullName: "Integration User",
		Role:     "staff",
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.Greater(t, user.ID, 0)

	retrieved, err := repo.GetByUsername("integrationuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, retrieved.Username)
}

func TestBarangRepositoryIntegration(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	testDB.Exec("TRUNCATE master_barang CASCADE")

	repo := repositories.NewBarangRepository(testDB)

	barang := &models.Barang{
		NamaBarang: "Integration Test Item",
		Satuan:     "unit",
		HargaBeli:  5000,
		HargaJual:  7000,
	}

	err := repo.Create(barang)
	assert.NoError(t, err)
	assert.Greater(t, barang.ID, 0)

	retrieved, err := repo.GetByID(barang.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, retrieved.KodeBarang)

	err = repo.Delete(barang.ID)
	assert.NoError(t, err)
}

func TestSearchPaginationIntegration(t *testing.T) {
	if testDB == nil {
		t.Skip("Database not available")
	}

	testDB.Exec("TRUNCATE master_barang CASCADE")

	repo := repositories.NewBarangRepository(testDB)

	for i := 1; i <= 15; i++ {
		barang := &models.Barang{
			NamaBarang: fmt.Sprintf("Search Item %d", i),
			Satuan:     "pcs",
			HargaBeli:  1000 * float64(i),
			HargaJual:  1500 * float64(i),
		}
		repo.Create(barang)
	}

	items, total, err := repo.GetAll("", 10, 0, "", "")
	assert.NoError(t, err)
	assert.Equal(t, 10, len(items))
	assert.Equal(t, 15, total)
}
