package controllers

import (
	"games/database"
	"games/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Profile(c *gin.Context) {

	var player models.Player

	email, _ := c.Get("email")

	result := database.GlobalDB.Where("email = ?", email.(string)).First(&player)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "Player Not Found",
		})
		c.Abort()
		return
	}

	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get Player Profile",
		})
		c.Abort()
		return
	}

	var account models.Account

	resultAccount := database.GlobalDB.Where("player_id = ?", player.ID).First(&account)

	if resultAccount.Error != nil {
		c.JSON(500, gin.H{
			"Error": "Could Not Get Player Account",
		})
		c.Abort()
		return
	}

	response := models.ProfileResponse{
		Username:      player.Username,
		Email:         player.Email,
		Password:      player.Password,
		Wallet:        player.Wallet,
		AccountName:   account.AccountName,
		AccountNumber: account.AccountNumber,
		BankName:      account.BankName,
	}

	c.JSON(200, response)
}

func Account(c *gin.Context) {
	var account models.Account
	err := c.ShouldBindJSON(&account)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs ",
		})
		c.Abort()
		return
	}

	playerId, _ := c.Get("playerId")
	account.PlayerID = playerId.(int)
	err = account.CreateAccountRecord()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Creating Account",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Message": "Sucessfully Register Account",
	})
}

func TopUpBalance(c *gin.Context) {

	var player models.Player

	playerId, _ := c.Get("playerId")

	result := database.GlobalDB.Where("ID = ?", playerId.(int)).First(&player)

	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "Player Not Found",
		})
		c.Abort()
		return
	}

	var balance models.Balance
	err := c.ShouldBindJSON(&balance)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"Error": "Invalid Inputs ",
		})
		c.Abort()
		return
	}

	player.Wallet += balance.Balance
	player.UpdatedAt = time.Now().Local()

	err = player.UpdatePlayerRecord()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Error": "Error Updating Balance",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"Message": "Sucessfully Update Balance",
	})
}

func GetListPlayer(c *gin.Context) {
	log.Println("error :")

	username := c.Query("username")
	accountName := c.Query("accountName")
	accountNumber := c.Query("accountNumber")
	bankName := c.Query("bankName")
	wallet := c.Query("wallet")

	query := database.GlobalDB.Debug().Table("players").Select("players.*, accounts.account_name, accounts.account_number, accounts.bank_name,accounts.created_at as account_created_at ")

	query = query.Joins("LEFT JOIN accounts ON players.ID = accounts.player_id")

	if username != "" {
		query = query.Where("players.username like ?", "%"+username+"%")
	}
	if accountName != "" {
		query = query.Where("accounts.account_name like ?", accountName)
	}
	if accountNumber != "" {
		query = query.Where("accounts.account_number like ?", accountNumber)
	}
	if bankName != "" {
		query = query.Where("accounts.bank_name like ?", bankName)
	}
	if wallet != "" {
		query = query.Where("players.wallet <= ?", wallet)
	}

	var playerResponses []models.PlayerAccountResponse

	if err := query.Scan(&playerResponses).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if playerResponses != nil {
		c.JSON(200, playerResponses)
	} else {
		c.JSON(200, "Data tidak ditemukan")

	}

}
