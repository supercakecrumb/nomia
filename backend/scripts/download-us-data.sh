#!/bin/bash
set -e

echo "========================================="
echo "Affirm Name - US SSA Data Downloader"
echo "========================================="
echo ""

# Create directory
echo "1. Creating directory structure..."
mkdir -p ../names-example/us

# Download the full dataset (all years)
echo "2. Downloading names.zip (1880-2024, all years)..."
curl -L "https://www.ssa.gov/oact/babynames/names.zip" -o ../names-example/us/names.zip

# Check if download was successful
if [ ! -f ../names-example/us/names.zip ]; then
    echo "   ERROR: Download failed"
    exit 1
fi

echo "   Download complete ($(du -h ../names-example/us/names.zip | cut -f1))"

# Extract
echo "3. Extracting files..."
cd ../names-example/us
unzip -o names.zip
rm names.zip

# Count files
FILE_COUNT=$(ls -1 yob*.txt 2>/dev/null | wc -l)
echo "4. Extraction complete!"
echo ""
echo "========================================="
echo "Downloaded $FILE_COUNT year files"
echo "Files location: names-example/us/"
echo "========================================="
echo ""
echo "Next steps:"
echo "  1. Run import script to load data into database:"
echo "     bash backend/scripts/import-us-data.sh all"
echo ""
echo "  2. Or import specific years:"
echo "     bash backend/scripts/import-us-data.sh 2020 2024"
echo ""