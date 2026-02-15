package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5" // Descomentar quando implementar validação real JWT
)

// AuthMiddleware valida o token JWT do Supabase
// NOTA: Em produção, você deve validar a assinatura do JWT usando a chave secreta do Supabase.
// Para este exemplo, estamos apenas verificando a presença do token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticação não fornecido"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato do token inválido"})
			return
		}

		tokenString := parts[1]

		// TODO: Validar o token com a biblioteca JWT e a chave secreta do Supabase
		// Por enquanto, assumimos que se o token existe, é válido (MOCK para desenvolvimento inicial)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			return
		}

		// Extrair UserID do token (simulado)
		// Em produção: claims := token.Claims.(jwt.MapClaims); userID := claims["sub"]
		// c.Set("userID", userID)

		c.Next()
	}
}
