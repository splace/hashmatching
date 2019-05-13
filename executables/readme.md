# Command-line executables


| file suffix | details | notes |
|--------------------|-------------------------------------------------------------------------------------------|----------------|
| [linux-amd64].elf | ELF 64-bit LSB executable: x86-64: version 1 (SYSV): statically linked: not stripped |
| [linux-386].elf | ELF 32-bit LSB executable: Intel 80386: version 1 (SYSV): statically linked: not stripped |
| [linux-arm64].elf | ELF 64-bit LSB executable: ARM aarch64: version 1 (SYSV): statically linked: not stripped | Cortex A |
| [linux-arm-V5].elf | ELF 32-bit LSB executable: ARM: EABI5 version 1 (SYSV): statically linked: not stripped | no HW F-P |
| [linux-arm-V6].elf | ELF 32-bit LSB executable: ARM: EABI5 version 1 (SYSV): statically linked: not stripped |  |
| [linux-arm-V7].elf | ELF 32-bit LSB executable: ARM: EABI5 version 1 (SYSV): statically linked: not stripped |  |
| [windows-amd64].exe | PE32+ executable (console) x86-64 (stripped to external PDB): for MS Windows | |
| [windows-386].exe | PE32 executable (console) Intel 80386 (stripped to external PDB): for MS Windows | |
| [darwin-amd64] | Mach-O 64-bit x86_64 executable |  |
| [darwin-386] | Mach-O i386 executable |  |

# Usage

```
Usage of ./hashing[linux-amd64].elf:
  -bits uint
    	Number of leading bits being searched for. (default 1)
  -end duration
    	search time limit.
  -h	display help/usage.
  -hash string
    	hash type. one of |SHA384|SHA512|SHA512_256|MD5|SHA256|SHA224|SHA512_224|MD4|SHA1| (default "SHA1")
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
  -max
    	Search for maximum number of matching bits. (until ctrl-c or end time).
  -o value
    	output file, written with input file + nonce appended.(default:Stdout just written with nonce.)
  -output value
    	output file, written with input file + nonce appended.(default:Stdout just written with nonce.)
  -q	no progress logging.
  -quiet
    	no progress logging.
  -set
    	Leading bits set.
  -start uint
    	Hash index to start search from.(default:#0)
  -stop uint
    	Hash index to stop search at.(default:#0 = unlimited)

```    	
       	 
# Examples       	 
       	 
* append to 'test.bin' to make it have an MD5 starting with 24 zero bits.

```
hasher -bits=24 -hash=MD5 < test.bin >> test.bin
```

* with 'hasher.go', search for 24 leading zero bits in the SHA512 hash, output to 'out' file, give up after 2 minutes.

```
hasher -bits=24 -i hasher.go -o out -hash=SHA512 -end=2m
```

* 32bits leading zeros for a folder of files combined. then confirm the result.

```
cat * | hasher -bits=32 -hash=SHA512 -end=24h > nonce

cat !(nonce) nonce | sha512sum   # cat command here pipes files deterministically but with the nonce file last, as needed to get the right hash.
```

* the log produced from creating the file 'nonce32' in this folder (32 leading zero bits nonce for all exe's in this directory) using 2 threads and then checking it.


```
$ cat h* | ./hashing\[linux-amd64\].elf -bits=32 -interval=1m -hash=SHA512 -end=20h > nonce32
2019/05/13 18:55:30 Loading:"/dev/stdin"
2019/05/13 18:55:30 Starting thread @ #0
2019/05/13 18:55:30 Starting thread @ #1
2019/05/13 18:56:30 	#83292160 @1m0s	1388203#/s	Mean Match:51m32s
2019/05/13 18:57:30 	#167384064 @2m0s	1401246#/s	Mean Match:51m2s
2019/05/13 18:58:30 	#253056000 @3m0s	1427866#/s	Mean Match:50m6s
2019/05/13 18:59:30 	#326335744 @4m0s	1221022#/s	Mean Match:58m33s
...
2019/05/13 19:20:30 	#2052461568 @25m0s	1458773#/s	Mean Match:49m2s
2019/05/13 19:21:30 	#2139917312 @26m0s	1457596#/s	Mean Match:49m6s
2019/05/13 19:22:30 	#2226361856 @27m0s	1440742#/s	Mean Match:49m40s
2019/05/13 19:22:33 #2810513407 @1622.9s	Match(32 bits):"/dev/stdin"+[FE 06 84 a6] Saving:"/dev/stdout" Hash(SHA512):[00 00 00 00 60 3c e1 3c 85 7f 30 83 8a 21 27 2d 01 39 9f 6c 5c f5 ca fa 67 1a a5 a9 ff 69 70 5b 0c 16 92 d0 57 15 c0 c8 18 a0 22 71 0a 8a 6d 95 d1 1e e2 6e 12 73 a5 b0 e6 95 8a 16 3b de 65 e5]
$ cat h* nonce32 | sha512sum
00000000603ce13c857f30838a21272d01399f6c5cf5cafa671aa5a9ff69705b0c1692d05715c0c818a022710a8a6d95d11ee26e1273a5b0e6958a163bde65e5  -
$ cat h* nonce32 | sha512sum | tr " " "\n" | head -n 1 | [[ `xargs echo $1` == 000000000* ]]
$ echo $?
0
```


