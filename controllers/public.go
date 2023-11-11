package controllers

import (
	"games/auth"
	"games/database"
	"games/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var redisClient *redis.Client

func init() {
	

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",               
		DB:       0,                
	})

}


type LoginPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}


type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}



func RegisterPlayer(c *gin.Context) {
	var player models.Player
	err := c.ShouldBindJSON(&player)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs ",
		})
		c.Abort()
		return
	}
	err = player.HashPassword(player.Password)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{
			"Error": "Error Hashing Password",
		})
		c.Abort()
		return
	}
	err = player.CreatePlayerRecord()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Creating Player",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Message": "Sucessfully Register",
	})
}



func Login(c *gin.Context) {
	var payload LoginPayload
	var player models.Player
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs",
		})
		c.Abort()
		return
	}
	result := database.GlobalDB.Where("email = ?", payload.Email).First(&player)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(401, gin.H{
			"Error": "Invalid Player Credentials",
		})
		c.Abort()
		return
	}
	err = player.CheckPassword(payload.Password)
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{
			"Error": "Invalid Player Credentials",
		})
		c.Abort()
		return
	}
	jwtWrapper := auth.JwtWrapper{
		SecretKey:         "rikkikey",
		Issuer:            "AuthService",
		ExpirationMinutes: 15,
		ExpirationHours:   12,
	}
	signedToken, err := jwtWrapper.GenerateToken(player.Email, player.ID)
	log.Println("Signed Token:", signedToken)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}

	
	setErr := redisClient.Set(c, player.Email+"_access_token", signedToken, time.Duration(jwtWrapper.ExpirationMinutes)*time.Minute)
	log.Println("Set Error Access Token:", setErr.Err())
	if setErr.Err() != nil {
		log.Println(setErr.Err())
		c.JSON(500, gin.H{
			"Error": "Gagal menyimpan token akses di Redis",
		})
		c.Abort()
		return
	}

	signedtoken, err := jwtWrapper.RefreshToken(player.Email,player.ID)
	log.Println("Signed Refresh Token:", signedtoken)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Signing Token",
		})
		c.Abort()
		return
	}

	
	setErrRefresh := redisClient.Set(c, player.Email+"_refresh_token", signedtoken, time.Duration(jwtWrapper.ExpirationHours)*time.Hour)
	log.Println("Set Error Refresh Token:", setErrRefresh.Err())
	if setErrRefresh.Err() != nil {
		log.Println(setErrRefresh.Err())
		c.JSON(500, gin.H{
			"Error": "Gagal menyimpan token pembaharuan di Redis",
		})
		c.Abort()
		return
	}

	tokenResponse := LoginResponse{
		Token:        signedToken,
		RefreshToken: signedtoken,
	}
	c.JSON(200, tokenResponse)
}



func Logout(c *gin.Context) {

	
	email, ok := c.Get("email")
	if !ok || email == nil {
		log.Println("Error getting email from context")
		c.JSON(500, gin.H{"Error": "Gagal mendapatkan email dari konteks"})
		c.Abort()
		return
	}
	
	log.Println("Email from context:", email)

	
	err := redisClient.Del(c, email.(string)+"_access_token").Err()
	if err != nil {
		
		
	
		c.JSON(500, gin.H{"Error": "Gagal menghapus token akses dari Redis"})
		c.Abort()
		return
	}
	
	err = redisClient.Del(c, email.(string)+"_refresh_token").Err()
	if err != nil {
		
		log.Println("Error deleting refresh token from Redis:", err)
		c.JSON(500, gin.H{"Error": "Gagal menghapus token pembaharuan dari Redis"})
		c.Abort()
		return
	}

	
	log.Println("Logout berhasil")

	c.JSON(200, gin.H{
		"Message": "Berhasil logout",
	})
}
