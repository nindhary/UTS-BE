package models

import "time"

type PekerjaanAlumni struct {
	ID                  int       `json:"id"`
	AlumniID            int       `json:"alumni_id"`
	UserID              int       `json:"user_id"`
	NamaPerusahaan      string    `json:"nama_perusahaan"`
	PosisiJabatan       string    `json:"posisi_jabatan"`
	BidangIndustri      string    `json:"bidang_industri"`
	LokasiKerja         string    `json:"lokasi_kerja"`
	GajiRange           string    `json:"gaji_range"`
	TanggalMulaiKerja   string    `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja string    `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string    `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string    `json:"deskripsi_pekerjaan"`
	IsDeleted           bool      `json:"isdeleted"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type CreatePekerjaanRequest struct {
	AlumniID            int    `json:"alumni_id"`
	UserID              int    `json:"user_id"`
	NamaPerusahaan      string `json:"nama_perusahaan"`
	PosisiJabatan       string `json:"posisi_jabatan"`
	BidangIndustri      string `json:"bidang_industri"`
	LokasiKerja         string `json:"lokasi_kerja"`
	GajiRange           string `json:"gaji_range"`
	TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanRequest struct {
	NamaPerusahaan      string `json:"nama_perusahaan"`
	UserID              int    `json:"user_id"`
	PosisiJabatan       string `json:"posisi_jabatan"`
	BidangIndustri      string `json:"bidang_industri"`
	LokasiKerja         string `json:"lokasi_kerja"`
	GajiRange           string `json:"gaji_range"`
	TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
	IsDeleted           bool   `json:"isdeleted"`
}
