/*

TODO : update

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

Note that this library is only a starting point of source code, and not intended to be used without customization. Users of the HashTable will have to customize it by changing the definitions of Key_t and Val_t to suite their needs. I'm experimenting next with storing objects in Capnproto serialized format, but this branch (branch capnp) isn't quite ready for use.

Related ideas:

https://gist.github.com/mish15/9822474 (using CGO)

https://groups.google.com/forum/#!topic/golang-nuts/kCQP6S6ZGh0

not fully off-heap, but using a slice instead of a map appears to help GC quite alot too:

https://github.com/cespare/kvcache/blob/master/refmap.go


*/
package offheap
