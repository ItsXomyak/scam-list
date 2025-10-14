#!/bin/bash

URL="http://localhost:4000/parse-domain"
DOMAIN="youtube-com"
OUTDIR="./responses"

# Создаём папку для результатов (если её нет)
mkdir -p "$OUTDIR"

for i in {1..10}
do
  echo "Request #$i"
  curl -s -X POST "$URL" \
       -H "Content-Type: application/json" \
       -d "{\"domain\": \"$DOMAIN\"}" \
       > "$OUTDIR/response_$i.json" &
done

# Ждём завершения всех фоновых процессов
wait
echo "All requests finished. Results saved in $OUTDIR/"
