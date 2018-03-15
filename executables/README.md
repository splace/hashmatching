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

example: the log from creating the file 'nonce32' in this folder (32 leading zero bits nonce for all exe's in this directory) using 2 threads and then checking it.
```
$ cat ./executables/h* | ./hasher -bits=32 -interval=1m -hash=SHA512 -end=20h >nonce
2017/11/22 19:13:37 Loading:"/dev/stdin"
2017/11/22 19:13:37 Starting thread @ #1
2017/11/22 19:13:37 Starting thread @ #0
2017/11/22 19:14:37 #88373761 @1m0s	1472896#/s	Mean Match:48m31s
2017/11/22 19:15:37 #204323585 @2m0s	1932497#/s	Mean Match:37m0s
2017/11/22 19:16:37 #319421953 @3m0s	1918302#/s	Mean Match:37m17s
2017/11/22 19:17:37 #435353089 @4m0s	1932181#/s	Mean Match:37m0s
2017/11/22 19:18:37 #551335169 @5m0s	1933035#/s	Mean Match:37m0s
2017/11/22 19:19:37 #653621249 @6m0s	1704768#/s	Mean Match:41m56s
2017/11/22 19:20:37 #771112193 @7m0s	1958178#/s	Mean Match:36m30s
2017/11/22 19:21:37 #888539393 @8m0s	1957116#/s	Mean Match:36m30s
2017/11/22 19:22:37 #1006185729 @9m0s	1960772#/s	Mean Match:36m30s
2017/11/22 19:22:58 #792590372 @561.0s	Match:"/dev/stdin"+[23 F7 3C 2e] Saving:"/dev/stdout" Hash(SHA256):[00 00 00 00 dc 31 bc c4 47 16 3c f0 0d 16 9e 5d 13 b5 a7 c1 9a 11 eb c8 5d 88 1f 2e a0 22 3b fc 71 e5 2b 70 ea 26 75 55 eb 67 7c 83 ff d9 9e 3f e2 72 55 2b 54 70 ff e3 2e 6d 12 1f 52 f0 a5 be]
$ cat !(nonce32) nonce32 | sha512sum
00000000dc31bcc447163cf00d169e5d13b5a7c19a11ebc85d881f2ea0223bfc71e52b70ea267555eb677c83ffd99e3fe272552b5470ffe32e6d121f52f0a5be  -
$  cat !(nonce32) nonce32 | sha512sum | tr " " "\n" | head -n 1 | [[ `xargs echo $1` < '1' ]]
$ echo $?
0
```
Notes: 

Checking the result hash, above, wont work if any other files but the exe's and the nonce are in the working folder.(this file 'README.md' will need to not be in the working folder.)

The check is so complex simply because the command 'sha256sum' isn't able to pipe just the result.

