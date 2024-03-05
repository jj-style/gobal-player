_default:
    @just --list

# run unit tests
test:
    go list -f '{{{{.Dir}}' -m | xargs go test -cover

# run unit tests and produce coverage report
test-cover: mocks
    go install github.com/axw/gocov/gocov@latest
    go install github.com/AlekSi/gocov-xml@latest
    go list -f '{{{{.Dir}}' -m | xargs gocov test | gocov-xml > coverage.xml

# run unit tests and produce coverage report (docker)
test-cover-ci:
    docker compose -f docker/docker-compose.build.yml run --rm -v $(pwd):/work -w /work go sh -c "just test-cover"

# tidy modules
tidy:
    go list -f '{{{{.Dir}}' -m | xargs -L1 go mod tidy -C

# format code
fmt:
    gofmt -w .

# generate mocks
mocks:
    @go install github.com/vektra/mockery/v2@v2.42.0
    mockery