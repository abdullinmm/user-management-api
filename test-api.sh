#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== User Management API Testing ==="
echo ""

# 1. Health Check
echo "1. Health Check"
curl -s $BASE_URL/health | jq
echo ""

# 2. List Tasks
echo "2. List Active Tasks"
curl -s $BASE_URL/api/v1/tasks | jq
echo ""

# 3. Register User
echo "3. Register User"
RANDOM_USERNAME="testuser$RANDOM"
USER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$RANDOM_USERNAME\"}")
echo $USER_RESPONSE | jq
USER_ID=$(echo $USER_RESPONSE | jq -r '.id')
echo "Created User ID: $USER_ID"
echo ""

echo "=== Generate JWT Token ==="
echo ""
echo "Run this command to generate JWT token:"
echo ""
echo "cat > gen_token.go << 'GOEOF'"
echo 'package main'
echo 'import ('
echo '  "fmt"'
echo '  "time"'
echo '  "github.com/golang-jwt/jwt/v5"'
echo ')'
echo 'func main() {'
echo "  claims := jwt.MapClaims{\"user_id\": $USER_ID, \"exp\": time.Now().Add(24 * time.Hour).Unix()}"
echo '  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)'
echo '  tokenString, _ := token.SignedString([]byte("your-secret-key-change-in-production"))'
echo '  fmt.Println(tokenString)'
echo '}'
echo 'GOEOF'
echo ""
echo "go run gen_token.go"
echo "export TOKEN=\"\$(go run gen_token.go)\""
echo ""
echo "=== Manual Tests ==="
echo ""
echo "After generating token, test these endpoints:"
echo ""
echo "# Get User Status"
echo "curl -H \"Authorization: Bearer \$TOKEN\" $BASE_URL/api/v1/users/$USER_ID/status | jq"
echo ""
echo "# Complete Task"
echo "curl -X POST -H \"Authorization: Bearer \$TOKEN\" -H \"Content-Type: application/json\" \\"
echo "  -d '{\"task_id\":1}' $BASE_URL/api/v1/users/$USER_ID/task/complete | jq"
echo ""
echo "# Get Leaderboard"
echo "curl -H \"Authorization: Bearer \$TOKEN\" $BASE_URL/api/v1/users/leaderboard?limit=5 | jq"
echo ""
