_default:
    @just --list

# run unit tests
test:
    go test -cover ./...

# run unit tests and produce coverage report
test-cover: mocks
    go install github.com/axw/gocov/gocov@latest
    go install github.com/AlekSi/gocov-xml@latest
    gocov test ./... | gocov-xml > coverage.xml

# run unit tests and produce coverage report (docker)
test-cover-ci:
    docker compose -f docker/docker-compose.build.yml run --rm -v $(pwd):/work -w /work go sh -c "just test-cover"

# tidy modules
tidy:
    go mod tidy

# format code
fmt:
    gofmt -w .

# build executables
build:
    @mkdir -p bin/
    go build -o bin/ ./...

# generate mocks
mocks:
    @go install github.com/vektra/mockery/v2@v2.42.0
    mockery

# bump
tag-bump semver='bump':
    git tag -a v$(convco version --{{semver}})

# changelog
changelog:
    convco changelog

