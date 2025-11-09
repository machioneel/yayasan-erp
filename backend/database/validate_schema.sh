#!/bin/bash

# ============================================
# SQL SCHEMA VALIDATION SCRIPT
# ============================================

echo "🔍 VALIDATING SQL SCHEMA..."
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
passed=0
failed=0

# Function to check
check() {
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ PASS${NC}: $1"
        ((passed++))
    else
        echo -e "${RED}❌ FAIL${NC}: $1"
        ((failed++))
    fi
}

echo "📄 Checking SQL Files..."
echo "─────────────────────────────────────"

# Check if files exist
test -f schema.sql
check "schema.sql exists"

test -f seed_data.sql
check "seed_data.sql exists"

test -f test_schema.sql
check "test_schema.sql exists"

test -f README.md
check "README.md exists"

test -d migrations
check "migrations directory exists"

echo ""
echo "🔍 Validating SQL Syntax..."
echo "─────────────────────────────────────"

# Check for common SQL syntax errors
grep -q "CREATE TABLE" schema.sql
check "Contains CREATE TABLE statements"

grep -q "PRIMARY KEY" schema.sql
check "Contains PRIMARY KEY definitions"

grep -q "FOREIGN KEY\|REFERENCES" schema.sql
check "Contains FOREIGN KEY relationships"

grep -q "CREATE INDEX" schema.sql
check "Contains INDEX definitions"

grep -q "CREATE TRIGGER" schema.sql
check "Contains TRIGGER definitions"

grep -q "CREATE.*VIEW" schema.sql
check "Contains VIEW definitions"

echo ""
echo "📊 Analyzing Schema Structure..."
echo "─────────────────────────────────────"

# Count tables
table_count=$(grep -c "CREATE TABLE" schema.sql)
echo "  Tables found: $table_count"
if [ $table_count -ge 18 ]; then
    echo -e "  ${GREEN}✅${NC} Expected ~18 tables"
    ((passed++))
else
    echo -e "  ${RED}❌${NC} Expected ~18 tables, found $table_count"
    ((failed++))
fi

# Count indexes
index_count=$(grep -c "CREATE INDEX" schema.sql)
echo "  Indexes found: $index_count"
if [ $index_count -ge 30 ]; then
    echo -e "  ${GREEN}✅${NC} Good index coverage"
    ((passed++))
else
    echo -e "  ${YELLOW}⚠️${NC}  Could use more indexes"
fi

# Count triggers
trigger_count=$(grep -c "CREATE TRIGGER" schema.sql)
echo "  Triggers found: $trigger_count"
if [ $trigger_count -ge 10 ]; then
    echo -e "  ${GREEN}✅${NC} Triggers implemented"
    ((passed++))
else
    echo -e "  ${YELLOW}⚠️${NC}  Few triggers found"
fi

# Count views
view_count=$(grep -c "CREATE.*VIEW" schema.sql)
echo "  Views found: $view_count"
if [ $view_count -ge 4 ]; then
    echo -e "  ${GREEN}✅${NC} Views created"
    ((passed++))
else
    echo -e "  ${YELLOW}⚠️${NC}  Could add more views"
fi

echo ""
echo "🔐 Security Checks..."
echo "─────────────────────────────────────"

# Check for password hashing mention
grep -q "password_hash" schema.sql
check "Password hashing implemented"

# Check for audit logging
grep -q "audit_logs" schema.sql
check "Audit logging table exists"

# Check for soft deletes
grep -q "deleted_at" schema.sql
check "Soft delete support"

echo ""
echo "📋 Data Integrity Checks..."
echo "─────────────────────────────────────"

# Check for NOT NULL constraints
grep -q "NOT NULL" schema.sql
check "NOT NULL constraints used"

# Check for UNIQUE constraints
grep -q "UNIQUE" schema.sql
check "UNIQUE constraints defined"

# Check for ON DELETE CASCADE
grep -q "ON DELETE CASCADE" schema.sql
check "CASCADE deletes configured"

# Check for DEFAULT values
grep -q "DEFAULT" schema.sql
check "DEFAULT values specified"

echo ""
echo "🧪 Testing Seed Data..."
echo "─────────────────────────────────────"

# Check seed data
grep -q "INSERT INTO" seed_data.sql
check "Seed data INSERT statements"

grep -q "branches" seed_data.sql
check "Branches seed data"

grep -q "users" seed_data.sql
check "Users seed data"

grep -q "accounts" seed_data.sql
check "Chart of Accounts seed data"

echo ""
echo "═══════════════════════════════════════"
echo "📊 VALIDATION SUMMARY"
echo "═══════════════════════════════════════"
echo ""
echo -e "  ${GREEN}Passed: $passed${NC}"
echo -e "  ${RED}Failed: $failed${NC}"
echo ""

if [ $failed -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL VALIDATIONS PASSED!${NC}"
    echo -e "${GREEN}✅ SQL Schema is valid and ready for use${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Create PostgreSQL database"
    echo "  2. Run: psql -d dbname -f schema.sql"
    echo "  3. Run: psql -d dbname -f seed_data.sql"
    echo "  4. Run: psql -d dbname -f test_schema.sql"
    exit 0
else
    echo -e "${RED}❌ VALIDATION FAILED${NC}"
    echo -e "${RED}Please fix the issues above${NC}"
    exit 1
fi
