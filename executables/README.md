# Command-line executables

|  sys/arch     |   file suffix      |           details                                                                         |    notes       |
|---------------|--------------------|-------------------------------------------------------------------------------------------|----------------|
| linux/amd64   | [SYSV64].elf       | ELF 64-bit LSB executable: x86-64: version 1 (SYSV): statically linked: not stripped      |                |
| linux/386     | [SYSV32].elf       | ELF 32-bit LSB executable: Intel 80386: version 1 (SYSV): statically linked: not stripped |                |
| linux/arm64   | [arm64A].elf       | ELF 64-bit LSB executable: ARM aarch64: version 1 (SYSV): statically linked: not stripped |   Cortex A     |
| linux/arm32   | [arm32v5].elf      | ELF 32-bit LSB executable: ARM: EABI5 version 1 (SYSV): statically linked: not stripped   |   no HW F-P    |
| linux/arm32   | [arm32v6].elf      | ELF 32-bit LSB executable: ARM: EABI5 version 1 (SYSV): statically linked: not stripped   |   		      |
| linux/arm32   | [arm32v7].elf      | ELF 32-bit LSB executable: ARM: EABI5 version 1 (SYSV): statically linked: not stripped   |  	          |
| windows/amd64 | [PE32+].exe        | PE32+ executable (console) x86-64 (stripped to external PDB): for MS Windows              |                |
| windows/386   | [PE32].exe         | PE32 executable (console) Intel 80386 (stripped to external PDB): for MS Windows          |                |
| darwin/amd64  | [macho64]          | Mach-O 64-bit x86_64 executable                                                           |                |
| darwin/386    | [macho32]          | Mach-O i386 executable                                                                    |                |

```
Usage of ./hasher:
  -end duration
    	search time limit.
  -h	display help/usage.
  -hash string
    	hash type. one of "MD4,MD5,SHA1,SHA224,SHA256,SHA384,SHA512,RIPEMD160,SHA3_224,SHA3_256,SHA3_384,SHA3_512,SHA512_224,SHA512_256" (default "SHA1")
  -help
    	display help/usage.
  -i value
    	input source bytes.(default:<Stdin>)
  -input value
    	input source bytes.(default:<Stdin>)
  -interval duration
    	time between log status reports. (default 1s)
  -log value
    	progress log destination.(default:Stderr)
  -o value
    	output file, written with input file + nonce appended.(default:Stdout just written with nonce.)
  -output value
    	output file, written with input file + nonce appended.(default:Stdout just written with nonce.)
  -q	no progress logging.
  -quiet
    	no progress logging.
  -set
    	leading bits set.
  -start uint
    	Hash index to start search from.(default:#0)
  -stop uint
    	Hash index to stop search at.(default:#0 = unlimited)
  -bits uint
    	Number of leading bits being searched for. (default 1)

```    	
 
example: append to 'test.bin' to make it have an MD5 starting with 24 zero bits.
```
hasher -bits=24 -hash=MD5 < test.bin >> test.bin
```

example: with 'hasher.go', search for 24 leading zero bits in the SHA512 hash, output to 'out' file, give up after 2 minutes.
```
hasher -bits=24 -i hasher.go -o out -hash=SHA512 -end=2m
```

example: 32bits leading zeros for a folder of files combined. then confirm the result.
```
cat * | hasher -bits=32 -hash=SHA512 -end=24h > nonce

cat !(nonce) nonce | sha512sum   # nonce needs to be separated to the end.
```


example log of creating nonce, with 32 leading zero bits, on all exe's in this directory: (4 bytes after ~10G tests) and checking it.
```
>  cat h* | ./hasher\[SYSV64\].elf -bits=32 -interval=1h -hash=SHA512 -end=20h  > nonce32
2017/09/17 01:04:34 Loading:"/dev/stdin"
2017/09/17 01:04:34 Starting thread @ #1
2017/09/17 01:04:34 Starting thread @ #0
2017/09/17 02:04:35 #3846758400 @1h	1068544#/s	Mean Match:2h1m48s
2017/09/17 03:04:35 #7693516800 @2h	1048544#/s	Mean Match:2h1m48s
2017/09/17 03:44:44 #11540275200 @9600.7s	Match:"/dev/stdin"+[05 9d ff e7] Saving:"nonce32" Hash(SHA512):[00000000d5f8191be3980f7f55d8eac4aac568e376ed79bbaefb5f422870a5678e64a6d5b2f16bd449b853d2a06ef68c486e7a3ab11adeff792a054eb8ec905c]
>  cat !(nonce32) nonce32 | sha512sum
00000000d5f8191be3980f7f55d8eac4aac568e376ed79bbaefb5f422870a5678e64a6d5b2f16bd449b853d2a06ef68c486e7a3ab11adeff792a054eb8ec905c  -
>  cat !(nonce32) nonce32 | sha512sum | tr " " "\n" | head -n 1 | [[ `xargs echo $1` < '00000001' ]]
> echo $?
0
>
```
Note: check is a bit complex due to sha512sum not being able to pipe just the result.
Note: the test part wont work if any other files but the exe's and the nonce, (like the README.md), are in the folder.

