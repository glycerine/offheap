/*


Offheap

An off-heap hash-table for Go (golang). Originally called go-offheap-hashtable,
but now shortened to just offheap.

The purpose here is to have a hash table that can work away
from Go's Garbage Collector, to avoid long GC pause times.

We accomplish this by writing our own Malloc() and Free() implementation
(see malloc.go) which requests memory directly from the OS.

The keys, values, and entire hash table is kept on off-heap
storage. This storage can also optionally be backed by memory mapped file
for speedy persistence and fast startup times.

Initial HashTable implementation inspired by the public domain C++ code of
    https://github.com/preshing/CompareIntegerMaps
See also
    http://preshing.com/20130107/this-hash-table-is-faster-than-a-judy-array/
for performance studies of the C++ code.


HashTable

The implementation is mostly in offheap.go, read that to start.

Maps pointer-sized integers to Cell structures, which in turn hold Val_t
as well as Key_t structures.

Uses open addressing with linear probing. This makes it very cache
friendly and thus very fast.

In the t.Cells array, UnHashedKey = 0 is reserved to indicate an unused cell.
Actual value for key 0 (if any) is stored in t.ZeroCell.
The hash table automatically doubles in size when it becomes 75% full.
The hash table never shrinks in size, even after Clear(), unless you explicitly
call Compact().

Basic operations: Lookup(), Insert(), DeleteKey(). These are the
equivalent of the builtin map[uint64]interface{}.

As an example of how to specialize for a map[string]*Cell equivalent,
see the following functions in the bytekey.go file:

    func (t *HashTable) InsertStringKey(strkey string, value interface{}) bool
    func (t *HashTable) LookupStringKey(strkey string) (Val_t, bool)
    func (t *HashTable) DeleteStringKey(strkey string) bool


Example use:

    h := offheap.NewHashTable(2)

    // basic three operations are:
    h.InsertStringKey("My number", 43)
    val, ok := h.LookupStringKey("My number")
    h.DeleteStringKey("My number")

Note that this library is only a starting point of source code, and not intended to be used without customization. Users of the HashTable will have to customize it by changing the definitions of Key_t and Val_t to suite their needs.

On Save(), serialization of the HashTable itself is done using msgpack to write bytes to the first page (4k bytes) of the memory mapped file. This uses github.com/tinylib/msgp which is a blazing fast msgpack serialization library. It is fast because it avoids reflection and pre-computes the serializations (using go generate based inspection of your go source). If you need to serialize your values into the Val_t, I would suggest evaluating the msgp for serialization and deserialization. The author, Philip Hofer, has done a terrific job and put alot of effort into tuning it for performance. If you are still pressed for speed, consider also omitting the field labels using the '//msgp:tuple MyValueType' annotation. As Mr. Hofer says, "For smaller objects, tuple encoding can yield serious performance improvements." [https://github.com/tinylib/msgp/wiki/Preprocessor-Directives].

Related ideas:

https://gist.github.com/mish15/9822474 (using CGO)

CGO note: the cgo-malloc branch of this github repo has an implementation that uses CGO to
call the malloc/calloc/free functions in the C stdlib. Using CGO
gives up the save-to-disk instantly feature and creates a portability issue where
you have linked against a specific version of the C stdlib. However if you
are making/destroying alot of tables, the CGO approach may be faster. This
is because calling malloc and free in the standard C library are much faster than
making repeated system calls to mmap().

more related ideas:

https://groups.google.com/forum/#!topic/golang-nuts/kCQP6S6ZGh0

not fully off-heap, but using a slice instead of a map appears to help GC quite alot too:

https://github.com/cespare/kvcache/blob/master/refmap.go


*/
package offheap
