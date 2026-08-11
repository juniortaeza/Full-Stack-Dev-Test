package main

import (
	"encoding/json"
	"os"

	// use gorm as mysql driver
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// load customers structure into database
func loadCustomers() ([]Customer, error) {
	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		return nil, err
	}

	var customers []Customer

	err = json.Unmarshal(data, &customers)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

// load equipment structure into database
func loadEquipment() ([]Equipment, error) {
	data, err := os.ReadFile("data/equipment.json")
	if err != nil {
		return nil, err
	}

	var equipment []Equipment

	err = json.Unmarshal(data, &equipment)
	if err != nil {
		return nil, err
	}

	return equipment, nil
}

// load labor rates structure into database
func loadLaborRates() ([]LaborRate, error) {
	data, err := os.ReadFile("data/labor_rates.json")
	if err != nil {
		return nil, err
	}

	var laborRates []LaborRate

	if err := json.Unmarshal(data, &laborRates); err != nil {
		return nil, err
	}

	return laborRates, nil
}

// connect backend to the database
func connectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/hvac_service"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// generate tables for each json file
func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&Customer{},
		&Equipment{},
		&LaborRate{},
	)
}

// populate customers table
func seedCustomers(db *gorm.DB, customers []Customer) error {
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&customers).Error
}

// populate equipment table
func seedEquipment(db *gorm.DB, equipment []Equipment) error {
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&equipment).Error
}

// populate labor rates table
func seedLaborRates(db *gorm.DB, laborRate []LaborRate) error {
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&laborRate).Error
}
