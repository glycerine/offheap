go-offheap-hashtable
====================

When GC pauses are long because you've got big hash tables in Go, you need an off-heap hash-table.


 The purpose here is to have hash table that can work away
 from Go's Garbage Collector, to avoid long GC pause times.

 We accomplish this by writing our own Malloc() and Free() implementation
 (see malloc.go) which requests memory directly from the OS.
 The keys, values, and entire hash table is kept on off-heap
 storage. This storage can also optionally be backed by memory mapped file
 for speedy persistence and fast startup times.

 Initial HashTable implementation inspired by the public domain C++ code of

 [https://github.com/preshing/CompareIntegerMaps](https://github.com/preshing/CompareIntegerMaps)

 See also

 [http://preshing.com/20130107/this-hash-table-is-faster-than-a-judy-array/](http://preshing.com/20130107/this-hash-table-is-faster-than-a-judy-array/)

 for performance studies of the C++ code.


Copyright (C) 2015 by Jason E. Aten, Ph.D.

License: MIT.