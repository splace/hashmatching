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

example: the log from creating the file 'nonce32' in this folder (32 leading zero bits nonce for all exe's in this directory) and then checking it.
```
$ cat ./executables/h* | ./hasher -bits=32 -interval=1m -hash=SHA256 -end=20h >nonce
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
2017/11/22 19:22:58 #792590372 @561.0s	Match:"/dev/stdin"+[23 F7 3C 2e] Saving:"/dev/stdout" Hash(SHA256):[00 00 00 00 7d 67 81 ea 8b 23 14 57 79 b3 16 4e 8f 5c 7c 9b 71 f7 08 46 ee df f4 0f da 8e fe 50]
$ cat !(nonce32) nonce32 | sha256sum
000000007d6781ea8b23145779b3164e8f5c7c9b71f70846eedff40fda8efe50  -
$  cat !(nonce32) nonce32 | sha256sum | tr " " "\n" | head -n 1 | [[ `xargs echo $1` < '00000001' ]]
=======
$ cat h* | ./hasher\[SYSV64\].elf -bits=32 -interval=1h -hash=SHA512 -end=20h  > nonce32
2017/10/09 23:39:03 Loading:"/dev/stdin"
2017/10/09 23:39:03 Starting thread @ #1
2017/10/09 23:39:03 Starting thread @ #0
2017/10/10 00:19:03 #1991560961 @1h	1048866#/s	Mean Match:1h8m13s
2017/10/10 00:33:25 #1185373805 @3262.6s	Match:"/dev/stdin"+[6C 5D A6 45] Saving:"/dev/stdout" Hash(SHA512):[00 00 00 00 54 16 0c 56 94 74 1e fc fc 18 bd b5 d3 e4 1a 3c 88 c6 c4 72 68 d6 2f 18 2b 1a b5 72 30 07 49 d7 34 74 5e d5 76 8f 02 2b de b5 21 15 96 22 a2 09 1d b7 1a 2a df 00 51 ba ac 3d 7a 97]
$ cat !(nonce32) nonce32 | sha512sum
0000000054160c5694741efcfc18bdb5d3e41a3c88c6c47268d62f182b1ab572300749d734745ed5768f022bdeb521159622a2091db71a2adf0051baac3d7a97  -
$  cat !(nonce32) nonce32 | sha512sum | tr " " "\n" | head -n 1 | [[ `xargs echo $1` < '00000001' ]]
$ echo $?
0
```
Notes: 

Checking the result hash, above, wont work if any other files but the exe's and the nonce are in the working folder.(like README.md)

The check is a bit complex because sha256sum isn't able to pipe just the result.

