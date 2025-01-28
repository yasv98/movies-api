#!/bin/bash

# Wait for mongodb to start up.
sleep 5

echo "Starting data import"
mongoimport --drop --db sample_mflix --collection comments --file /docker-entrypoint-initdb.d/comments.json
mongoimport --drop --db sample_mflix --collection movies --file /docker-entrypoint-initdb.d/movies.json
echo "Finished importing"

echo "Creating indexes"
mongosh sample_mflix --eval '
    // Index for title searches.
    db.movies.createIndex({"title": 1});

    // Index to speed up comment lookups by movie and comment ID.
    db.comments.createIndex({"movie_id": 1, "_id": 1});
'
echo "Finished creating indexes"