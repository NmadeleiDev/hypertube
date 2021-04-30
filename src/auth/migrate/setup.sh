#!/bin/bash
sudo -i -u postgres psql < createDatabase.sql # Это для Linux
# psql -d postgres < createDatabase.sql # Это для MacOs