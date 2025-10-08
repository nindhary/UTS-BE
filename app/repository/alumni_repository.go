package repository

import (
	models "crud-app/app/model"
	"crud-app/database"
	"fmt"
)

type AlumniRepository interface {
	GetAll() ([]models.Alumni, error)
	GetByID(id int) (models.Alumni, error)
	Create(input models.CreateAlumniRequest) (models.Alumni, error)
	Update(id int, input models.UpdateAlumniRequest) error
	Delete(id int) error
	GetAlumniRepo(search, sortBy, order string, limit, offset int) ([]models.Alumni, error)
	CountAlumniRepo(search string) (int, error)
}

type alumniRepository struct{}

func NewAlumniRepository() AlumniRepository {
	return &alumniRepository{}
}

// Ambil semua data alumni
func (r *alumniRepository) GetAll() ([]models.Alumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		       no_telepon, alamat, created_at, updated_at 
		FROM alumni 
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumni []models.Alumni
	for rows.Next() {
		var a models.Alumni
		if err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}
	return alumni, nil
}

// Ambil alumni berdasarkan ID
func (r *alumniRepository) GetByID(id int) (models.Alumni, error) {
	var a models.Alumni
	err := database.DB.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		       no_telepon, alamat, created_at, updated_at 
		FROM alumni 
		WHERE id = $1`, id).
		Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt,
		)
	if err != nil {
		return a, err
	}
	return a, nil
}

// Tambah data alumni baru
func (r *alumniRepository) Create(input models.CreateAlumniRequest) (models.Alumni, error) {
	var a models.Alumni
	err := database.DB.QueryRow(`
		INSERT INTO alumni 
		    (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW()) 
		RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		          no_telepon, alamat, created_at, updated_at`,
		input.NIM, input.Nama, input.Jurusan, input.Angkatan,
		input.TahunLulus, input.Email, input.NoTelepon, input.Alamat,
	).Scan(
		&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
		&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
		&a.CreatedAt, &a.UpdatedAt,
	)
	return a, err
}

// Update data alumni (tanpa ubah NIM)
func (r *alumniRepository) Update(id int, input models.UpdateAlumniRequest) error {
	_, err := database.DB.Exec(`
		UPDATE alumni 
		SET nama=$1, jurusan=$2, angkatan=$3, tahun_lulus=$4, 
		    email=$5, no_telepon=$6, alamat=$7, updated_at=NOW() 
		WHERE id=$8`,
		input.Nama, input.Jurusan, input.Angkatan,
		input.TahunLulus, input.Email, input.NoTelepon, input.Alamat, id,
	)
	return err
}

// Hapus data alumni
func (r *alumniRepository) Delete(id int) error {
	_, err := database.DB.Exec(`DELETE FROM alumni WHERE id=$1`, id)
	return err
}

// Ambil alumni dengan pagination, search, dan sorting
func (r *alumniRepository) GetAlumniRepo(search, sortBy, order string, limit, offset int) ([]models.Alumni, error) {
	var alumni []models.Alumni

	// Pastikan sortBy valid
	switch sortBy {
	case "id", "nim", "nama", "jurusan", "angkatan", "tahun_lulus", "email", "created_at":
	default:
		sortBy = "id"
	}

	// Pastikan order valid
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Jika search kosong, ambil semua
	searchPattern := "%"
	if search != "" {
		searchPattern = "%" + search + "%"
	}

	query := fmt.Sprintf(`
        SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email,
               no_telepon, alamat, created_at, updated_at
        FROM alumni
        WHERE nama ILIKE $1 OR email ILIKE $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

	rows, err := database.DB.Query(query, searchPattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Alumni
		if err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}

	return alumni, nil
}

func (r *alumniRepository) CountAlumniRepo(search string) (int, error) {
	var total int
	searchPattern := "%"
	if search != "" {
		searchPattern = "%" + search + "%"
	}
	query := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR email ILIKE $1`
	err := database.DB.QueryRow(query, searchPattern).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
