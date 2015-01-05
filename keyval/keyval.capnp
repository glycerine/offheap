@0x887a749f8d712028;
using Go = import "go.capnp";
$Go.package("keyval");
$Go.import("testpkg");


struct AccountCapn { 
   id            @0:   Int64; 
   dty           @1:   Int64; 
   acctId        @2:   Text; 
   openedFromIP  @3:   Text; 
   name          @4:   Text; 
   email         @5:   Text; 
   disabled      @6:   Int64; 
} 

struct AcctIdCapn { 
   acctId  @0:   Text; 
} 

##compile with:

##
##
##   capnp compile -ogo ./schema.capnp

