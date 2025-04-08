package storage

import (
	"fmt"
	"os"
	"time"
	"gorm.io/gorm"

	"inventory/src/types"

	"github.com/jung-kurt/gofpdf"
)

func generatePDF(db *gorm.DB, duration time.Duration) error {
	// Ambil data yang lebih dari 3 bulan
	var sebaranBarangList []types.Inventaris
	cutOffTime := time.Now().Add(-duration)

	// Ambil data yang lebih lama dari 3 bulan
	result := db.Where("created_at < ?", cutOffTime).Find(&sebaranBarangList)
	if result.Error != nil {
		return fmt.Errorf("error retrieving data: %v", result.Error)
	}

	// Mengelompokkan data berdasarkan tahun, bulan, dan hari
	for _, sebaran := range sebaranBarangList {
		// Mendapatkan tahun, bulan, dan hari dari `created_at`
		year, month, day := sebaran.CreatedAt.Date()

		// Menyusun path folder berdasarkan tahun/bulan/hari
		folderPath := fmt.Sprintf("./exports/%d/%02d/%02d", year, month, day)

		// Pastikan folder tujuan ada
		err := os.MkdirAll(folderPath, os.ModePerm) // Membuat folder ./exports/tahun/bulan/hari jika belum ada
		if err != nil {
			return fmt.Errorf("error creating folder: %v", err)
		}

		// Membuat file PDF di folder berdasarkan tahun/bulan/hari
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()

		// Menambahkan konten PDF
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("ID: %d", sebaran.ID))
		pdf.Ln(10) // Line break
		pdf.Cell(40, 10, fmt.Sprintf("CreatedAt: %s", sebaran.CreatedAt))
		pdf.Ln(10)

		// Menyimpan PDF ke folder berdasarkan tahun/bulan/hari
		err = pdf.OutputFileAndClose(fmt.Sprintf("%s/data_sebaran_%d.pdf", folderPath, sebaran.ID))
		if err != nil {
			return fmt.Errorf("error generating PDF: %v", err)
		}
	}

	// Hapus data dari database
	err := db.Delete(&types.SebaranBarang{}, "created_at < ?", cutOffTime).Error
	if err != nil {
		return fmt.Errorf("error deleting data from database: %v", err)
	}

	fmt.Println("PDF telah berhasil dibuat dan data telah dihapus")
	return nil
}
