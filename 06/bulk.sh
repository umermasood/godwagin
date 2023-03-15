#!/bin/bash

while IFS= read -r thread
do
  printf "\n%s\n" "$thread"
  curl -X POST -H "Content-Type: application/json" -d "{\"url\":\"$thread\"}" http://localhost:5000/parse
done < "threads"