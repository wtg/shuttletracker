To import the sample database:
1. Run mongoDB first
2. ./import.sh

Use 

chmod +x import.sh

to fix the permission problem
=====================================================================================
import the sample database:

CAUTIOUS!
   import.sh **WILL REPLACE** existing records with same _id as records in sample database!

if you don't want to replace the existing records:
   **REMOVE --upsert** from every mongoimport statement 

=====================================================================================
export the shuttle tracking database from mongoDB:
running export.sh

=====================================================================================
Notice:
the scripts will only dump to json format.
if BSON is used in the database, the JSON result is not guaranteed to be matched with records in database.
