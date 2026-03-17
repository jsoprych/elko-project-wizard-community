package db

const schema = `
CREATE TABLE IF NOT EXISTS directives (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    category    TEXT NOT NULL,
    description TEXT,
    content     TEXT NOT NULL,
    tags        TEXT DEFAULT '[]',
    variables   TEXT DEFAULT '{}',
    builtin     INTEGER DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profiles (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    description  TEXT DEFAULT '',
    project_name TEXT DEFAULT '',
    variables    TEXT DEFAULT '{}',
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS profile_directives (
    profile_id   TEXT NOT NULL,
    directive_id TEXT NOT NULL,
    sort_order   INTEGER DEFAULT 0,
    PRIMARY KEY (profile_id, directive_id),
    FOREIGN KEY (profile_id)   REFERENCES profiles(id)   ON DELETE CASCADE,
    FOREIGN KEY (directive_id) REFERENCES directives(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS generations (
    id           TEXT PRIMARY KEY,
    profile_id   TEXT,
    project_name TEXT NOT NULL,
    files        TEXT DEFAULT '[]',
    created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
);
`
