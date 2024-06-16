import http.client
import json
import os
import sqlite3
import sys
from datetime import datetime
from pprint import pprint

API_KEY = None

try:
    API_KEY = os.environ["API_KEY"]
except KeyError:
    print("You should load API_KEY to the environment")
    sys.exit(1)

NOW = datetime.now()
NOW_TS = int(NOW.timestamp())

DB_FILE = "local.db"
SCHEMA_FILE = "scripts/create_schema.sql"
USER_GROUP_FILE = "scripts/create_users_and_groups.sql"
LEAGUES_FILE = "scripts/create_leagues.sql"

try:
    os.rename(DB_FILE, f"backup-{NOW_TS}.db")
except FileNotFoundError:
    pass

try:
    os.remove(f"{DB_FILE}-shm")
except FileNotFoundError:
    pass

try:
    os.remove(f"{DB_FILE}-wal")
except FileNotFoundError:
    pass

conn = sqlite3.connect(DB_FILE)
cursor = conn.cursor()

with open(SCHEMA_FILE, "r", encoding="utf-8") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()
print("Created sqlite database")

with open(USER_GROUP_FILE, "r", encoding="utf-8") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()
print("Created users and groups")

with open(LEAGUES_FILE, "r", encoding="utf-8") as f:
    sql_script = f.read()

cursor.executescript(sql_script)
conn.commit()
print("Created leagues")

cursor.execute("SELECT * FROM Leagues")
leagues_rows = cursor.fetchall()

leagues = {
    api_football_id: league_id for league_id, api_football_id, *_ in leagues_rows
}
teams = {}
venues = {}

files = []

one_week_ago = datetime.fromtimestamp(NOW_TS - 7 * 24 * 60 * 60).strftime("%Y-%m-%d")
one_week_ahead = datetime.fromtimestamp(NOW_TS + 7 * 24 * 60 * 60).strftime("%Y-%m-%d")

# 1- Load fixtures for desired leagues

for league_id, league_api_football_id, league_name, *_ in leagues_rows:
    leagues[league_api_football_id] = league_id
    query = f"league={league_api_football_id}&season=2024&from={one_week_ago}&to={one_week_ahead}"
    filename = os.path.join("json", f"{query}.json")
    files.append(filename)
    if os.path.exists(filename):
        print(f"{filename} exists. Skipping...")
        continue
    client = http.client.HTTPSConnection("api-football-v1.p.rapidapi.com")
    headers = {
        "X-RapidAPI-Key": API_KEY,
        "X-RapidAPI-Host": "api-football-v1.p.rapidapi.com",
    }
    print(f"Requesting 2 weeks worth of fixtures for league {league_name}")
    client.request("GET", f"/v3/fixtures?{query}", headers=headers)
    res = client.getresponse()
    with open(filename, "w", encoding="utf-8") as f:
        f.write(res.read().decode("utf-8"))

    print(f"Requesting standings for league {league_name}")
    client.request(
        "GET",
        f"/v3/standings?league={league_api_football_id}&season=2024",
        headers=headers,
    )
    res = client.getresponse()
    filename = os.path.join("json", f"standings-{league_name}-{str(NOW_TS)}.json")
    with open(filename, "w", encoding="utf-8") as f:
        f.write(res.read().decode("utf-8"))

# 2- Read each file and store the contents in DB

for json_file in files:
    with open(json_file, "r", encoding="utf-8") as f:
        fixtures = json.load(f)["response"]
        for info in fixtures:
            fixture = info["fixture"]
            venue = fixture["venue"]
            goals = info["goals"]
            league = info["league"]
            score = info["score"]
            home_team = info["teams"]["home"]
            away_team = info["teams"]["away"]

            # 2.1 Read (or create) venue
            venue_id = venues.get(venue["id"])
            if venue_id is None and venue["id"] is not None:
                cursor.execute(
                    "INSERT INTO Venues(api_football_id, name, city) VALUES(?, ?, ?)",
                    (venue["id"], venue["name"], venue["city"]),
                )
                conn.commit()
                venue_id = cursor.lastrowid
                venues[venue["id"]] = venue_id

            # 2.2 Read (or create) home team
            home_team_id = teams.get(home_team["id"])
            if home_team_id is None:
                cursor.execute(
                    "INSERT INTO Teams(api_football_id, name, logo_url) VALUES(?, ?, ?)",
                    (home_team["id"], home_team["name"], home_team["logo"]),
                )
                conn.commit()
                home_team_id = cursor.lastrowid
                teams[home_team["id"]] = home_team_id

            # 2.3 Read (or create) away team
            away_team_id = teams.get(away_team["id"])
            if away_team_id is None:
                cursor.execute(
                    "INSERT INTO Teams(api_football_id, name, logo_url) VALUES(?, ?, ?)",
                    (away_team["id"], away_team["name"], away_team["logo"]),
                )
                conn.commit()
                away_team_id = cursor.lastrowid
                teams[away_team["id"]] = away_team_id

            # 2.4 Create fixture
            api_football_id = fixture["id"]
            league_id = leagues[league["id"]]
            season = 2024
            timestamp_numb = fixture["timestamp"]
            status = {
                "FT": 0,
                "PST": 1,
                "PEN": 2,
                "NS": 3,
            }.get(fixture["status"]["short"])
            if status is None:
                print("Ignoring fixture:")
                pprint(info)
                # Game is in progress, so for now let's ignore it
                continue
            referee = fixture["referee"]
            home_score = goals["home"]
            away_score = goals["away"]
            home_winner = home_team["winner"]
            away_winner = away_team["winner"]
            league_round = league["round"]
            cursor.execute(
                """
                INSERT INTO Fixtures(
                api_football_id, league_id, season, home_team_id,
                away_team_id, timestamp_numb, venue_id, status, referee,
                home_score, away_score, home_winner, away_winner, round
                ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """,
                (
                    fixture["id"],
                    league_id,
                    season,
                    home_team_id,
                    away_team_id,
                    timestamp_numb,
                    venue_id,
                    status,
                    referee,
                    home_score,
                    away_score,
                    home_winner,
                    away_winner,
                    league_round,
                ),
            )
conn.close()
