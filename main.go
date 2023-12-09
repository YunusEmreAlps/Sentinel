package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"sentinel/config"
	"sentinel/helpers"
	"sentinel/logger"
	"sentinel/mail"
	"sentinel/models"

	_ "github.com/lib/pq"
	"github.com/roylee0704/gron"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// var dbConn *gorm.DB
var isConfigSuccess = false
var equals string = strings.Repeat("=", 50)

// Second: 	gron.Every(1*time.Second)
// Minute: 	gron.Every(1*time.Minute)
// Hour: 		gron.Every(1*time.Hour)
// Day: 		gron.Every(1 * xtime.Day)
// Week: 		gron.Every(1 * xtime.Week)
// gron.Every(30 * xtime.Day).At("00:00")
// gron.Every(1 * xtime.Week).At("23:59")
var repeatTime = gron.Every(5 * time.Minute)

// ignored error messages array
var ignoredErrorMessages = []string{
	// ignored error messages here
}

// To Users
var toUsers = []string{
	// team members here
}

// CC Users
var ccUsers = []string{
	// team members here
}

func main() {
	// Connect to the database
	// dbConn = dbConnection()

	// Run the task repeat time and check the changes
	c := gron.New()
	c.AddFunc(repeatTime, func() {
		// Query the TARGET table and retrieve changes
		changes, err := getChanges()
		if err != nil {
			panic(err)
		}

		// Handle the changes
		fmt.Println(equals)
		if len(changes) > 0 {
			for _, change := range changes {
				logger.CLogger.Infof("INFO: %s:%d - %s", change.Domain, change.Port, change.Message)
			}

			// Filter the changes
			filteredChanges := helpers.FilterChanges(changes, ignoredErrorMessages)

			if len(filteredChanges) > 0 {
				for _, v := range filteredChanges {
					logger.CLogger.Tracef("TRACE: %s:%d - %s", v.Domain, v.Port, v.Message)
				}
				f := helpers.SetChangesToExcel(filteredChanges)
				sendMailWithAttachment(filteredChanges, f)
			}
		} else {
			logger.CLogger.Info("INFO: No changes in the last minute.")
		}
		fmt.Println(equals)
	})
	c.Start()

	// Infinite loop to keep the program running
	select {}
}

// Initialize Application
func init() {
	isConfigSuccess = configureApplication()
	if !isConfigSuccess {
		logger.CLogger.Error("INIT: Application configuration failed. Please check your config file.")
		os.Exit(1)
	}
}

// Configure Application
func configureApplication() bool {
	// Clear the terminal screen
	fmt.Println(equals)
	dir, err := os.Getwd()
	if err != nil {
		logger.CLogger.Error("INIT: Cannot get current working directory os.Getwd()")
		return false
	} else {
		config.ReadConfig(dir)
		logger.CLogger.Info("INIT: Application configuration file read success.")
		return true
	}
}

// DB Connection
func dbConnection() *gorm.DB {
	env := config.C.DB

	// String to Int
	port, err := strconv.Atoi(env.Port)
	if err != nil {
		logger.CLogger.Error("ERROR: ", err)
		os.Exit(1)
	}

	// Connect to the "postgres" database
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", env.Host, port, env.Username, env.Password, env.DBName, env.SSLMode)
	db, err := gorm.Open(postgres.Open(dbInfo), &gorm.Config{})
	if err != nil {
		logger.CLogger.Info("ERROR: ", err)
		os.Exit(1)
	}

	// Connection Success
	logger.CLogger.Success("PostgreSQL Database Connection Success")
	return db
}

func getChanges() ([]models.Log, error) {
	var logs []models.Log

	for _, domain := range helpers.DomainList {
		const maxRetries = 3
		for i := 0; i < maxRetries; i++ {
			isOK, data := helpers.CheckDomainCertificate(domain, 30)
			if isOK && data != nil {
				// Successfully retrieved certificate, break out of the loop
				logs = append(logs, *data)
				break
			} else {
				if data != nil {
					// Certificate will not expired in 30 days
					logger.CLogger.Info("INFO: ", domain+" - "+data.Message)
					break
				} else {
					// Connection Error
					logger.CLogger.Error("ERROR: ", domain+" - Connection Error Attempt: "+strconv.Itoa(i+1)+"/"+strconv.Itoa(maxRetries))
				}
			}
			// Wait for a brief period before retrying
			time.Sleep(2 * time.Second)
		}
	}

	return logs, nil
}

// Send Mail with Excel File
func sendMailWithAttachment(logs []models.Log, f *excelize.File) {
	mailContent := &models.Mail{
		Sender:  config.C.Mail.FromMail,
		To:      toUsers,
		Cc:      ccUsers,
		Bcc:     []string{},
		Subject: config.C.App.TargetApp + " Error Logs",
	}

	mail.SendMail(mailContent, logs, f)
}
