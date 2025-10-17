package repository

import (
	models "crud-app/app/model"
	"crud-app/database"
	"database/sql"
	"fmt"
)

type PekerjaanRepository interface {
	GetAll() ([]models.PekerjaanAlumni, error)
	GetByID(id int) (models.PekerjaanAlumni, error)
	Create(input models.CreatePekerjaanRequest) (models.PekerjaanAlumni, error)
	Update(id int, input models.UpdatePekerjaanRequest) (models.PekerjaanAlumni, error)
	Delete(id int) error
	GetByAlumni(alumniID int) ([]models.PekerjaanAlumni, error)
	GetPekerjaanRepo(search, sortBy, order string, limit, offset int) ([]models.PekerjaanAlumni, error)
	CountPekerjaanRepo(search string) (int, error)
	SoftDelete(id int, req models.UpdatePekerjaanRequest) error
	GetTrash() ([]models.PekerjaanAlumni, error)
	Restore(id int) error
	HardDelete(id int) error
}

type pekerjaanRepository struct{}

func NewPekerjaanRepository() PekerjaanRepository {
	return &pekerjaanRepository{}
}

// ===== Ambil semua data =====
func (r *pekerjaanRepository) GetAll() ([]models.PekerjaanAlumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
			&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// ===== Ambil by ID + user pemilik =====
func (r *pekerjaanRepository) GetByID(id int) (models.PekerjaanAlumni, error) {
	var p models.PekerjaanAlumni
	err := database.DB.QueryRow(`
		SELECT p.id,
		       p.alumni_id,
		       a.user_id,
		       p.nama_perusahaan,
		       p.posisi_jabatan,
		       p.bidang_industri,
		       p.lokasi_kerja,
		       p.gaji_range,
		       p.tanggal_mulai_kerja,
		       p.tanggal_selesai_kerja,
		       p.status_pekerjaan,
		       p.deskripsi_pekerjaan,
		       p.created_at,
		       p.updated_at
		FROM pekerjaan_alumni p
		JOIN alumni a ON p.alumni_id = a.id
		WHERE p.id = $1
	`, id).Scan(
		&p.ID,
		&p.AlumniID,
		&p.UserID, // penting: harus ada di struct
		&p.NamaPerusahaan,
		&p.PosisiJabatan,
		&p.BidangIndustri,
		&p.LokasiKerja,
		&p.GajiRange,
		&p.TanggalMulaiKerja,
		&p.TanggalSelesaiKerja,
		&p.StatusPekerjaan,
		&p.DeskripsiPekerjaan,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	return p, err
}

// ===== Create =====
func (r *pekerjaanRepository) Create(input models.CreatePekerjaanRequest) (models.PekerjaanAlumni, error) {
	var p models.PekerjaanAlumni
	err := database.DB.QueryRow(`
		INSERT INTO pekerjaan_alumni
		(alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
		 gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
		 deskripsi_pekerjaan, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,NOW(),NOW())
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		          lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		          status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at`,
		input.AlumniID, input.NamaPerusahaan, input.PosisiJabatan, input.BidangIndustri,
		input.LokasiKerja, input.GajiRange, input.TanggalMulaiKerja, input.TanggalSelesaiKerja,
		input.StatusPekerjaan, input.DeskripsiPekerjaan,
	).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

// ===== Update =====
func (r *pekerjaanRepository) Update(id int, input models.UpdatePekerjaanRequest) (models.PekerjaanAlumni, error) {
	var p models.PekerjaanAlumni
	err := database.DB.QueryRow(`
		UPDATE pekerjaan_alumni
		SET nama_perusahaan=$1, posisi_jabatan=$2, bidang_industri=$3, lokasi_kerja=$4,
		    gaji_range=$5, tanggal_mulai_kerja=$6, tanggal_selesai_kerja=$7,
		    status_pekerjaan=$8, deskripsi_pekerjaan=$9, updated_at=NOW()
		WHERE id=$10
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		          lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		          status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at`,
		input.NamaPerusahaan, input.PosisiJabatan, input.BidangIndustri, input.LokasiKerja,
		input.GajiRange, input.TanggalMulaiKerja, input.TanggalSelesaiKerja,
		input.StatusPekerjaan, input.DeskripsiPekerjaan, id,
	).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
		&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
		&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return p, fmt.Errorf("pekerjaan alumni dengan id %d tidak ditemukan", id)
	}
	return p, err
}

// ===== Delete permanen =====
func (r *pekerjaanRepository) Delete(id int) error {
	result, err := database.DB.Exec(`DELETE FROM pekerjaan_alumni WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("pekerjaan alumni dengan id %d tidak ditemukan", id)
	}
	return nil
}

// ===== Get by Alumni =====
func (r *pekerjaanRepository) GetByAlumni(alumniID int) ([]models.PekerjaanAlumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE alumni_id = $1`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
			&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, p)
	}
	return data, nil
}

// ===== Pagination + search =====
func (r *pekerjaanRepository) GetPekerjaanRepo(search, sortBy, order string, limit, offset int) ([]models.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`, sortBy, order)

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
			&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

// ===== Count untuk pagination =====
func (r *pekerjaanRepository) CountPekerjaanRepo(search string) (int, error) {
	var total int
	err := database.DB.QueryRow(
		`SELECT COUNT(*) FROM pekerjaan_alumni
		 WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1`,
		"%"+search+"%").Scan(&total)
	return total, err
}

// ===== SoftDelete =====
func (r *pekerjaanRepository) SoftDelete(id int, _ models.UpdatePekerjaanRequest) error {
	_, err := database.DB.Exec(`
		UPDATE pekerjaan_alumni
		SET isdeleted = TRUE,
		    updated_at = NOW()
		WHERE id = $1`, id)
	return err
}

// get soft deleted data
func (r *pekerjaanRepository) GetTrash() ([]models.PekerjaanAlumni, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE isdeleted = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.PekerjaanAlumni
	for rows.Next() {
		var p models.PekerjaanAlumni
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange,
			&p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan,
			&p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

// restore soft deleted data
func (r *pekerjaanRepository) Restore(id int) error {
	query := "UPDATE pekerjaan_alumni SET isdeleted = TRUE WHERE id = $1 AND isdeleted = NULL"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("data not found or not soft deleted")
	}

	return nil
}

// hard delete
func (r *pekerjaanRepository) HardDelete(id int) error {
	query := "DELETE FROM pekerjaan_alumni WHERE id = $1 AND isdeleted = TRUE"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("data not found or not soft deleted")
	}

	return nil
}
