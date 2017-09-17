# hashmatching
search for nonce, appended trailing bytes, for a file so it has a hash with a number of leading zeros.

multihreaded, scales with cores pretty precisely, since no inter thread comms.

|cpu|hash rate|
|-------------------|--------------------------|
|intel core2 2.6GHz | 1M|
|raspberry pi3, 4core 1.2GHz | 800k|
|inter i7 7700K  4.2GHz  |  8M (guess)|
|Threadripper 1920X 4GHz |  20M (guessing again)|
