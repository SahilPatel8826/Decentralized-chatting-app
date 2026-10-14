# Running locally (Day 1)

1. Start Postgres and create the database (`createdb chat_db`).
2. `cp user_service/.env.example user_service/.env` and
   `cp chat_service/.env.example chat_service/.env`. Edit the values.
   `JWT_ACCESS_SECRET` must be the same in both files.
3. Terminal 1: `cd user_service && go run ./cmd/server`  (HTTP :8080, gRPC :50051)
4. Terminal 2: `cd chat_service && go run ./cmd/server`  (HTTP/WS :8081)
5. Smoke test: `bash scripts/smoke_test.sh`
6. Live check with two terminals:
   `wscat -c "ws://localhost:8081/ws?token=<TOKEN_B>"`
   then send a message as user A (see the script output) and watch it arrive.
