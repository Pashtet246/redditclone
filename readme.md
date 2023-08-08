Reddit Clone

Project from the course: Backend for a Reddit clone (to be used with the provided frontend - https://asperitas.now.sh/).

Store sessions in a K/V database such as redis/memcached/tarantool (redis is commonly used). In JWT sessions, stateless mode can no longer be used. Now they will be stateful (stored in the database), meaning that when a session arrives, you go to the database and check if there is a record with that ID.
Store users in mysql/postgres.
Store posts with comments in MongoDB (use github.com/mongodb/mongo-go-driver).

This assignment will allow you to work with fundamental aspects used in application development:

Designing simple tables in MySQL.
Writing basic queries for MySQL.
Testing the repository pattern with MySQL and MongoDB.
Testing HTTP handlers that are closer to reality, not fake ones like in the 4th assignment.
You can implement this through separate structures that satisfy the specified interface. For example, alongside another repository, you could have session_jwt.go + session_mysql.go + session_redis.go.

Requirements:
go v16 or higher
docker

How to run:
1. docker-compose up -d && sleep 10 && cd cmd/api && go run main.go
2. Execute init.sql in database/mysql
