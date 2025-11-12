#!/bin/bash

# Bulk upload script for US SSA name data files
# Usage: ./scripts/bulk-upload.sh <country_id>

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COUNTRY_ID=$1
API_URL="${API_URL:-http://localhost:8080}"
DATA_DIR="sample-data/real-us-data"
WAIT_TIME=2  # Seconds to wait between uploads

# Validate input
if [ -z "$COUNTRY_ID" ]; then
    echo -e "${RED}Error: Country ID is required${NC}"
    echo "Usage: $0 <country_id>"
    echo ""
    echo "Example:"
    echo "  $0 550e8400-e29b-41d4-a716-446655440000"
    echo ""
    echo "Get your country_id by creating a country first:"
    echo "  curl -X POST http://localhost:8080/v1/countries \\"
    echo "    -H 'Content-Type: application/json' \\"
    echo "    -d '{\"code\":\"US\",\"name\":\"United States\"}'"
    exit 1
fi

# Check if data directory exists
if [ ! -d "$DATA_DIR" ]; then
    echo -e "${RED}Error: Data directory not found: $DATA_DIR${NC}"
    exit 1
fi

# Check if API is reachable
echo -e "${BLUE}Checking API health...${NC}"
if ! curl -s -f "$API_URL/health" > /dev/null; then
    echo -e "${RED}Error: API is not reachable at $API_URL${NC}"
    echo "Make sure the services are running: docker-compose up -d"
    exit 1
fi
echo -e "${GREEN}✓ API is healthy${NC}"
echo ""

# Count total files
TOTAL_FILES=$(ls -1 "$DATA_DIR"/yob*.txt 2>/dev/null | wc -l | tr -d ' ')
if [ "$TOTAL_FILES" -eq 0 ]; then
    echo -e "${RED}Error: No data files found in $DATA_DIR${NC}"
    exit 1
fi

echo -e "${BLUE}Found $TOTAL_FILES files to upload${NC}"
echo -e "${BLUE}Country ID: $COUNTRY_ID${NC}"
echo -e "${BLUE}API URL: $API_URL${NC}"
echo ""

# Initialize counters
SUCCESS_COUNT=0
FAIL_COUNT=0
CURRENT=0

# Start time
START_TIME=$(date +%s)

# Loop through all yob*.txt files
for file in "$DATA_DIR"/yob*.txt; do
    CURRENT=$((CURRENT + 1))
    FILENAME=$(basename "$file")
    
    echo -e "${YELLOW}[$CURRENT/$TOTAL_FILES]${NC} Uploading $FILENAME..."
    
    # Upload file
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
        -X POST "$API_URL/v1/datasets/upload" \
        -F "file=@$file" \
        -F "country_id=$COUNTRY_ID")
    
    # Check response
    if [ "$HTTP_CODE" -eq 202 ]; then
        echo -e "${GREEN}✓ Success${NC} (HTTP $HTTP_CODE)"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        echo -e "${RED}✗ Failed${NC} (HTTP $HTTP_CODE)"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    # Wait between uploads (except for the last file)
    if [ "$CURRENT" -lt "$TOTAL_FILES" ]; then
        sleep $WAIT_TIME
    fi
done

# End time
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

# Print summary
echo ""
echo "=========================================="
echo -e "${BLUE}Upload Summary${NC}"
echo "=========================================="
echo -e "Total files:    $TOTAL_FILES"
echo -e "${GREEN}Successful:     $SUCCESS_COUNT${NC}"
if [ "$FAIL_COUNT" -gt 0 ]; then
    echo -e "${RED}Failed:         $FAIL_COUNT${NC}"
else
    echo -e "Failed:         $FAIL_COUNT"
fi
echo -e "Duration:       ${DURATION}s"
echo "=========================================="
echo ""

if [ "$FAIL_COUNT" -eq 0 ]; then
    echo -e "${GREEN}✓ All files uploaded successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Check job status:"
    echo "     curl $API_URL/v1/jobs"
    echo ""
    echo "  2. Search for names:"
    echo "     curl '$API_URL/v1/names/search?query=Emma&limit=10'"
    echo ""
    echo "  3. Get top names by year:"
    echo "     curl '$API_URL/v1/names/top?year=2020&limit=10'"
    echo ""
else
    echo -e "${YELLOW}⚠ Some uploads failed. Check the logs above for details.${NC}"
    exit 1
fi