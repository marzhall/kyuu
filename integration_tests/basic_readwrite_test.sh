#!/bin/bash

# setup
sh start.sh
touch test/test_queue

firstMessage="firstMessage"
secondMessage="secondMessage"

# testing
echo $firstMessage > test/test_queue
echo $secondMessage > test/test_queue

msg1=`cat test/test_queue`
if [ $msg1 != $firstMessage ]
then
    echo "Message 1 was not produced correctly by the test_queue."
    echo $msg1
    sh stop.sh
    exit 1
else
    echo "Test 1 passed."
fi

msg2=`cat test/test_queue`
if [ $msg2 != $secondMessage ]
then
    echo "Message 2 was not produced correctly by the test_queue."
    echo $msg2
    sh stop.sh
    exit 1
else
    echo "Test 2 passed."
fi

#teardown
echo "Tearing down."
sh stop.sh
