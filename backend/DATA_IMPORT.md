# Data Import Guide

This guide explains how to import baby name data from various countries into the Affirm Name database.

## Quick Start

### 1. Download US Data (All Years 1880-2024)

```bash
bash backend/scripts/download-us-data.sh
```

This downloads ~150 files covering 145 years of US baby name data (~30MB total).

### 2. Import All US Data

```bash
bash backend/scripts/import-us-data.sh all
```

Or import specific year range:

```bash
bash backend/scripts/import-us-data.sh 2020 2024
```

## Data Sources

See [`data-sources.yml`](data-sources.yml) for complete list of supported data sources.

### United States (SSA)
- **Format**: CSV (name,gender,count)
- **Years**: 1880-2024
- **Files**: One per year (yobYYYY.txt)
- **Download**: https://www.ssa.gov/oact/babynames/names.zip
- **Size**: ~150 files, ~30MB total

### United Kingdom (ONS)
- **Format**: Excel (.xlsx)
- **Years**: 1996-2022
- **Files**: Separate for boys/girls per year
- **Download**: Manual from ONS website
- **Notes**: Requires Excel parsing library

### Sweden (SCB)
- **Format**: Excel (.xlsx) or CSV
- **Years**: 1998-2023
- **Download**: Manual from SCB website
- **Notes**: Contains Swedish characters (å, ä, ö)

### Canada (StatCan)
- **Format**: CSV
- **Years**: Varies by province
- **Download**: Manual from Statistics Canada
- **Notes**: Provincial data may have different formats

## Import Tool Usage

### Basic Import

```bash
cd backend
go run cmd/import/main.go -country=US -dir=../names-example
```

### Command-Line Options

```
-country string    Country code (US, UK, SE, CA) (default "US")
-dir string        Directory containing data files (default "../names-example")
-year-from int     Start year (optional, 0 means all)
-year-to int       End year (optional, 0 means all)
-dry-run           Validate files without importing
-verbose           Show detailed progress
```

### Examples

```bash
# Import all US data
go run cmd/import/main.go -country=US -dir=../names-example/us

# Import specific years
go run cmd/import/main.go -country=US -year-from=2020 -year-to=2024

# Import single year
go run cmd/import/main.go -country=US -year-from=2023 -year-to=2023

# Dry run to check files
go run cmd/import/main.go -country=US -dry-run

# Verbose output
go run cmd/import/main.go -country=US -verbose

# Import UK data (when available)
go run cmd/import/main.go -country=UK -dir=../names-example/uk
```

## File Organization

Organize downloaded files by country:

```
names-example/
├── us/
│   ├── yob1880.txt
│   ├── yob1881.txt
│   └── ... (all years)
├── uk/
│   ├── boys-2020.xlsx
│   ├── girls-2020.xlsx
│   └── ...
├── se/
│   ├── sweden-2020.xlsx
│   └── ...
└── ca/
    ├── canada-2020.csv
    └── ...
```

## Import Process

### Step 1: Prepare Data

1. Download data files from official sources
2. Organize by country in `names-example/` directory
3. Verify file formats match expected structure

### Step 2: Seed Country Metadata

Run migrations to ensure countries exist:

```bash
# Apply migration 003
docker-compose exec postgres psql -U postgres -d affirm_name -f /migrations/003_seed_all_countries.sql
```

Or ensure countries are in database:

```sql
INSERT INTO countries (code, name, data_source_name, data_source_url, ...)
VALUES ('US', 'United States', 'SSA', 'https://...');
```

### Step 3: Run Import

```bash
bash backend/scripts/import-us-data.sh all
```

### Step 4: Verify Import

```bash
# Check imported data
docker-compose exec -T postgres psql -U postgres -d affirm_name -c "
SELECT 
    c.name as country,
    MIN(n.year) as first_year,
    MAX(n.year) as last_year,
    COUNT(DISTINCT n.year) as total_years,
    COUNT(DISTINCT n.name) as unique_names,
    COUNT(*) as total_records,
    SUM(n.count) as total_occurrences
FROM names n
JOIN countries c ON n.country_id = c.id
GROUP BY c.name
ORDER BY c.name;
"
```

## Data Validation

### Check for Duplicates

```sql
SELECT name, year, country_id, gender, COUNT(*)
FROM names
GROUP BY name, year, country_id, gender
HAVING COUNT(*) > 1;
```

### Verify Gender Balance

```sql
SELECT 
    gender,
    COUNT(*) as records,
    SUM(count) as total_count
FROM names
GROUP BY gender;
```

### Check Year Coverage

```sql
SELECT 
    year,
    COUNT(DISTINCT name) as unique_names
FROM names
GROUP BY year
ORDER BY year;
```

### Verify Data by Country

```sql
SELECT 
    c.code,
    c.name,
    COUNT(DISTINCT n.year) as years,
    COUNT(DISTINCT n.name) as unique_names,
    SUM(n.count) as total_occurrences
FROM countries c
LEFT JOIN names n ON c.id = n.country_id
GROUP BY c.code, c.name
ORDER BY c.code;
```

## Performance Tips

### Large Dataset Import

For importing 100+ years of data:

1. **Use batch imports** (already implemented in import tool using `COPY FROM`)
2. **Disable indexes temporarily** for faster inserts (optional)
3. **Import during off-peak hours** if running on shared infrastructure
4. **Monitor disk space** - full US dataset requires ~500MB-1GB database space

### Monitor Progress

```bash
# Watch import progress in real-time
watch -n 1 'docker-compose exec -T postgres psql -U postgres -d affirm_name -c "SELECT COUNT(*) FROM names;"'

# Check most recent imports
docker-compose exec -T postgres psql -U postgres -d affirm_name -c "
SELECT year, COUNT(*) as records
FROM names 
GROUP BY year 
ORDER BY year DESC 
LIMIT 10;
"
```

## Troubleshooting

### File Format Errors

Check file encoding and format:

```bash
# Check file encoding
file names-example/us/yob2023.txt

# View first few lines
head -5 names-example/us/yob2023.txt

# Check for unexpected characters
cat names-example/us/yob2023.txt | od -c | head
```

### Import Errors

Check import tool logs:

```bash
# Run with verbose flag and save logs
go run cmd/import/main.go -country=US -verbose 2>&1 | tee import.log

# Check for specific errors
grep "ERROR" import.log
grep "❌" import.log
```

### Database Connection Errors

```bash
# Check database is running
docker-compose ps

# Check database logs
docker-compose logs postgres

# Test connection manually
docker-compose exec postgres psql -U postgres -d affirm_name -c "SELECT 1;"
```

### Year Range Issues

If year filtering doesn't work as expected:

```bash
# Verify year extraction works
go run cmd/import/main.go -country=US -year-from=2023 -year-to=2023 -verbose

# Check files match pattern
ls -la names-example/us/yob202*.txt
```

## Expected Data Volumes

| Country | Years | Approx Files | Approx Records | Approx Size | Import Time |
|---------|-------|--------------|----------------|-------------|-------------|
| US      | 1880-2024 | 145 | ~2,000,000 | ~500MB | 2-5 min |
| UK      | 1996-2022 | ~54 | ~500,000 | ~150MB | 1-2 min |
| SE      | 1998-2023 | ~26 | ~300,000 | ~100MB | 1 min |
| CA      | varies | varies | varies | varies | varies |

**Note**: Import times are approximate and depend on system performance.

## Import Examples

### Example 1: Fresh Import of All US Data

```bash
# 1. Download data
bash backend/scripts/download-us-data.sh

# 2. Start database if needed
docker-compose up -d

# 3. Run migrations
docker-compose exec postgres psql -U postgres -d affirm_name -c "\i /migrations/003_seed_all_countries.sql"

# 4. Import all data
bash backend/scripts/import-us-data.sh all

# 5. Verify
docker-compose exec -T postgres psql -U postgres -d affirm_name -c "
SELECT COUNT(*) as total_records, 
       COUNT(DISTINCT year) as years, 
       COUNT(DISTINCT name) as unique_names 
FROM names;
"
```

### Example 2: Import Recent Years Only

```bash
# Import last 5 years
bash backend/scripts/import-us-data.sh 2020 2024

# Or use import tool directly
cd backend
go run cmd/import/main.go -country=US -year-from=2020 -year-to=2024 -verbose
```

### Example 3: Dry Run Before Import

```bash
# Validate files without importing
cd backend
go run cmd/import/main.go -country=US -dry-run -verbose

# If validation passes, run actual import
go run cmd/import/main.go -country=US
```

### Example 4: Re-import Single Year

```bash
# First, delete existing data for that year
docker-compose exec -T postgres psql -U postgres -d affirm_name -c "
DELETE FROM names WHERE year = 2023;
DELETE FROM name_datasets WHERE year_from = 2023;
"

# Then import
bash backend/scripts/import-us-data.sh 2023
```

## Data Format Examples

### US (SSA) - CSV Format

File: `yob2023.txt`

```csv
Mary,F,7065
Anna,F,2604
Emma,F,2003
Elizabeth,F,1998
James,M,11670
John,M,9386
Robert,M,7857
```

Format: `name,gender,count` (no header row)

### UK (ONS) - Excel Format

File: `boys-2020.xlsx`

| Rank | Name | Count |
|------|------|-------|
| 1 | Oliver | 6623 |
| 2 | George | 5912 |
| 3 | Noah | 5471 |

**Note**: UK data requires Excel parsing (not yet implemented in import tool).

## Future Enhancements

- [ ] Automatic download from official APIs
- [ ] Excel file parsing for UK/Sweden data
- [ ] Incremental updates (only new years)
- [ ] Progress bars for CLI (using tqdm or similar)
- [ ] Data quality validation rules
- [ ] Automatic deduplication
- [ ] Parallel imports for multiple countries
- [ ] Support for additional countries (Australia, France, Germany, etc.)
- [ ] Web interface for data import monitoring
- [ ] Scheduled automatic updates

## Resources

### Official Data Sources

- **US**: https://www.ssa.gov/oact/babynames/
- **UK**: https://www.ons.gov.uk/
- **Sweden**: https://www.scb.se/
- **Canada**: https://www.statcan.gc.ca/

### Related Documentation

- [`README.md`](../README.md) - Project overview
- [`DATABASE.md`](DATABASE.md) - Database schema details
- [`data-sources.yml`](data-sources.yml) - Data source specifications
- [`ARCHITECTURE.md`](../ARCHITECTURE.md) - System architecture

### Support

For issues or questions:
1. Check this documentation
2. Review import logs with `-verbose` flag
3. Check database logs: `docker-compose logs postgres`
4. Verify data files match expected format
5. Open an issue in the project repository

---

**Last Updated**: 2024-11-16