set -euo pipefail

API_URL="http://localhost:8080/categories"

declare -a categories=(
  "Smartphone"
  "Ordinateur"
  "Tablette"
  "Ecouteur"
  "Clavier"
  "Souris"
  "Enceinte"
  "Chargeur"
)

for name in "${categories[@]}"; do
  echo -n "Creating category: $name … "
  http_status=$(  
    curl -s -o /dev/null -w "%{http_code}" \
      -X POST "$API_URL" \
      -H "Content-Type: application/json" \
      -d "{\"name\":\"${name//\"/\\\"}\"}"
  )
  if [[ "$http_status" -eq 201 ]]; then
    echo "✔️"
  elif [[ "$http_status" -eq 409 ]]; then
    echo "⚠️ already exists"
  else
    echo "❌ failed (HTTP $http_status)"
  fi
done
