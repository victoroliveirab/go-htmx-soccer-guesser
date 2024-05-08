import os
import sqlite3
from datetime import datetime


def offset_days(dt, offset=0):
    timestamp = dt.timestamp()
    new_timestamp = timestamp + offset * 3600 * 24
    return datetime.fromtimestamp(new_timestamp)


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

# 6.1 Update fixtures to mix between past, current, and future
# 9 past (5 two days ago + 4 yesterday)
# 2 current
# 9 future (5 tomorrow + 4 two days from today)

now = datetime.now()
for i in range(20):
    row_id = i + 1
    if i < 5:  # two days ago
        two_days_ago = offset_days(now, -2)
        two_days_ago = two_days_ago.replace(minute=30, second=0, microsecond=0)

        timestamp_numb = int(two_days_ago.timestamp() * 1000)
        match_date = "T".join(str(two_days_ago).split(" ")) + "+00:00"
        cursor.execute(
            """
            UPDATE fixtures
            SET timestamp_numb = ?,
                match_date = ?
            WHERE id = ?
        """,
            (timestamp_numb, match_date, row_id),
        )
    elif i < 9:  # yesterday
        yesterday = offset_days(now, -1)
        yesterday = yesterday.replace(minute=0, second=0, microsecond=0)

        timestamp_numb = int(yesterday.timestamp() * 1000)
        match_date = "T".join(str(yesterday).split(" ")) + "+00:00"
        cursor.execute(
            """
            UPDATE fixtures
            SET timestamp_numb = ?,
                match_date = ?
            WHERE id = ?
        """,
            (timestamp_numb, match_date, row_id),
        )
    elif i < 11:  # now
        today = now.replace(minute=0, microsecond=0)
        timestamp_numb = int(today.timestamp() * 1000)
        match_date = "T".join(str(today).split(" ")) + "+00:00"
        cursor.execute(
            """
            UPDATE fixtures
            SET timestamp_numb = ?,
                match_date = ?
            WHERE id = ?
        """,
            (timestamp_numb, match_date, row_id),
        )
    elif i < 16:
        tomorrow = offset_days(now, 1)
        tomorrow = tomorrow.replace(minute=30, second=0, microsecond=0)
        timestamp_numb = int(tomorrow.timestamp() * 1000)
        match_date = "T".join(str(tomorrow).split(" ")) + "+00:00"
        cursor.execute(
            """
            UPDATE fixtures
            SET timestamp_numb = ?,
                match_date = ?
            WHERE id = ?
        """,
            (timestamp_numb, match_date, row_id),
        )
    else:
        two_days_ahead = offset_days(now, 2)
        two_days_ahead = two_days_ahead.replace(minute=30, second=0, microsecond=0)
        timestamp_numb = int(two_days_ahead.timestamp() * 1000)
        match_date = "T".join(str(two_days_ahead).split(" ")) + "+00:00"
        cursor.execute(
            """
            UPDATE fixtures
            SET timestamp_numb = ?,
                match_date = ?
            WHERE id = ?
        """,
            (timestamp_numb, match_date, row_id),
        )

print("Update Fixtures to current date")


# 7- Create guesses
with open(GUESSES_FILE, "r") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()

print("Created guesses")


conn.commit()
conn.close()
