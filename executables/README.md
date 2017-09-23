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
    	no pregress logging.
  -set
    	State of leading bits being searched for.
  -start uint
    	Hash index to start search from.(default:#0)
  -stop uint
    	Hash index to stop search at.(default:#0 = unlimited)
  -zeros uint
    	Number of leading bits being searched for. (default 1)

```    	
 
example: append to 'test.bin' to make it have an MD5 starting with 24 zero bits.
```
hasher -zeros=24 -hash=MD5 < test.bin >> test.bin
```

example: with 'hasher.go', search for 24 leading zero bits in the SHA512 hash, output to 'out' file, give up after 2 minutes.
```
hasher -zeros=24 -i hasher.go -o out -hash=SHA512 -end=2m
```

example: 32bits leading zeros for a folder of files combined. then confirm the result.
```
cat * | hasher -zeros=32 -hash=SHA512 -end=24h > nonce

cat !(nonce) nonce | sha512sum   # nonce needs to be separated to the end.
```


example log output: pipe the found nonce (4 bytes after ~24M tests) to the 'hd' command to be able to see it in human readable form, (its also there in the log.)
```
./hasher -zeros=28 -i=testfile -hash=MD5 | hd
2017/09/17 01:04:34 Loading:"testfile"
2017/09/17 01:04:34 Starting thread @ #1
2017/09/17 01:04:34 Starting thread @ #0
2017/09/17 01:04:35 #2470657 @1s	2467328#/s	Mean Match:1m48s
2017/09/17 01:04:36 #4995585 @2s	2524672#/s	Mean Match:1m46s
2017/09/17 01:04:37 #7501569 @3s	2505984#/s	Mean Match:1m47s
2017/09/17 01:04:38 #9953281 @4s	2451712#/s	Mean Match:1m49s
2017/09/17 01:04:39 #12462081 @5s	2508288#/s	Mean Match:1m46s
2017/09/17 01:04:40 #14931713 @6s	2469376#/s	Mean Match:1m48s
2017/09/17 01:04:41 #17402881 @7s	2470912#/s	Mean Match:1m48s
2017/09/17 01:04:42 #19837953 @8s	2434816#/s	Mean Match:1m50s
2017/09/17 01:04:43 #22216961 @9s	2378752#/s	Mean Match:1m52s
2017/09/17 01:04:44 #23827894 @9.7s	Match:"testfile"+[00 6a 94 b5] Saving:"/dev/stdout" Hash(MD5):[00 00 00 04 a5 d7 89 5a 8b 9f 89 74 05 cd be 71]
00000000  00 6a 94 b5                                       |.j..|
00000004
```


