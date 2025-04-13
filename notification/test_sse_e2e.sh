#!/bin/bash
# End-to-end SSE test script

set -e

BASE_URL="http://localhost:3000"
USER_ID="e2euser"
SSE_URL="$BASE_URL/api/notifications/subscribe?userId=$USER_ID"
NOTIFY_URL="$BASE_URL/api/notifications"
TMP_SSE_OUT="sse_output_$$.txt"

echo "1. Subscribing to SSE endpoint..."
curl -N "$SSE_URL" > "$TMP_SSE_OUT" 2>&1 &
SSE_PID=$!

# Give curl a moment to connect
sleep 1

echo "2. Sending notification..."
NOTIF_ID="notif-$(date +%s%N)"
curl -s -X POST "$NOTIFY_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"id\": \"$NOTIF_ID\",
    \"timestamp\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\",
    \"title\": \"E2E Test Notification\",
    \"message\": \"This is an end-to-end SSE test\",
    \"priority\": \"normal\",
    \"read\": false,
    \"recipients\": [{\"type\": \"user\", \"id\": \"$USER_ID\"}]
  }"

echo "3. Waiting for notification via SSE..."
TIMEOUT=10
FOUND=0
for i in $(seq 1 $TIMEOUT); do
  if grep -q "E2E Test Notification" "$TMP_SSE_OUT"; then
    FOUND=1
    break
  fi
  sleep 1
done

kill $SSE_PID || true

if [ $FOUND -eq 1 ]; then
  echo "✅ Notification received via SSE!"
  grep "E2E Test Notification" "$TMP_SSE_OUT"
  rm "$TMP_SSE_OUT"
  exit 0
else
  echo "❌ Notification NOT received via SSE within $TIMEOUT seconds."
  cat "$TMP_SSE_OUT"
  rm "$TMP_SSE_OUT"
  exit 1
fi
