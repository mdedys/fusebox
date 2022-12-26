@up:
  docker-compose up --build -d

@down:
  docker-compose down -v

@build-relay:
  docker build --file ./cmd/relay/Dockerfile -t fuse/relay:latest .

@install-fkit: 
  go build -o bin/fkit ./cmd/fkit/fkit.go
