import random
import sqlite3
from datetime import datetime
from pprint import pprint

GOALS_PROBABILITIES = [0.35, 0.30, 0.25, 0.10]
OUTCOMES = [0, 1, 2, 3]

now = datetime.now().timestamp()


def random_number_of_goals():
    return random.choices(OUTCOMES, GOALS_PROBABILITIES)[0]


# This should come from the groups table, but for now let's do this way
POINTS = {
    "perfect": 20,
    "diff+winner": 15,
    "wg+winner": 12,
    "lg+winner": 11,
    "winner": 10,
    "draw": 15,
    "1g+draw": 4,
    "opposite": -10,
    "diff+opposite": -5,
    "none": 0,
}


DB_FILE = "local.db"

conn = sqlite3.connect(DB_FILE)
cursor = conn.cursor()

# 1- Get all fixtures
cursor.execute("SELECT * FROM Fixtures")
fixture_rows = cursor.fetchall()

fixtures = [
    {
        "id": f_id,
        "status": f_status,
        "home_score": f_hs,
        "away_score": f_as,
        "timestamp_numb": f_ts,
    }
    for f_id, _, _, _, _, _, f_ts, _, f_status, _, f_hs, f_as, *_ in fixture_rows
]

# 2- Get all users and groups
cursor.execute("SELECT * FROM Users")
users_rows = cursor.fetchall()

cursor.execute("SELECT * FROM Groups")
groups_rows = cursor.fetchall()

# 3- Expand users and groups for the groups they participate + individual guess
cursor.execute("SELECT * FROM User_Groups")
user_groups_rows = cursor.fetchall()

ugr_index = 0
users_guesses = []
for user_id, *_ in users_rows:
    from_agg_table_user_id, from_agg_table_group_id = user_groups_rows[ugr_index]
    while from_agg_table_user_id == user_id:
        users_guesses.append(
            {
                "user_id": user_id,
                "group_id": from_agg_table_group_id,
            }
        )
        ugr_index += 1
        if ugr_index == len(user_groups_rows):
            break
        from_agg_table_user_id, from_agg_table_group_id = user_groups_rows[ugr_index]

print("Combos:")
pprint(users_guesses)

for combo in users_guesses:
    user_id = combo["user_id"]
    group_id = combo["group_id"]
    for fixture in fixtures:
        if fixture["timestamp_numb"] > now + 24 * 60 * 60:
            print(f"Skipping fixture {fixture['id']}")
            continue
        guessed_hg = random_number_of_goals()
        guessed_ag = random_number_of_goals()

        fixture_id = fixture["id"]
        fixture_status = fixture["status"]

        locked = fixture_status in (0, 2)

        cursor.execute(
            """
                INSERT INTO Guesses(
                user_id, group_id, fixture_id, locked, home_goals, away_goals
                ) VALUES (?, ?, ?, ?, ?, ?)
            """,
            (user_id, group_id, fixture_id, locked, guessed_hg, guessed_ag),
        )
        conn.commit()
