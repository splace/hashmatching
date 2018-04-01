searches for a sequence of bytes, usually starting short and growing in length, that when appended to some source data they together have a hash with, a number of, its leading bits equal to either zero or one.

for usage see; executables/readme/usage

# examples

to search for a hash less than a value, say 2^^x, use (hash bit length)-x leading zero bits.

eg. for a SHA512 hash value less-than 2^^500, search for 12 leading zero bits.

to search for a hash greater than a value, say (2^^hash-bit-length - 2^^x), use x leading set bits.(-set option)

# characteristics

each increment of the bit count halves to chance of matching, so on average doubles the searches needed.  

uses the standard libs hashing routines, so supports what they do, see; executables/readme/usage/hash

multi-threaded, scales with cores pretty precisely, since no inter-thread comms.

sequences of bytes have a reference number, uint64, called the hash index, this can be used to start a new search where another left off, without duplication.

# performance

|cpu|hash rate SHA512|Watts|#/j|
|-|-|-|-|
|E4700, 2-Core 2.6GHz|1.45M|60|24k|
|raspberry pi3, 4-Core 1.2GHz|240k|2.5|100k|
|i7 7700K  4-core 4.2GHz|10M (guess)|90|110k|
|Threadripper 1920X 12-Core 4GHz|24M (big guess)|180|130k|
|AMD A10-9600P, 4-Core 2.3GHz|4M|15|250k|

