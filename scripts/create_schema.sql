-- Enable Foreign Key Support (if needed in your SQLite environment)
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_football_id INTEGER NOT NULL UNIQUE,
    name TEXT NOT NULL,
    logo_url TEXT  -- optional logo for the team
);

CREATE TABLE IF NOT EXISTS Leagues (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_football_id INTEGER NOT NULL UNIQUE,
  name TEXT NOT NULL,
  logo_url TEXT,
  country TEXT,
  country_flag_url TEXT,
  league_type TEXT,
  meta TEXT
);

CREATE TABLE IF NOT EXISTS Leagues_Seasons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    league_id INTEGER NOT NULL,
    season INTEGER NOT NULL,
    standings TEXT,
    FOREIGN KEY (league_id) REFERENCES Leagues(id)
);

CREATE TABLE IF NOT EXISTS Venues (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_football_id INTEGER NOT NULL UNIQUE,
  name TEXT NOT NULL,
  city TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Odds (
    fixture_id INTEGER PRIMARY KEY,
    home_win_odd INTEGER NOT NULL,
    draw_odd INTEGER NOT NULL,
    away_win_odd INTEGER NOT NULL,
    FOREIGN KEY (fixture_id) REFERENCES Fixtures(id)
);

CREATE TABLE IF NOT EXISTS Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now'))
);

CREATE TABLE IF NOT EXISTS Groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    admin INTEGER NOT NULL,
    points_table TEXT NOT NULL,
    ranking TEXT NOT NULL,
    ranking_up_to_date INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (admin) REFERENCES Users(id)
);

CREATE TABLE IF NOT EXISTS User_Groups (
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, group_id),
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (group_id) REFERENCES Groups(id)
);


CREATE TABLE IF NOT EXISTS Fixtures (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_football_id INTEGER NOT NULL,
    league_season_id INTEGER,
    home_team_id INTEGER NOT NULL,
    away_team_id INTEGER NOT NULL,
    timestamp_numb INTEGER,
    venue_id INTEGER,
    status INTEGER NOT NULL,
    referee TEXT,
    home_score INTEGER,
    away_score INTEGER,
    home_winner INTEGER,
    away_winner INTEGER,
    round TEXT,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (league_season_id) REFERENCES Leagues_Seasons(id),
    FOREIGN KEY (home_team_id) REFERENCES Teams(id),
    FOREIGN KEY (away_team_id) REFERENCES Teams(id),
    FOREIGN KEY (venue_id) REFERENCES Venues(id)
);

CREATE TABLE IF NOT EXISTS Guesses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    fixture_id INTEGER NOT NULL,
    locked INTEGER NOT NULL,
    home_goals INTEGER NOT NULL,
    away_goals INTEGER NOT NULL,
    points INTEGER DEFAULT 0,
    created_at INTEGER DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER DEFAULT (strftime('%s', 'now')),
    outcome INTEGER,
    counted INTEGER DEFAULT 0,
    CHECK (outcome BETWEEN 0 AND 9)
    FOREIGN KEY (user_id) REFERENCES Users(id),
    FOREIGN KEY (group_id) REFERENCES Groups(id),
    FOREIGN KEY (fixture_id) REFERENCES Fixtures(id),
    UNIQUE (user_id, group_id, fixture_id)
);

CREATE INDEX IF NOT EXISTS idx_teams_api_football_id on Teams(api_football_id);

CREATE INDEX IF NOT EXISTS idx_leagues_api_football_id on Leagues(api_football_id);

CREATE INDEX IF NOT EXISTS idx_fixtures_api_football_id on Fixtures(api_football_id);
CREATE INDEX IF NOT EXISTS idx_fixtures_home_team_id ON Fixtures(home_team_id);
CREATE INDEX IF NOT EXISTS idx_fixtures_away_team_id ON Fixtures(away_team_id);

CREATE INDEX IF NOT EXISTS idx_user_groups_user_id ON User_Groups(user_id);
CREATE INDEX IF NOT EXISTS idx_user_groups_group_id ON User_Groups(group_id);

CREATE INDEX IF NOT EXISTS idx_guesses_user_id ON Guesses(user_id);
CREATE INDEX IF NOT EXISTS idx_guesses_group_id ON Guesses(group_id);
CREATE INDEX IF NOT EXISTS idx_guesses_fixture_id ON Guesses(fixture_id);
