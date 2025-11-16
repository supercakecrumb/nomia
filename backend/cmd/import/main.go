package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type NameRecord struct {
	Year      int
	Name      string
	Gender    string
	Count     int
	CountryID int
	DatasetID int
}

func main() {
	fmt.Println("🚀 Starting US Name Data Import...")

	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://postgres:postgres@localhost:5432/affirm_name?sslmode=disable"
		fmt.Println("⚠️  DATABASE_URL not set, using default:", databaseURL)
	}

	// Connect to database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	fmt.Println("✅ Connected to database")

	// Ensure US country exists
	countryID, err := ensureUSCountry(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to ensure US country: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ US country ID: %d\n", countryID)

	// Find all SSA data files in names-example directory (at project root)
	dataDir := "../names-example"
	files, err := filepath.Glob(filepath.Join(dataDir, "yob*.txt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to find data files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "❌ No data files found in %s\n", dataDir)
		os.Exit(1)
	}

	fmt.Printf("📂 Found %d data files\n", len(files))

	// Process each file
	totalRecords := 0
	for _, filePath := range files {
		fmt.Printf("\n📄 Processing %s...\n", filepath.Base(filePath))

		// Extract year from filename
		year, err := extractYearFromFilename(filepath.Base(filePath))
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Skipping file %s: %v\n", filePath, err)
			continue
		}

		// Check if dataset already exists
		exists, err := datasetExists(ctx, conn, countryID, year, filepath.Base(filePath))
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to check dataset: %v\n", err)
			continue
		}
		if exists {
			fmt.Printf("⏭️  Dataset for year %d already exists, skipping\n", year)
			continue
		}

		// Create dataset record
		datasetID, err := insertDataset(ctx, conn, countryID, year, filepath.Base(filePath), filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to create dataset: %v\n", err)
			continue
		}
		fmt.Printf("✅ Created dataset ID: %d for year %d\n", datasetID, year)

		// Parse file
		records, err := parseSSAFile(filePath, year, countryID, datasetID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to parse file: %v\n", err)
			continue
		}
		fmt.Printf("📊 Parsed %d records\n", len(records))

		// Batch insert records
		err = batchInsertNames(ctx, conn, records)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to insert records: %v\n", err)
			continue
		}

		totalRecords += len(records)
		fmt.Printf("✅ Imported %d records for year %d\n", len(records), year)
	}

	fmt.Printf("\n🎉 Import complete! Total records imported: %d\n", totalRecords)
}

// ensureUSCountry ensures the US country record exists and returns its ID
func ensureUSCountry(ctx context.Context, conn *pgx.Conn) (int, error) {
	var countryID int
	err := conn.QueryRow(ctx, `
		SELECT id FROM countries WHERE code = 'US'
	`).Scan(&countryID)

	if err != nil {
		return 0, fmt.Errorf("US country not found in database. Please run migrations first")
	}

	return countryID, nil
}

// extractYearFromFilename extracts the year from SSA filename format (yob2023.txt)
func extractYearFromFilename(filename string) (int, error) {
	re := regexp.MustCompile(`yob(\d{4})\.txt`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) != 2 {
		return 0, fmt.Errorf("filename does not match expected format 'yobYYYY.txt'")
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid year in filename: %v", err)
	}

	return year, nil
}

// datasetExists checks if a dataset for the given year already exists
func datasetExists(ctx context.Context, conn *pgx.Conn, countryID, year int, filename string) (bool, error) {
	var count int
	err := conn.QueryRow(ctx, `
		SELECT COUNT(*) FROM name_datasets 
		WHERE country_id = $1 AND year_from = $2 AND source_file_name = $3
	`, countryID, year, filename).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// insertDataset creates a dataset record and returns its ID
func insertDataset(ctx context.Context, conn *pgx.Conn, countryID, year int, filename, storagePath string) (int, error) {
	var datasetID int
	err := conn.QueryRow(ctx, `
		INSERT INTO name_datasets (
			country_id, 
			source_file_name, 
			year_from, 
			year_to, 
			file_type, 
			storage_path,
			parse_status,
			uploaded_at,
			parsed_at
		) VALUES ($1, $2, $3, $3, 'SSA-TXT', $4, 'parsed', NOW(), NOW())
		RETURNING id
	`, countryID, filename, year, storagePath).Scan(&datasetID)

	return datasetID, err
}

// parseSSAFile parses an SSA format file (name,gender,count) and returns records
func parseSSAFile(path string, year, countryID, datasetID int) ([]NameRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []NameRecord
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse CSV: name,gender,count
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid format at line %d: expected 3 fields, got %d", lineNum, len(parts))
		}

		name := strings.TrimSpace(parts[0])
		gender := strings.TrimSpace(parts[1])
		countStr := strings.TrimSpace(parts[2])

		// Validate gender
		if gender != "F" && gender != "M" {
			return nil, fmt.Errorf("invalid gender at line %d: %s (expected F or M)", lineNum, gender)
		}

		// Parse count
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return nil, fmt.Errorf("invalid count at line %d: %v", lineNum, err)
		}

		if count <= 0 {
			return nil, fmt.Errorf("invalid count at line %d: must be positive", lineNum)
		}

		records = append(records, NameRecord{
			Year:      year,
			Name:      name,
			Gender:    gender,
			Count:     count,
			CountryID: countryID,
			DatasetID: datasetID,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// batchInsertNames efficiently inserts records using pgx.CopyFrom
func batchInsertNames(ctx context.Context, conn *pgx.Conn, records []NameRecord) error {
	if len(records) == 0 {
		return nil
	}

	// Use COPY for maximum performance
	copyCount, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"names"},
		[]string{"country_id", "dataset_id", "year", "name", "gender", "count"},
		pgx.CopyFromSlice(len(records), func(i int) ([]interface{}, error) {
			r := records[i]
			return []interface{}{r.CountryID, r.DatasetID, r.Year, r.Name, r.Gender, r.Count}, nil
		}),
	)

	if err != nil {
		return fmt.Errorf("copy failed: %v", err)
	}

	// Show progress for large imports
	if len(records) != int(copyCount) {
		return fmt.Errorf("expected to copy %d rows, but copied %d", len(records), copyCount)
	}

	// Log progress every 1000 records
	if len(records) > 1000 {
		fmt.Printf("   💾 Inserted %d records in %.2fs\n", copyCount, float64(time.Now().Unix()))
	}

	return nil
}
