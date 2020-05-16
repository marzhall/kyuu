#!/bin/bash

sh setup_tests.sh

echo
echo "##############################################"
echo RUNNING PERSISTANCE TESTS
echo "##############################################"


touch test_kyuus/test_queue

firstMessage="firstMessage"

# testing
echo $firstMessage > test_kyuus/test_queue
num_lines=$(cat kyuu_storage/queued | wc -l)
if [ $num_lines != "1" ]
then
	echo "FAILURE: No message was stored to the persistance file when a message was added to the test queue."
	sh teardown_test.sh
	exit 1
else
	echo "SUCCESS: Message was written to the persistance file when added to the queue."
fi
	
cat test_kyuus/test_queue
num_lines=$(cat kyuu_storage/dequeued | wc -l)
if [ $num_lines != "1" ]
then
    echo "FAILURE: A message was not added to the 'dequeued' file after being dequeued."
    echo $msg1
    sh teardown_test.sh
    exit 1
else
    echo "SUCCESS: A message was added to the 'dequeued' file after being dequeued."
fi

undeliveredMessages=`diff kyuu_storage/queued kyuu_storage/dequeued`
if [ $undeliveredMessages ]
then
    echo "FAILURE: After queueing and dequeueing a message, the persistance files disagree on what messages have been sent/received."
    sh teardown_tests.sh
    exit 1
else
    echo "SUCCESS: Persitance files agree on what messages were received and sent."
fi

#teardown
sh teardown_tests.sh
