# PayList-API

Pay your bill easily!

## Setup 

1. Install golang and setup `$GOPATH`. Download in [here](#https://golang.org/dl/) and Installation instruction can be found in [here](#https://golang.org/doc/install)
2. `go get -u github.com/ariebrainware/paylist-api`
3. Setup `connString`  variable inside `main.go`
4. `go run main.go`

## API Design

| Endpoint         | Method | Description                                           |
| ---------------  | ------ | ----------------------------------------------------- |
| /paylist         | GET    | Show all paylist                                      |
| /paylist         | POST   | Add new paylist                                       |
| /paylist/:id     | PUT    | Update paylist based on `id` in parameter request     |
| /paylist/:id     | DELETE | Delete paylist based on input `id` parameter          |
| /paylist/:id     | GET    | Show single paylist based on id in parameter request  |
| /users           | GET    | Show all user                                         |
| /users/:id       | GET    | Show single user bases on `id` in parameter request   |
| /users/:id       | PUT    | Update paylist based on `id` in parameter request     |
| /users/:id       | DELETE | Delete user based on `id` in parameter                |
| /users/signup    | POST   | Sign Up user or create new user                       |
| /users/signin    | POST   | Sign in for user                                      |



## Database Design

Database design will automatically create by using `db.AutoMigrate`. So you just need to config database connection string inside `main.go`, then run `main.go`

```
+------------+------------------+------+-----+---------+----------------+
| Field      | Type             | Null | Key | Default | Extra          |
+------------+------------------+------+-----+---------+----------------+
| id         | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
| created_at | timestamp        | YES  |     | NULL    |                |
| updated_at | timestamp        | YES  |     | NULL    |                |
| deleted_at | timestamp        | YES  | MUL | NULL    |                |
| name       | varchar(255)     | YES  |     | NULL    |                |
| amount     | int(11)          | YES  |     | NULL    |                |
+------------+------------------+------+-----+---------+----------------+
```