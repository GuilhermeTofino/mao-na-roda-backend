package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware valida o token JWT do Supabase usando a chave secreta
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
		jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar o algoritmo de assinatura
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			return
		}

		// Extrair Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Erro ao ler claims do token"})
			return
		}

		// Adicionar ID do usuário ao contexto (Supabase usa 'sub' como ID)
		if sub, ok := claims["sub"].(string); ok {
			// Em alguns casos pode ser necessário simular que esse ID veio de Header (como fizemos no mock)
			// Mas idealmente o handler leria do contexto.
			// Vamos setar um Header interno para compatibilidade com o código atual que lê X-User-ID
			c.Request.Header.Set("X-User-ID", sub)
			c.Set("userID", sub)
		}

		c.Next()
	}
}
