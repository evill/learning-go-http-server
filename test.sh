#!/bin/sh

go clean -testcache && go test -count=1 ./e2e