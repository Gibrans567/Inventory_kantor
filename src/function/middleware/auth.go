package middleware

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"net/http"
	"strings"

	"time"
	"github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
    var req struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db := database.GetDB()
    var user types.User
    if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := GenerateJWT(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token, "message": "Login successful","user": gin.H{
        "id":       user.ID,
        "role":     user.Role,  // Assuming user.Role exists in your User struct
        "division_id": user.IdDivisi,  // Add division_id or tenant_id if applicable
    },})
}

func GenerateJWT(user types.User) (string, error) {
    claims := types.JWTClaims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token berlaku selama 1 hari
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte("7b7a4159f6d56da06e2033d33c411b6482743cde")) // Ganti dengan secret key yang lebih aman
}

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte("7b7a4159f6d56da06e2033d33c411b6482743cde"), nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        claims, ok := token.Claims.(*types.JWTClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            return
        }

        // Inject user info ke context
        c.Set("userID", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists || role != requiredRole {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
            return
        }
        c.Next()
    }
}

func RegisterHandler(c *gin.Context) {
    var req types.User

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    req.Password = string(hashedPassword)

    db := database.GetDB()
    if err := db.Create(&req).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout success, silakan hapus token di client",
	})
}


