@0xb2ce27958b10688c;
using Go = import "go.capnp";
$Go.package("keyval");
$Go.import("testpkg");


# the value payload
struct AccountCapn { 
   id            @0:   Int64; 
   dty           @1:   Int64; 
   acctId        @2:   Text; 
   openedFromIP  @3:   Text; 
   name          @4:   Text; 
   email         @5:   Text; 
   disabled      @6:   Int64; 
} 
