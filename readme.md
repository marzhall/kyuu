# Kyuu

## Message queue using 9p virtual filesystems

Kyuu creates message queues that expose themselves as files.

Writing to a queue file places a message on it, and reading
from the file returns a message.

The kyuu daemon opens a socket, allowing it to be mounted
throughout the network. It can be mounted by any machine running
a 9p fileystem driver, such as the built-in linux 9p driver, or,
more easily (and without root priveleges), the 9pfuse FUSE filesystem
driver.

Once mounted, a queue can be read from and be written to simulatenously
at an arbitrary number of mount points, and entirely using the normal
filesystem operations of `READ` and `WRITE`, such as through using `cat` and
`echo message > queue` in bash.

When the kyuu daemon is spun up by the `kyuu` client, by default, its
virtual filesystem will be mounted in the directory `$KYUUPATH`; if
`$KYUUPATH` is unset, it will create a new directory `$HOME/kyuus` and
mount there.

# Usage

Building this package creates a server binary `kyuuD`, and a client tool `kyuu`.

Creating, writing to, and reading from queues is handled by touching,
reading, and writing files to wherever the kyuu daemon is mounted (usually,
`$KYUUPATH` or `$HOME/kyuus`). Touching a file in `$KYUUPATH`, e.g.:

        $: touch $KYUUPATH/newqueue
        $: ls $KYUUPATH
        >  newqueue
        
creates a new file as expected - but writing and reading the file acts like
a message queue:

        $: cd $KYUUPATH
        $: echo message1 > newqueue
        $: echo message2 > newqueue
        $: cat newqueue
        > message1
        $: cat newqueue
        > message2
        $: cat newqueue
        > 

Note that last bit - if the queue is empty, you just get an `EOF`!

Likewise, deleting a queue is as might be expected:

        $: rm newqueue

Directories can also be created to house related queues.

The `kyuu` client is a bash script that glues together the above file
operations in convenient ways, as well as automatically starting the
daemon and dynamically linking queue files to the current directory for 
ease-of-use, but it does nothing you can't do with a `bash` prompt.

## "Everyday Usage" Example

Okay, you've got a need for a new queue in your current directory. You run:

        $: kyuu testqueue
        $: ls
        > testqueue

The new file in your directory is a message queue named "testqueue" -
you can add a messages to the queue by writing to it:

        $: echo message1 > testqueue
        $: echo message2 > testqueue

And get those messages by reading from it:

        $: cat testqueue
        > message1
        $: cat testqueue
        > message2
        $: cat testqueue
        > 

Delivery is guaranteed once-only; queues are entirely in-memory and
non-persistent (for now).

## A full example of starting and mounting the daemon and its virtual directory

As opposed to using the client, this is an example of starting a kyuu
daemon, mounting it to the folder "kyuus," creating a new queue, and
then controlling the new queue directly. The only difference between this
and using the `kyuu` client script is that we choose where to mount the kyuu
directory, and we don't dynamically link a queue file into the current
directory as the kyuu client would normally do.

        $: ls
        > kyuus

        # start and mount the daemon; default port is 5640.
        $: kyuuD&
        $: 9pfuse 127.0.0.1:5640 kyuus

        # create the new queue by `touch`ing it
        $: touch kyuus/newqueue

        $: echo message1 > kyuus/newqueue
        $: echo message2 > kyuus/newqueue
        $: cat kyuus/newqueue
        > message1
        $: cat kyuus/newqueue
        > message2
        $: cat kyuus/newqueue
        >

The `integration_tests/basic_read_write_tests.sh` bash script is also
a good example.

# Dependencies

Kyuu is built and tested with go `1.13.3`.

It requires the styx 9p server library:
    https://godoc.org/aqwari.net/net/styx

## Secondary Dependency

For easier mounting, the 9pfuse package from plan9port is very helpful:

    https://9fans.github.io/plan9port/

The "kyuu" script depends upon it, but the kyuu daemon can be mounted
without it by using the 9p driver in the Linux kernel - notably, this
requires root prileges.

The linux 9 driver command to mount the kyu directory would be:

    sudo mount -t 9p -o tcp,trans=tcp,port=5640 127.0.0.1 $KYUUPATH

# Building

With tests:

    ./build.sh tests

Without tests:

    ./build.sh
