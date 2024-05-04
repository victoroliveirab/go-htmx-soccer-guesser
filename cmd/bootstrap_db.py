import os
import sqlite3

SCHEMA_FILE = "scripts/create_schema.sql"
USER_GROUP_FILE = "scripts/create_users_and_groups.sql"
TEAMS_FILE = "scripts/create_teams.sql"
LEAGUES_FILE = "scripts/create_leagues.sql"
VENUES_FILE = "scripts/create_venues.sql"
FIXTURES_FILE = "scripts/create_fixtures.sql"
GUESSES_FILE = "scripts/create_guesses.sql"


DB_FILE = "local.db"

os.remove(DB_FILE)

conn = sqlite3.connect(DB_FILE)
cursor = conn.cursor()

with open(SCHEMA_FILE, "r") as f:
    sql_script = f.read()

# 1- Create database
cursor.executescript(sql_script)
conn.commit()
print("Created sqlite database")

# 2- Create Users and groups
with open(USER_GROUP_FILE, "r") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()
print("Created users and groups")

# 3- Create Teams
with open(TEAMS_FILE, "r") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()
print("Created teams")

# 4- Create Leagues
with open(LEAGUES_FILE, "r") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()
print("Created leagues")

# 5- Create Venues
with open(VENUES_FILE, "r") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()
print("Created venues")

# 6- Create fixtures
with open(FIXTURES_FILE, "r") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()

print("Created fixtures")

# 7- Create guesses
with open(GUESSES_FILE, "r") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()

print("Created guesses")


conn.commit()
conn.close()
