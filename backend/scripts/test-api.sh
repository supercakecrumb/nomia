#!/bin/bash
# Comprehensive API testing script for Affirm Name backend

set -e

BASE_URL="${1:-http://localhost:8083}"

echo "========================================="
echo "Affirm Name Backend - API Tests"
echo "Base URL: $BASE_URL"
echo "========================================="
echo ""

echo "1. Health Check"
echo "   GET $BASE_URL/health"
curl -s "$BASE_URL/health" | jq '.'
echo ""

echo "2. Meta - Year Range"
echo "   GET $BASE_URL/api/meta/years"
curl -s "$BASE_URL/api/meta/years" | jq '.'
echo ""

echo "3. Meta - Countries"
echo "   GET $BASE_URL/api/meta/countries"
curl -s "$BASE_URL/api/meta/countries" | jq '.countries | length'
echo "   Countries available"
echo ""

echo "4. Top 10 Most Popular Names (All Time)"
echo "   GET $BASE_URL/api/names?top_n=10"
curl -s "$BASE_URL/api/names?top_n=10" | jq '.names[:5] | map({rank, name, total_count})'
echo "   (showing top 5 of 10)"
echo ""

echo "5. Unisex Names (48-52% gender balance)"
echo "   GET $BASE_URL/api/names?gender_balance_min=48&gender_balance_max=52&page_size=5"
curl -s "$BASE_URL/api/names?gender_balance_min=48&gender_balance_max=52&page_size=5" | jq '.names | map({name, gender_balance})'
echo ""

echo "6. Name Trend for 'Mary'"
echo "   GET $BASE_URL/api/names/trend?name=Mary"
MARY_DATA=$(curl -s "$BASE_URL/api/names/trend?name=Mary")
echo "$MARY_DATA" | jq '{total: .summary.total_count, years: (.summary.name_end - .summary.name_start + 1), gender: .summary.gender_balance}'
echo ""

echo "7. Name Pattern Matching 'Alex*'"
echo "   GET $BASE_URL/api/names?name_glob=Alex*&page_size=5"
curl -s "$BASE_URL/api/names?name_glob=Alex*&page_size=5" | jq '.names | map(.name)'
echo ""

echo "8. Decade Filter - 1990s Top 3"
echo "   GET $BASE_URL/api/names?year_from=1990&year_to=1999&top_n=3"
curl -s "$BASE_URL/api/names?year_from=1990&year_to=1999&top_n=3" | jq '.names | map({name, total_count})'
echo ""

echo "9. Coverage Filter - 90% of population"
echo "   GET $BASE_URL/api/names?coverage_percent=90&page_size=5"
COVERAGE_DATA=$(curl -s "$BASE_URL/api/names?coverage_percent=90&page_size=5")
echo "$COVERAGE_DATA" | jq '{names_covering_90pct: .meta.popularity_summary.derived_top_n, first_5: .names[:5] | map(.name)}'
echo ""

echo "10. Pagination Test"
echo "   GET $BASE_URL/api/names?page=1&page_size=20"
curl -s "$BASE_URL/api/names?page=1&page_size=20" | jq '{page: .meta.page, page_size: .meta.page_size, total_pages: .meta.total_pages, total_count: .meta.total_count}'
echo ""

echo "========================================="
echo "✅ All API tests complete!"
echo "========================================="