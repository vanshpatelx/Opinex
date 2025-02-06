#!/bin/bash
echo "🔧 Modifying pg_hba.conf to allow trust authentication..."

# Locate the pg_hba.conf file
PG_HBA_FILE="/var/lib/postgresql/data/pg_hba.conf"

# Update authentication method from 'md5' to 'trust'
sed -i 's/md5/trust/g' "$PG_HBA_FILE"

echo "✅ Authentication changed to trust mode."

# Restart PostgreSQL to apply changes
pg_ctl reload
