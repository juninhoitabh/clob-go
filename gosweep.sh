#!/bin/zsh
set -e

echo 'mode: count' > coverage.out

max_steps=3

for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path './.history*' -path '*/internal/*' -type d);
do
if ls $dir/*.go &> /dev/null; then
    if [[ $dir != (*"dtos"*|*"fakers"*|*"infra/http-server"*|*"middle-wares"*|*"mocks"*) ]]; then
        echo "Running go test race $dir ... (1/$max_steps)"
        ENVIRONMENT=test go test -tags=all -v -short -race "$dir"

        echo "Running go test $dir ... (2/$max_steps)"
        ENVIRONMENT=test go test -tags=all -v -short -covermode=count -coverprofile=$dir/coverage.tmp "$dir"
        if [ -f $dir/coverage.tmp ]
        then
            cat $dir/coverage.tmp | tail -n +2 >> coverage.out
            rm $dir/coverage.tmp
        fi
    fi
fi
done

go tool cover -func coverage.out
go tool cover -html=coverage.out -o coverage.html

echo "Running golangci-lint $dir ... (3/$max_steps)"

golangci-lint run ./...