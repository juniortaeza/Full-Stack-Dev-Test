package main

import "encoding/json"

// capture Customer data structure
type Customer struct {
	ID              string `json:"id" gorm:"primaryKey"`
	Name            string `json:"name"`
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	PropertyType    string `json:"propertyType"`
	SquareFootage   *int   `json:"squareFootage"`
	SystemType      string `json:"systemType"`
	SystemAge       *int   `json:"systemAge"`
	LastServiceDate string `json:"lastServiceDate"`
}

// modify Unmarshal to categorize SquareFootage and Sqft
func (c *Customer) UnmarshalJSON(data []byte) error {
	type Alias Customer

	temp := struct {
		*Alias
		PropertyTypeAlt string `json:"property_type"`
		Sqft            *int   `json:"sqft"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if c.PropertyType == "" {
		c.PropertyType = temp.PropertyTypeAlt
	}

	if c.SquareFootage == nil {
		c.SquareFootage = temp.Sqft
	}

	return nil
}

// capture Equipment data structure
type Equipment struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Brand       string `json:"brand"`
	ModelNumber string `json:"modelNumber"`
	BaseCost    int    `json:"baseCost"`
}

// modify Unmarshal to categorize baseCost and base_cost
func (e *Equipment) UnmarshalJSON(data []byte) error {
	type Alias Equipment

	temp := struct {
		*Alias
		BaseCostAlt *int `json:"base_cost"`
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if e.BaseCost == 0 && temp.BaseCostAlt != nil {
		e.BaseCost = *temp.BaseCostAlt
	}

	return nil
}

// capture EstimatedHours data structure for LaborRate
type EstimatedHours struct {
	Min float64 `json:"min" gorm:"column:min_hours"`
	Max float64 `json:"max" gorm:"column:max_hours"`
}

// capture LaborRate data structure
type LaborRate struct {
	JobType        string         `json:"jobType" gorm:"primaryKey"`
	Level          string         `json:"level" gorm:"primaryKey"`
	HourlyRate     int            `json:"hourlyRate"`
	EstimatedHours EstimatedHours `json:"estimatedHours" gorm:"embedded"`
}
