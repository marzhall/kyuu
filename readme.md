# Kyuu

## 9 queue

Kyuu is a message queue that exposes itself for mount onto a filesystem.

Once mounted as a file, writes to the queue file place a message on it,
and reads from the file produce a message.

The queue can be mounted, read from, and written to simulatenously on
multiple places on the filesystem and from multiple machines.

# Usage

        $: kyuu myqueue &
        $: 9pfuse 127.0.0.1:5640 test

        $: echo msg1 > test/myqueue
        $: echo msg2 > test/myqueue

        $: cat test/myqueue
           msg1
        $: cat test/myqueue
           msg2

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
