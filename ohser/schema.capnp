@0x8b53f943fbc3527f;
using Go = import "go.capnp";
$Go.package("offheap");
$Go.import("testpkg");

struct CellCapn { 
   unHashedKey  @0:   UInt64;
   byteKey      @1:   Data;
   value        @2:   Data;
} 

struct HashTableCapn { 
   cells       @0:   UInt64;
   cellSz      @1:   UInt64; 
   arraySize   @2:   UInt64; 
   population  @3:   UInt64; 
   zeroUsed    @4:   Bool; 
   zeroCell    @5:   CellCapn; 
   offheap     @6:   Data;
   mmm         @7:   MmapMallocCapn; 
} 

struct IteratorCapn { 
   tab  @0:   HashTableCapn; 
   pos  @1:   Int64; 
   cur  @2:   CellCapn; 
} 

struct MmapMallocCapn { 
   path          @0:   Text; 
   fd            @1:   Int64; 
   fileBytesLen  @2:   Int64; 
   bytesAlloc    @3:   Int64; 
   mem           @4:   Data;
} 

##compile with:

##
##
##   capnp compile -ogo odir/schema.capnp

