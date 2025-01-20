# Real-time chat server

Hi! It's a real-time chat server, at the same time a large number of users can read and send messages, by using websocket technology, making two-way communication between client and server over TCP protocol. All users connections, before broadcast, upgrade from http connection to ws in real-time.

### Go | Gorilla/Websocket | Redis

- Websocket protocol (powered by Gorilla framework).
- NoSQL database Redis.
- Logging users and messages by ID.
- Catching runtime errors.

### Execute app (mac)
You have to install:
- Go (https://go.dev/doc/install)
- Redis v8 (https://github.com/redis/go-redis)
- Gorilla/Websocket (https://github.com/gorilla/websocket)
 Make sure, that you run "redis-server", after launch application, you can check connection with Redis server with command in terminal "netstat -an | grep 6379".

##### `1 redis-cli ping (answer: PONG)`
##### `2 redis-server`
##### `3 go run ./cmd/web`
##### `4 netstat -an | grep 6379`
##### `5 http://localhost:8081/`

<b>For test application you can use tools like a Postman or websocat, but also your browser</b>