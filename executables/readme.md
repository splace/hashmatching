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
