# PayList-API

Pay your bill easily!

## Setup 

1. Install golang and setup `$GOPATH`. Download in [here](https://golang.org/dl/) and Installation instruction can be found in [here](https://golang.org/doc/install)
2. `go get -u github.com/ariebrainware/paylist-api`
3. Setup `connString`  variable inside `config.go`
4. `go run main.go`
5. Install go get -u github.com/pressly/goose/cmd/goose for goose migration or you can find Installation instruction in [here](https://github.com/pressly/goose)

## API Design

| Endpoint              | Method | Description                                           |
| ---------------       | ------ | ----------------------------------------------------- |
| /paylist              | GET    | Show all user-paylist data                            |
| /paylist              | POST   | Add new user-paylist data                             |
| /paylist/:id          | PUT    | Update paylist based on `id` in parameter request     |
| /paylist/:id          | DELETE | Delete user-paylist based on input `id` parameter     |
| /paylist/:id          | GET    | Show single paylist based on id in parameter request  |
| /users                | GET    | Show all user                                         |
| /users/:id            | GET    | Show single user bases on `id` in parameter request   |
| /users/:id            | PUT    | Update paylist based on `id` in parameter request     |
| /users/:id            | DELETE | Delete user based on `id` in parameter                |
| /users/signup         | POST   | Sign Up user or create new user                       |
| /users/signin         | POST   | Sign in for user                                      |
| /user/signout         | GET    | Sign out user or logout                               |
| /user-paylist/:id     | PUT    | Update user-paylist status(complete or not)           |
| /users/refresh-token  | POST   | Refresh Expired Token                                 |


## Database Design

Use GORM

Database design will automatically create by using `db.AutoMigrate`. So you just need to config database connection string inside `config.go`, then run `main.go`

Goose Migration

Create a new SQL migration.
```
goose create add_some_column sql
```

Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND
```
goose mysql "user:password@/dbname?parseTime=true" status
```

Commands:
```
up                   Migrate the DB to the most recent version available
down                 Roll back the version by 1
status               Dump the migration status for the current DB
create NAME [sql|go] Creates new migration file with the current timestamp
```

```
Paylists Table
+------------+------------------+------+-----+---------+----------------+
| Field      | Type             | Null | Key | Default | Extra          |
+------------+------------------+------+-----+---------+----------------+
| id         | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| created_at | timestamp        | YES  |     | NULL    |                |
| updated_at | timestamp        | YES  |     | NULL    |                |
| deleted_at | timestamp        | YES  | MUL | NULL    |                |
| name       | varchar(255)     | YES  |     | NULL    |                |
| amount     | int(11)          | YES  |     | NULL    |                |
| username   | varchar(255)     | YES  |     | NULL    |                |
| completed  | tinyint(1)       | YES  |     | NULL    |                |
+------------+------------------+------+-----+---------+----------------+
```

```
Users Table
+------------+------------------+------+-----+---------+----------------+
| Field      | Type             | Null | Key | Default | Extra          |
+------------+------------------+------+-----+---------+----------------+
| id         | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| created_at | timestamp        | YES  |     | NULL    |                |
| updated_at | timestamp        | YES  |     | NULL    |                |
| deleted_at | timestamp        | YES  | MUL | NULL    |                |
| email      | varchar(255)     | YES  |     | NULL    |                |
| name       | varchar(255)     | YES  |     | NULL    |                |
| username   | varchar(255)     | NO   | PRI | NULL    |                |
| password   | varchar(255)     | YES  |     | NULL    |                |
| balance    | int(11)          | YES  |     | NULL    |                |
+------------+------------------+------+-----+---------+----------------+
```

```
Loggings Table
+------------+------------------+------+-----+---------+----------------+
| Field      | Type             | Null | Key | Default | Extra          |
+------------+------------------+------+-----+---------+----------------+
| username   | varchar(255)     | YES  |     | NULL    |                |
| token      | varchar(255)     | YES  |     | NULL    |                |
| user_status| tinyint(1)       | YES  |     | NULL    |                |
| created_at | timestamp        | YES  |     | NULL    |                |
| deleted_at | timestamp        | YES  | MUL | NULL    |                |
+------------+------------------+------+-----+---------+----------------+
```