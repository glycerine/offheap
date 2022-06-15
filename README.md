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

Update: the cgo-malloc branch of this github repo has an implementation that uses CGO to
call the malloc/calloc/free functions in the C stdlib. Using CGO 
gives up the save-to-disk instantly feature and creates a portability issue where
you have linked against a specific version of the C stdlib. However if you
are making/destroying alot of tables, the CGO approach may be faster.

 See offheap.go for all the interesting code. Modify val_t to hold
 you values, and key_t to contain your keys. Current sample code
 for three types of keys (int64, []byte, and strings) is provided (see bytekey.go). 
 For the hashing function itself, the incredibly fast [xxhash64](https://github.com/OneOfOne/xxhash) is used to produce uint64 hashes of strings and []byte.

 Note that all your key and values should be inline in the Cell. If you
 point back into the go-heap, such values maybe garbage collected by
 the Go runtime without notice.

 On Save(), serialization of the HashTable itself is done using msgpack to write bytes to the first page (4k bytes) of the memory mapped file. This uses github.com/tinylib/msgp which is a blazing fast msgpack serialization library. It is fast because it avoids reflection and pre-computes the serializations (using go generate based inspection of your go source). If you need to serialize your values into the Val_t, I would suggest evaluating the msgp for serialization and deserialization. The author, Philip Hofer, has done a terrific job and put alot of effort into tuning it for performance. If you are still pressed for speed, consider also omitting the field labels using the '//msgp:tuple MyValueType' annotation. As Mr. Hofer says, "For smaller objects, tuple encoding can yield serious performance improvements." [https://github.com/tinylib/msgp/wiki/Preprocessor-Directives].

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