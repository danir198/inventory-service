#!/bin/bash

# Define the API endpoint
API_URL="http://localhost:8002/products"

# Read the JWT token from an environment variable
AUTH_TOKEN=${JWT_TOKEN}

# Check if the JWT token is set
if [ -z "$AUTH_TOKEN" ]; then
    echo "Error: JWT token is not set. Please set the JWT_TOKEN environment variable."
    exit 1
fi

# Read the JSON file and loop through each product
cat input_products.json | jq -c '.[]' | while read -r product; do
    # Make the curl request to create the product
    curl -X POST $API_URL \
        -H "Authorization: Bearer $AUTH_TOKEN" \
        -H "Content-Type: application/json" \
        -d "$product"
    echo # Add a newline for readability
done