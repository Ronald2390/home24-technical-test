# home24-technical-test

## Requirement

1. Install and run Redis in local with port 6789 (or you can adjust the config refer to the redis). To install redis: sudo apt install redis.
2. Install and run PostgreSQL in local with port 5432 (or you can adjust the config refer to the postgresql), to simplify the process, postgresql can be run by docker by exec this command "docker run --name postgres-docker -e POSTGRES_HOST_AUTH_METHOD=trust -p 5432:5432 -d postgres
". This would trigger docker to install postgres in container and run the default database "postgres" with no password.

## How to run 

### To run the server:
`go run cmd/main.go`

Listen to port 8089 by default

### To access the html:
browse the html from specific path in browser or just double click or open the html in browser 

## HTML
HTML file located in internal/view/html. there is a login.html page in there which can be opened with browser.

## API

### Login
- [POST] 127.0.0.1:8089/v1/login 
- Request Body:
{
    "email":"user@home24.com",
    "password":"user"
}
- curl command: 
curl -X POST 127.0.0.1:8089/v1/login --data $'{"email":"user@home24.com","password":"user"}'

### Logout
- [POST] 127.0.0.1:8089/v1/logout (no need request body and url params)
- curl command:
curl -X POST 127.0.0.1:8089/v1/logout -H "Authorization:session {token_retrieved_on_login}"

### Get Login Session
- [GET] 127.0.0.1:8089/v1/session 
- curl command:
curl -X GET 127.0.0.1:8089/v1/session -H "Authorization:session {token_retrieved_on_login}"

### Change Password
- [PUT] 127.0.0.1:8089/v1/users/password
- Request Body
{
    "oldPassword": "user",
    "newPassword": "newPassword"
}
- curl command: 
curl -X PUT 127.0.0.1:8089/v1/users/password -H "Authorization:session {token_retrieved_on_login}" --data $'{"oldPassword": "user","newPassword": "newPassword"}'

## Notes

- Default user password is "user"
- Exposing Port, DB Connection String, and Redis Credential can be change in config/config.go
- I put the default env variable setup in main function, it can be replace by actual environment variable (export ENVIROMENT=development) and by implementing that, the hard coded setup can be removed
- For the HTML, I am just provide the event to do login.
- I apologize for not bring the good UI for the HTML, I too focused on the backend side while working on this.