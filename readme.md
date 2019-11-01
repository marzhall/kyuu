# Kyuu

## 9p-based virtual file queue

Kyuu creates message queues that expose themselves as files.

Writes to a queue file place a message on it, and reads from
the file return a message.

The queue server exposes itself to the network at an IP address
and port, and can be mounted by the built-in linux 9p driver, or,
more easily, the 9pfuse FUSE filesystem driver.

The queue can be mounted, read from, and written to simulatenously on
multiple places on the filesystem and from multiple machines.

# Usage

Starting a queue daemon, writing to a queue, then popping messages from
the queue by `cat`ting the file:

        $: kyuu myqueue &
        $: 9pfuse 127.0.0.1:5640 test

        $: echo msg1 > test/myqueue
        $: echo msg2 > test/myqueue

        $: cat test/myqueue
           msg1
        $: cat test/myqueue
           msg2

Creating a new queue, after having mounted the filesystem, is done by
just touching a new file in the kyuu folder:

        $: touch test/newqueue
        $: ls test
           myqueue
           newqueue


## Dependency

The styx 9p server for go library:
    https://godoc.org/aqwari.net/net/styx

## Optional Dependency

For easier mounting, the 9pfuse package from plan9port is very helpful:
    https://9fans.github.io/plan9port/

There is a 9p driver in the linux kernel that will allow the mounting
and use of kyuu without building 9pfuse, but it requires root and is
a pain in the butt.

The linux 9 driver would be:

    sudo mount -t 9p -o tcp,trans=tcp,port=5640 127.0.0.1 test
