package penyusutan

import (
	"errors"
	"net/http"
	"strconv"
	"inventory/src/types"


	"github.com/gin-gonic/gin"
)

func HitungPenyusutanHandler(c *gin.Context) {
	hargaBeliStr := c.Param("hargaBeli")
	hargaBeli, err := strconv.Atoi(hargaBeliStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harga beli harus berupa angka"})
		return
	}

	result, err := HitungPenyusutan(hargaBeli)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Fungsi perhitungan penyusutan
func HitungPenyusutan(hargaBeli int) ([]types.Penyusutan, error) {
	if hargaBeli <= 0 {
		return nil, errors.New("harga beli harus lebih dari 0")
	}

	var penyusutanList []types.Penyusutan
	penyusutanPersen := 2.5
	penyusutanNilai := hargaBeli * int(penyusutanPersen) / 100
	hargaSetelahPenyusutan := hargaBeli

	for hargaSetelahPenyusutan > 0 {
		hargaSetelahPenyusutan -= penyusutanNilai
		if hargaSetelahPenyusutan < 0 {
			hargaSetelahPenyusutan = 0
		}

		penyusutanList = append(penyusutanList, types.Penyusutan{
			NilaiPenyusutan: penyusutanNilai,
			HargaPenyusutan: hargaSetelahPenyusutan,
		})
	}

	return penyusutanList, nil
}
