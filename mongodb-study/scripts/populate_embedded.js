// MongoDB Script: Populate Embedded Schema for WASAText
// This script generates synthetic but realistic data for the embedded document model
// Run with: mongosh --file populate_embedded.js

// Database connection and setup
const dbName = "wasatext_embedded";
use(dbName);

// Clear existing collections
db.users.drop();
db.conversations.drop();

print("🚀 Starting WASAText Embedded Schema Population");
print("📊 Generating synthetic messaging data...");

// Utility functions for realistic data generation
function randomChoice(array) {
  return array[Math.floor(Math.random() * array.length)];
}

function randomDate(start, end) {
  return new Date(start.getTime() + Math.random() * (end.getTime() - start.getTime()));
}

function generateUsername(firstName, lastName) {
  const variations = [
    `${firstName.toLowerCase()}_${lastName.toLowerCase()}`,
    `${firstName.toLowerCase()}${lastName.toLowerCase()}`,
    `${firstName.toLowerCase()}.${lastName.toLowerCase()}`,
    `${firstName.toLowerCase()}${Math.floor(Math.random() * 999)}`
  ];
  return randomChoice(variations);
}

// Sample data arrays*
