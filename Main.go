package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type Promotion struct {
	ID             string `json:"id"`
	Price          string `json:"price"`
	ExpirationDate string `json:"expiration_date"`
}

type PromotionsMap struct {
	sync.RWMutex
	Data map[string]Promotion
}

var (
	db         *sql.DB
	promotions PromotionsMap
)

func main() {
	// Initialize the database connection
	var err error
	db, err = sql.Open("postgres", "postgres://dgalstyan:dgalstyan@localhost:5432/promotionsdb?sslmode=disable")
	if err != nil {

		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/promotions/", handlePromotionRequest)

	go func() {
		// Start the server
		log.Println("Server listening on port 1321")

		log.Fatal(http.ListenAndServe(":1321", nil))
	}()

	// Call the processCSVFilesPeriodically function synchronously
	processCSVFilesPeriodically("csv_dir", 1*time.Second)

	// Wait indefinitely
	select {}
}

// handlePromotionRequest handles the promotion requests
func handlePromotionRequest(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/promotions/"):]
	promotion, err := getPromotionByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Promotion not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(promotion)
}

func getPromotionByID(id string) (*Promotion, error) {
	row := db.QueryRow("SELECT cvs_id, price, expiration_date FROM promotions WHERE id = $1", id)

	var promotion Promotion
	err := row.Scan(&promotion.ID, &promotion.Price, &promotion.ExpirationDate)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("promotion not found")
	} else if err != nil {
		return nil, err
	}

	return &promotion, nil
}

// processCSVFilesPeriodically periodically checks for new CSV files in the specified directory and updates the database
func processCSVFilesPeriodically(directory string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := processCSVFiles(directory)
			if err != nil {
				log.Println("Error processing CSV files:", err)
			}
		}
	}
}

// processCSVFiles processes CSV files in the specified directory and updates the database
func processCSVFiles(directory string) error {
	files, err := getCSVFilesInDirectory(directory)
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", directory, file.Name())

		err := processCSVFile(filePath)
		if err != nil {
			log.Println("Error processing CSV file:", err)
			continue
		}

		err = os.Remove(filePath)
		if err != nil {
			log.Println("Error removing CSV file:", err)
		}
	}

	return nil
}

// getCSVFilesInDirectory returns a list of CSV files in the specified directory
func getCSVFilesInDirectory(directory string) ([]os.FileInfo, error) {
	dir, err := os.Open(directory)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var csvFiles []os.FileInfo
	for _, file := range files {
		if !file.IsDir() && isCSVFile(file.Name()) {
			csvFiles = append(csvFiles, file)
		}
	}

	return csvFiles, nil
}

// isCSVFile checks if the given file name has a CSV file extension
func isCSVFile(filename string) bool {
	return filepath.Ext(filename) == ".csv"
}

// processCSVFile processes a CSV file and updates the database
func processCSVFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO promotions(cvs_id, price, expiration_date) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		promotion := Promotion{
			ID:             record[0],
			Price:          record[1],
			ExpirationDate: record[2],
		}

		_, err = stmt.Exec(promotion.ID, promotion.Price, promotion.ExpirationDate)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
