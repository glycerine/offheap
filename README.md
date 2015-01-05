offheap (formerly go-offheap-hashtable)
====================

Docs:

http://godoc.org/github.com/glycerine/offheap

When GC pauses are long because you've got big hash tables in Go, you need an off-heap hash-table.


 The purpose here is to have a hash table that can work away
 from Go's Garbage Collector, to avoid long GC pause times.

 We accomplish this by writing our own Malloc() and Free() implementation
 (see malloc.go) which requests memory directly from the OS.
 The keys, the values, and the entire hash table itself is kept 
 in this off-heap storage. This storage can also optionally be backed by a memory mapped file for speedy persistence and fast startup times.

 See offheap.go for all the interesting code. Modify val_t to hold
 you values, and key_t to contain your keys. Current sample code
 for three types of keys (int64, []byte, and strings) is provided (see bytekey.go). 
 For the hashing function itself, the incredibly fast [xxhash64](https://github.com/OneOfOne/xxhash) is used to produce uint64 hashes of strings and []byte.

 Note that all your key and values should be inline in the Cell. If you
 point back into the go-heap, such values maybe garbage collected by
 the Go runtime without notice.

 Initial HashTable implementation inspired by the public domain C++ code of

 [https://github.com/preshing/CompareIntegerMaps](https://github.com/preshing/CompareIntegerMaps)

 See also

 [http://preshing.com/20130107/this-hash-table-is-faster-than-a-judy-array/](http://preshing.com/20130107/this-hash-table-is-faster-than-a-judy-array/)

 for performance studies of the C++ code.

installation
------------

     go get -t github.com/glycerine/go-offheap-hashtable

testing
--------

    go test -v


The implementation was test driven and includes over 4500 correctness checks.

Copyright (C) 2015 by Jason E. Aten, Ph.D.

License: MIT.