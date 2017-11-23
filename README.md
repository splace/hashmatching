# hashmatching
searches through sequences of bytes, usually growing in length, that when appended to some source data they together have a hash with, a number of, its leading bits equal to either zero or one.

for usage see; executables/readme/usage

to search for a hash less than a value, say 2^^x, use (hash bit length)-x leading zero bits.

so for a SHA512 hash value less-than 2^^500, search for 12 leading zero bits.

to search for a hash greater than a value, say (2^^hash-bit-length - 2^^x), use x leading set bits.(-set option)

each increment of the bit count halves to chance of matching, so on average doubles the searches needed.  

uses the standard libs hashing routines, so supports what they do, see; executables/readme/usage/hash

multi-threaded, scales with cores pretty precisely, since no inter-thread comms.

sequences of bytes have a reference number, uint64, called the hash index, this can be used to start a new search where another left off, without duplication.


|cpu|hash rate SHA512|watts|#/j|
|-|-|-|-|
|intel core2 2.6GHz|1.2M|60|20k|
|raspberry pi3, 4core 1.2GHz|200k|2.5|80k|
|inter i7 7700K  4.2GHz|8M (guess)|90|110k|
|Threadripper 1920X 4GHz|20M (big guess)|180|110k|

