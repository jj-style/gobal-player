_default:
    @just --list

# run unit tests
test:
    go test -cover `go list ./... | grep -v mocks`

# run unit tests and produce coverage report
test-cover: mocks
    go install github.com/axw/gocov/gocov@latest
    go install github.com/AlekSi/gocov-xml@latest
    gocov test `go list ./... | grep -v mocks` | gocov-xml > coverage.xml

# run unit tests and produce coverage report (docker)
test-cover-ci:
    docker compose -f docker/docker-compose.build.yml run --rm -v $(pwd):/work -w /work go sh -c "just test-cover"

# tidy modules
tidy:
    go mod tidy

# format code
fmt:
    gofmt -w .

# generate compile time code
generate:
    go install github.com/google/wire/cmd/wire@latest
    go generate ./...

# build executables
build:
    @mkdir -p bin/
    go build -o bin/ ./...

# generate mocks
mocks:
    @go install github.com/vektra/mockery/v2@v2.42.0
    mockery

# release
release semver='bump':
    #!/bin/bash
    currentVersion=$(convco version)
    nextVersion=$(convco version --{{semver}})

    convco changelog -u $nextVersion > CHANGELOG.md
    sed -i "s/@v$currentVersion/@v$nextVersion/g" README.md

    git add README.md CHANGELOG.md
    convco commit --chore --message "release $currentVersion -> $nextVersion"

    git tag v$nextVersion -m "v$nextVersion"

    git push --follow-tags


# changelog
changelog:
    convco changelog

# initialize git hooks
hooks:
    go install -v github.com/go-critic/go-critic/cmd/gocritic@latest
    command -v pre-commit
    pre-commit install
