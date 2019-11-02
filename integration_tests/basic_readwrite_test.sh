#!/bin/bash

sh setup_tests.sh
touch test_kyuus/test_queue

firstMessage="firstMessage"
secondMessage="secondMessage"

# testing
echo $firstMessage > test_kyuus/test_queue
echo $secondMessage > test_kyuus/test_queue

msg1=`cat test_kyuus/test_queue`
if [ $msg1 != $firstMessage ]
then
    echo "Message 1 was not produced correctly by the test_queue."
    echo $msg1
    sh teardown_test.sh
    exit 1
else
    echo "Test 1 passed."
fi

msg2=`cat test_kyuus/test_queue`
if [ $msg2 != $secondMessage ]
then
    echo "Message 2 was not produced correctly by the test_queue."
    echo $msg2
    sh teardown_tests.sh
    exit 1
else
    echo "Test 2 passed."
fi

#teardown
sh teardown_tests.sh
