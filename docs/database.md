# Database Setup

This application uses MySQL.

Choose one of the installation methods below.

---

## Option 1: Docker (Recommended)

### Start MySQL

```bash
docker run --name snippetbox-db \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=your_root_password \
  -v snippetbox_data:/var/lib/mysql \
  -d mysql:9.7.0
```

### Connect to MySQL

```bash
docker exec -it snippetbox-db mysql -u root -p
```

---

## Option 2: macOS (Homebrew)

Install MySQL:

```bash
brew install mysql
```

Start MySQL:

```bash
brew services start mysql
```

Connect:

```bash
mysql -u root -p
```

---

## Option 3: Ubuntu / Debian

Install MySQL:

```bash
sudo apt install mysql-server
```

Connect:

```bash
sudo mysql
```

Or:

```bash
mysql -u root -p
```

---

## Create the Database

```sql
CREATE DATABASE snippetbox
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE snippetbox;
```

---

## Create the Table

```sql
CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created
    ON snippets(created);

CREATE TABLE sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
```

---

## Seed Data

```sql
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);
```

---

## Create the Application User

### Local MySQL Installation

```sql
CREATE USER IF NOT EXISTS 'web'@'%'
IDENTIFIED BY 'pass';

GRANT SELECT, INSERT, UPDATE, DELETE
ON snippetbox.*
TO 'web'@'%';

FLUSH PRIVILEGES;
```

### Docker Installation

Because the application connects from outside the container:

```sql
CREATE USER IF NOT EXISTS 'web'@'%'
IDENTIFIED BY 'pass';

GRANT SELECT, INSERT, UPDATE, DELETE
ON snippetbox.*
TO 'web'@'%';

FLUSH PRIVILEGES;
```

---

## Verify the User

```bash
mysql -D snippetbox -u web -p
```

Run:

```sql
SELECT id, title, expires
FROM snippets;
```

Expected output:

```text
+----+------------------------+---------------------+
| id | title                  | expires             |
+----+------------------------+---------------------+
|  1 | An old silent pond     | 2027-xx-xx xx:xx:xx |
|  2 | Over the wintry forest | 2027-xx-xx xx:xx:xx |
|  3 | First autumn morning   | 2026-xx-xx xx:xx:xx |
+----+------------------------+---------------------+
```

---

## DataGrip Configuration

| Setting  | Value      |
| -------- | ---------- |
| Host     | localhost  |
| Port     | 3306       |
| User     | web        |
| Password | pass       |
| Database | snippetbox |

---

## WebStorm Configuration

| Setting  | Value      |
| -------- | ---------- |
| Host     | localhost  |
| Port     | 3306       |
| User     | web        |
| Password | pass       |
| Database | snippetbox |

---

## Test Database Setup

Integration tests in `internal/models` use a separate test database and user. Each test run creates and drops tables via `testdata/setup.sql` and `testdata/teardown.sql`, so the test user needs extra DDL privileges.

### Create the Test Database and User

Connect as root, then run:

```sql
CREATE DATABASE IF NOT EXISTS test_snippetbox
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'test_web'@'%' IDENTIFIED BY 'pass';

GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE, DELETE
ON test_snippetbox.*
TO 'test_web'@'%';

FLUSH PRIVILEGES;
```
