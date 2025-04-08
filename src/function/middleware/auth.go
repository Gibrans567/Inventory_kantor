package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"


	"inventory/src/types"
	"inventory/src/function/utils"
	"inventory/src/function/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// TokenMetadata adalah struktur untuk menyimpan informasi dari token
type TokenMetadata struct {
	UserID uint
	Role   string
	Exp    int64
}

// AuthMiddleware adalah middleware untuk memverifikasi token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan token dari header Authorization
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Token tidak ditemukan"})
			c.Abort()
			return
		}

		// Memverifikasi token
		tokenData, err := verifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Unauthorized: %v", err)})
			c.Abort()
			return
		}

		// Memeriksa apakah token sudah kadaluarsa
		if time.Now().Unix() > tokenData.Exp {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Token kadaluarsa"})
			c.Abort()
			return
		}

		// Menyimpan data token ke context untuk digunakan oleh handler berikutnya
		c.Set("user_id", tokenData.UserID)
		c.Set("role", tokenData.Role)
		
		c.Next()
	}
}

// RoleAuthMiddleware adalah middleware untuk memeriksa peran pengguna
func RoleAuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan peran dari context yang telah diset oleh AuthMiddleware
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Autentikasi diperlukan"})
			c.Abort()
			return
		}

		// Memeriksa apakah peran pengguna termasuk dalam peran yang diizinkan
		userRole := role.(string)
		isAllowed := false
		for _, r := range allowedRoles {
			if r == userRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Anda tidak memiliki akses"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractToken mengambil token dari header Authorization
func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// verifyToken memverifikasi token dan mengekstrak data dari token
func verifyToken(tokenString string) (*TokenMetadata, error) {
	// Mendapatkan secret key dari environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "default_secret_key" // Default key jika tidak diset di environment
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Memeriksa metode signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode signing tidak dikenali: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Memeriksa validitas token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Mengekstrak data dari claims
		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)
		exp := int64(claims["exp"].(float64))

		// Membuat struktur TokenMetadata
		tokenData := &TokenMetadata{
			UserID: userID,
			Role:   role,
			Exp:    exp,
		}

		return tokenData, nil
	}

	return nil, fmt.Errorf("token tidak valid")
}

	// GenerateToken membuat token JWT baru
	func GenerateToken(userID uint, role string) (string, error) {
		// Mendapatkan secret key dari environment variable
		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			secretKey = "default_secret_key" // Default key jika tidak diset di environment
		}
	
		// Membuat claims
		claims := jwt.MapClaims{
			"user_id": userID,
			"role":    role,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku selama 24 jam
		}
	
		// Membuat token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
		// Menandatangani token dengan secret key
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return "", err
		}
	
		return tokenString, nil
	}

	func LoginHandler(c *gin.Context) {
		var loginForm struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
	
		if err := c.ShouldBindJSON(&loginForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		db := database.GetDB()
		// Periksa kredensial pengguna dari database
		var user types.User
		if err := db.Where("email = ?", loginForm.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
	
		// Verifikasi password - pastikan fungsi ini didefinisikan di suatu tempat
		// atau gunakan bcrypt langsung di sini
		if !utils.CheckPasswordHash(loginForm.Password, user.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
	
		// Generate token JWT - pastikan GenerateToken diimpor dengan benar
		token, err := GenerateToken(user.ID, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"token":   token,
			"user": gin.H{
				"id":        user.ID,
				"name":      user.NamaUser,  // Menggunakan NamaUser dari struct
				"email":     user.Email,
				"role":      user.Role,
				"id_divisi": user.IdDivisi,
			},
		})
	}