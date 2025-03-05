#!/bin/bash

cd backend
echo "Starting backend server..."
go run ./cmd/api/ &

cd ../frontend

if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi

echo "Starting frontend development server..."
npm run dev