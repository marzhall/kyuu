#!/bin/bash
# setup
sh start.sh

firstMessage="firstMessage"
secondMessage="secondMessage"

# testing
echo $firstMessage > test/queue
echo $secondMessage > test/queue

msg1=`cat test/queue`
if [ $msg1 != $firstMessage ]
then
    echo "Message 1 was not produced correctly by the queue."
    echo $msg1
    sh stop.sh
    exit 1
else
    echo "Test 1 passed."
fi

msg2=`cat test/queue`
if [ $msg2 != $secondMessage ]
then
    echo "Message 2 was not produced correctly by the queue."
    echo $msg2
    sh stop.sh
    exit 1
else
    echo "Test 2 passed."
fi

#teardown
echo "Tearing down."
sh stop.sh
