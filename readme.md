9p server that queues messages on write to the file and prints messages on read from the file

Usage:

    # Using the linux kernel 9p driver
    ./kyuu myqueue &
    sudo mount -t 9p -o tcp,trans=tcp,port=5640 127.0.0.1 test
    echo foo > test/myqueue
    cat test/myqueue
    # responds: "foo"
    # When stopping:
    umount test

    # Using 9pfuse from the Plan 9 in Userspace toolset:
    sh start.sh
    echo foo > test/myqueue
    cat test/myqueue
    # responds: "foo"
    fusermount -u test
