package repository

import (
	models "crud-app/app/model"
	"database/sql"
)

type NganggurRepository struct {
	DB *sql.DB
}

func (r *NganggurRepository) GetAll() ([]models.Alumni, error) {
	rows, err := r.DB.Query(`SELECT a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.no_telepon, a.alamat FROM alumni a 
			LEFT JOIN pekerjaan_alumni p
			ON a.id = p.alumni_id
			WHERE p.id IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []models.Alumni
	for rows.Next() {
		var a models.Alumni
		if err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}
