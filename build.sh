#!/bin/bash
go build -o bin/kyuuD ./src
echo "Build complete."

if [ $# -eq 0 ]
then
    exit 0
fi

if [ $1 == "tests" ]
then
    echo "Testing."
    cd integration_tests
    sh basic_readwrite_test.sh
    echo "Testing complete."
fi
