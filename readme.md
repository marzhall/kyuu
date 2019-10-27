9p server that queues messages on write to the file and prints messages on read from the file

Usage:

    ./kyuu -D -v queueName &
    sudo mount -t 9p -o tcp,trans=tcp,port=5640 127.0.0.1 test
