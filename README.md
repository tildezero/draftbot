# draftbot

it's a bot for [draftout](https://draftoutmc.com) stats!

## commands
- `/stats player type=competitive|quick-play|lobby` -> gets stats for player on mode type
- more soon

## running

### docker
- `cp draftbot.env.example draftbot.env` and fill in
- `docker compose up -f docker-compose-local.yml -d` (builds docker image)
- `docker compose up -d` (uses remote docker image off of main)

### local
- `cp draftbot.env.example .env` and fill in
- `go mod download`
- `go run .`