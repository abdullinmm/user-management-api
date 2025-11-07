#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== User Management API Testing ==="
echo ""

# 1. Health Check
echo "1. Health Check"
curl -s $BASE_URL/health | jq
echo ""

# 2. List Tasks (public)
echo "2. List Active Tasks"
curl -s $BASE_URL/api/v1/tasks | jq
echo ""

# 3. Register User
echo "3. Register User"
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser1"}')
echo $REGISTER_RESPONSE | jq
USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.id')
echo "Created User ID: $USER_ID"
echo ""

echo "=== Quick Manual Tests ==="
echo ""
echo "Test protected endpoints manually:"
echo ""
echo "# Generate JWT at https://jwt.io"
echo "Payload: {\"user_id\": $USER_ID, \"exp\": $(date -v+1d +%s 2>/dev/null || date -d '+1 day' +%s)}"
echo "Secret: dev_secret"
echo ""
echo "# Get User Status"
echo "curl -H 'Authorization: Bearer YOUR_TOKEN' $BASE_URL/api/v1/users/$USER_ID/status | jq"
echo ""
echo "# Complete Task"
echo "curl -X POST -H 'Authorization: Bearer YOUR_TOKEN' -H 'Content-Type: application/json' \\"
echo "  -d '{\"task_id\":1}' $BASE_URL/api/v1/users/$USER_ID/task/complete | jq"
echo ""
echo "# Get Leaderboard"
echo "curl -H 'Authorization: Bearer YOUR_TOKEN' $BASE_URL/api/v1/users/leaderboard?limit=5 | jq"
echo ""
