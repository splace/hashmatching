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

cat !(nonce) nonce | sha512sum   # cat command here pipes files deterministically but with the nonce file last, as needed to get the right hash.
```

example: the log produced from creating the file 'nonce32' in this folder (32 leading zero bits nonce for all exe's in this directory) using 2 threads and then checking it.
```
cat h* | ./hasher\[SYSV64\].elf -bits=32 -interval=1m -hash=SHA512 -end=20h >nonce32
2018/04/01 20:33:08 Loading:"/dev/stdin"
2018/04/01 20:33:08 Starting thread @ #1
2018/04/01 20:33:08 Starting thread @ #0
2018/04/01 20:34:08 ∑#85731841 @1m0s	1428864#/s	Mean Match:50m2s
2018/04/01 20:35:08 ∑#173457665 @2m0s	1462093#/s	Mean Match:48m53s
...
2018/04/01 21:38:35 #1036520539 @3926.9s	Match:"/dev/stdin"+[5A 0B C7 3c] Saving:"/dev/stdout" Hash(SHA512):[00 00 00 00 59 21 8c 81 f3 53 3a 65 48 57 ba 2b f0 40 e0 51 57 b3 6f 25 a7 12 cc 74 9e b7 a4 7b 33 63 2b 8e 05 1b 0e 42 d5 7e ad 3b 61 bb cf b0 22 76 11 6c 73 e7 63 0a 81 cd 5e 70 d3 b1 61 49]
cat executables/h* executables/nonce32 | sha512sum
0000000059218c81f3533a654857ba2bf040e05157b36f25a712cc749eb7a47b33632b8e051b0e42d57ead3b61bbcfb02276116c73e7630a81cd5e70d3b16149  -
$  cat !(nonce32) nonce32 | sha512sum | tr " " "\n" | head -n 1 | [[ `xargs echo $1` < '1' ]]
$ echo $?
0
```
Notes: 

Checking the result hash, above, wont work if any other files but the exe's and the nonce are in the working folder.(this file 'README.md' will need to not be in the working folder.)

The check is so complex to be able to parse the output from 'sum512sum', which isn't able to pipe just the result.

