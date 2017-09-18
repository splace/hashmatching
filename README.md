# hashmatching
searchs for appended trailing bytes, the 'nonce', for some data so that its hash has at least the specified number equal leading bits.

for usage see; executables/readme/usage

to search for a hash less than a value, say 2^^x, use (hash bit length)-x leading zero bits.

for a SHA512 hash less-than 2^^500, search for 12 leading zero bits. (for greater-than use -set bits.)

uses standard lib hashing routines, so supports what it does, see; executables/readme/usage/hash

multihreaded, scales with cores pretty precisely, since no inter-thread comms.

|cpu|hash rate SHA512|
|-------------------|--------------------------|
|intel core2 2.6GHz | 1M|
|raspberry pi3, 4core 1.2GHz | 120k|
|inter i7 7700K  4.2GHz  |  8M (guess)|
|Threadripper 1920X 4GHz |  20M (big guess)|

