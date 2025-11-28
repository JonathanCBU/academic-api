#!/bin/bash

# Default values
NUM_ENTRIES=5
API_URL="http://127.0.0.1:8080"
AUTH_TOKEN="12345"

# Parse command line arguments
while getopts "n:u:t:h" opt; do
  case $opt in
    n)
      NUM_ENTRIES=$OPTARG
      ;;
    u)
      API_URL=$OPTARG
      ;;
    t)
      AUTH_TOKEN=$OPTARG
      ;;
    h)
      echo "Usage: $0 [-n num_entries] [-u api_url] [-t auth_token]"
      echo "  -n: Number of entries to create for each type (default: 10)"
      echo "  -u: API base URL (default: http://127.0.0.1:8080)"
      echo "  -t: Authorization token (default: 12345)"
      exit 0
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
  esac
done

# Arrays for random data generation
STATE_CODES=("MA" "TX" "AZ")
SCHOOL_NAMES=("Lincoln Elementary" "Washington High" "Jefferson Middle" "Roosevelt Academy" "Kennedy School" "Madison Institute" "Monroe Education Center" "Adams Learning Academy" "Jackson Preparatory" "Wilson Charter")
SUBJECTS=("ela" "math")
GRADE_LEVELS=("3" "4" "5" "6" "7" "8")
DEMOGRAPHIC_GROUPS=("all" "black" "hispanic" "economically_disadvantaged")

echo "========================================="
echo "API Dummy Data Generator"
echo "========================================="
echo "Creating $NUM_ENTRIES entries of each type"
echo "API URL: $API_URL"
echo ""

# Array to store created school IDs
SCHOOL_IDS=()

# Create Schools
echo "Creating Schools..."
for i in $(seq 1 $NUM_ENTRIES); do
  STATE=${STATE_CODES[$RANDOM % ${#STATE_CODES[@]}]}
  SCHOOL_NAME="${SCHOOL_NAMES[$RANDOM % ${#SCHOOL_NAMES[@]}]} $i"
  
  RESPONSE=$(curl -s -X POST "$API_URL/schools/put" \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -d "{\"state_code\": \"$STATE\", \"school_name\": \"$SCHOOL_NAME\"}")
  
  # Try to extract school ID from response (adjust based on your API's response format)
  SCHOOL_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*' | head -1)
  
  if [ -z "$SCHOOL_ID" ]; then
    SCHOOL_ID=$i
  fi
  
  SCHOOL_IDS+=($SCHOOL_ID)
  echo "  Created school: $SCHOOL_NAME (ID: $SCHOOL_ID) in $STATE"
done

echo ""
echo "Created ${#SCHOOL_IDS[@]} schools"
echo ""

# Create School Reports
echo "Creating School Reports..."
for i in $(seq 1 $NUM_ENTRIES); do
  # Use random school ID from created schools
  SCHOOL_ID=${SCHOOL_IDS[$RANDOM % ${#SCHOOL_IDS[@]}]}
  DATA_ID=$((RANDOM % 100 + 1))
  ACADEMIC_YEAR=$((2020 + RANDOM % 5))
  SUBJECT=${SUBJECTS[$RANDOM % ${#SUBJECTS[@]}]}
  GRADE=${GRADE_LEVELS[$RANDOM % ${#GRADE_LEVELS[@]}]}
  DEMOGRAPHIC=${DEMOGRAPHIC_GROUPS[$RANDOM % ${#DEMOGRAPHIC_GROUPS[@]}]}
  N_TESTED=$((RANDOM % 200 + 50))
  N_PROFICIENT=$((RANDOM % N_TESTED))
  
  curl -s -X POST "$API_URL/school-reports/put" \
    -H "Authorization: Bearer $AUTH_TOKEN" \
    -d "{\"school_id\": $SCHOOL_ID, \"data_id\": $DATA_ID, \"academic_year\": $ACADEMIC_YEAR, \"subject\": \"$SUBJECT\", \"grade_level\": \"$GRADE\", \"demographic_group\": \"$DEMOGRAPHIC\", \"n_tested\": $N_TESTED, \"n_proficient\": $N_PROFICIENT}" \
    > /dev/null
  
  echo "  Created report: School $SCHOOL_ID | $ACADEMIC_YEAR | Grade $GRADE | $SUBJECT | $DEMOGRAPHIC | $N_PROFICIENT/$N_TESTED proficient"
done

echo ""
echo "========================================="
echo "Data generation complete!"
echo "Created $NUM_ENTRIES schools and $NUM_ENTRIES school reports"
echo "========================================="