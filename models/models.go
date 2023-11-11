package models

import (
	"games/database"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	ID       int    `gorm:"primaryKey"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required" gorm:"unique"`
	Password string `json:"password" binding:"required"`
	Wallet   uint32 `json:"wallet" gorm:"default:0"`
}

type ProfileResponse struct {
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	Password         string    `json:"password"`
	Wallet           uint32    `json:"wallet"`
	AccountName      string    `json:"accountName"`
	AccountNumber    string    `json:"accountNumber"`
	BankName         string    `json:"bankName"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	AccountCreatedAt time.Time `json:"accountCreatedAt" `
}

type PlayerAccountResponse struct {
	Username         string    `json:"username"`
	Wallet           uint32    `json:"wallet"`
	AccountName      string    `json:"accountName"`
	AccountNumber    string    `json:"accountNumber"`
	BankName         string    `json:"bankName"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	AccountCreatedAt time.Time `json:"accountCreatedAt" `
}

type Balance struct {
	Balance uint32 `json:"balance" binding:"required"`
}
type Account struct {
	gorm.Model
	ID            int    `gorm:"primaryKey"`
	PlayerID      int    `json:"playerId"`
	AccountName   string `json:"accountName" binding:"required" gorm:"unique"`
	AccountNumber string `json:"accountNumber" binding:"required" gorm:"unique"`
	BankName      string `json:"bankName" binding:"required"`
}

func (account *Account) CreateAccountRecord() error {
	result := database.GlobalDB.Create(&account)
	log.Print("rikki")
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (player *Player) CreatePlayerRecord() error {
	result := database.GlobalDB.Create(&player)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (player *Player) UpdatePlayerRecord() error {
	result := database.GlobalDB.Save(&player)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (player *Player) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	player.Password = string(bytes)
	return nil
}

func (player *Player) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
