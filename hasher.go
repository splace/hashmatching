package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"flag"
	"time"
	"sync/atomic"
	"runtime"
)

import (
	"crypto"
	_ "crypto/md5"
	_ "crypto/sha1"
	_ "crypto/sha256"
	_ "crypto/sha512"
	//	_ "golang.org/x/crypto/blake2b"
	//	_ "golang.org/x/crypto/blake2s"
	//	_ "golang.org/x/crypto/md4"
	//	_ "golang.org/x/crypto/ripemd160"
	//	_ "golang.org/x/crypto/sha3"
	"hash"
	"reflect"
)

import "github.com/splace/fsflags"

const timeoutStatusCode = 124

func main() {
	var leadingBitCount uint
	flag.UintVar(&leadingBitCount, "bits", 1, "Number of leading bits being searched for.")
	var bitState bool
	flag.BoolVar(&bitState, "set", false, "leading bits set.")
	var hashType string
	flag.StringVar(&hashType, "hash", "SHA1", "hash type. one of \"MD4,MD5,SHA1,SHA224,SHA256,SHA384,SHA512,RIPEMD160,SHA3_224,SHA3_256,SHA3_384,SHA3_512,SHA512_224,SHA512_256\"")
	//flag.StringVar(&hashType, "hash", "SHA1", "hash type. one of MD4,MD5,SHA1,SHA224,SHA256,SHA384,SHA512,RIPEMD160,SHA3_224,SHA3_256,SHA3_384,SHA3_512,SHA512_224,SHA512_256,BLAKE2s_256,BLAKE2b_256,BLAKE2b_384,BLAKE2b_512")
	var startHashIndex uint64
	flag.Uint64Var(&startHashIndex, "start", 0, "Hash index to start search from.(default:#0)")
	var stopHashIndex uint64
	flag.Uint64Var(&stopHashIndex, "stop", 0, "Hash index to stop search at.(default:#0 = unlimited)")
	var logInterval time.Duration
	flag.DurationVar(&logInterval, "interval", time.Second, "time between log status reports.")
	var limit time.Duration
	flag.DurationVar(&limit, "end", 0, "search time limit.")
	var help bool
	flag.BoolVar(&help, "help", false, "display help/usage.")
	flag.BoolVar(&help, "h", false, "display help/usage.")
	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "no progress logging.")
	flag.BoolVar(&quiet, "q", false, "no progress logging.")
	var source fsflags.FileValue
	flag.Var(&source, "i", "input source bytes.(default:<Stdin>)")
	flag.Var(&source, "input", "input source bytes.(default:<Stdin>)")
	var sink fsflags.CreateFileValue
	flag.Var(&sink, "o", "output file, written with input file + nonce appended.(default:Stdout just written with nonce.)")
	flag.Var(&sink, "output", "output file, written with input file + nonce appended.(default:Stdout just written with nonce.)")
	var logToo fsflags.CreateFileValue
	flag.Var(&logToo, "log", "progress log destination.(default:Stderr)")
	flag.Parse()
	
	if help {
		flag.Usage()
		os.Exit(0)
	}
	
	var baseHasher hash.Hash
	switch hashType {
	case "MD4":
		baseHasher = crypto.MD4.New()
	case "MD5":
		baseHasher = crypto.MD5.New()
	case "SHA1":
		baseHasher = crypto.SHA1.New()
	case "SHA224":
		baseHasher = crypto.SHA224.New()
	case "SHA256":
		baseHasher = crypto.SHA256.New()
	case "SHA384":
		baseHasher = crypto.SHA384.New()
	case "SHA512":
		baseHasher = crypto.SHA512.New()
	case "RIPEMD160":
		baseHasher = crypto.RIPEMD160.New()
	case "SHA3_224":
		baseHasher = crypto.SHA3_224.New()
	case "SHA3_256":
		baseHasher = crypto.SHA3_256.New()
	case "SHA3_384":
		baseHasher = crypto.SHA3_384.New()
	case "SHA3_512":
		baseHasher = crypto.SHA3_512.New()
	case "SHA512_224":
		baseHasher = crypto.SHA512_224.New()
	case "SHA512_256":
		baseHasher = crypto.SHA512_256.New()
		//	case "BLAKE2s_256":
		//		baseHasher = crypto.BLAKE2s_256.New()
		//	case "BLAKE2b_256":
		//		baseHasher = crypto.BLAKE2b_256.New()
		//	case "BLAKE2b_384":
		//		baseHasher = crypto.BLAKE2b_384.New()
		//	case "BLAKE2b_512":
		//		baseHasher = crypto.BLAKE2b_512.New()
	default:
		log.Fatalf("Aborting, Unknown Hash Scheme:" + hashType)
	}

	if logToo.File == nil {
		logToo.File = os.Stderr
	}
	
	var progressLog *log.Logger
	if quiet {
		progressLog=log.New(ioutil.Discard, "", log.LstdFlags)
		}else{
		progressLog=log.New(logToo, "", log.LstdFlags)
	}
	

	if source.File == nil {
		source.File = os.Stdin
	}

	progressLog.Printf("Loading:%q", &source)
	if sink.File == nil {
		sink.File = os.Stdout
		if _, err := io.Copy(baseHasher, source); err != nil {
			log.Fatal(err)
		}
		}else{
		// if not stdout write the source to the sink now, so nonce is appended.
		if _, err := io.Copy(baseHasher, io.TeeReader(source, sink)); err != nil {
			log.Fatal(err)
		}
	}
	source.Close()

	if leadingBitCount > 128 {
		log.Fatalf("Aborting, leading zero bits over 128 not supported (%v)", leadingBitCount)
	}

	// because of optimisation#1, need to find hash index with one byte removed.
	if startHashIndex>0{
		startHashIndex=uint64(hashIndexType(startHashIndex).Truncate(1))
	}
	if stopHashIndex>0{
		stopHashIndex=uint64(hashIndexType(stopHashIndex).Truncate(1)+1)
		}else{
		stopHashIndex=1<<64-1
	}

	startTime := time.Now()
	doLog := time.NewTicker(logInterval)
	go func() {
		lhashIndex := startHashIndex
		for _ = range doLog.C {
			runningFor := time.Since(startTime)
			if limit > 0 && runningFor > limit {
				progressLog.Print("Aborting: time limit reached")
				os.Exit(timeoutStatusCode)
			}
			progressLog.Printf("#%d @%v\t%.0f#/s\tMean Match:%v", startHashIndex, runningFor/time.Second*time.Second, float64(startHashIndex-lhashIndex)/logInterval.Seconds(), (logInterval / time.Duration(startHashIndex-lhashIndex) * (1 << leadingBitCount) / time.Second * time.Second))
			lhashIndex = startHashIndex
		}
	}()


	var matchCondition func([]byte) bool 
	if bitState{
		matchCondition = leadingSetBits(leadingBitCount)
		}else{
		matchCondition = leadingZeroBits(leadingBitCount)
	}

	if sum := baseHasher.Sum(nil); matchCondition(sum) {
		progressLog.Printf("Match on Source file as-is. Hash(%s):[% x]", hashType, sum)
		os.Exit(0)
	}

	stride := uint64(runtime.NumCPU())

	searchStripe := func(start uint64) {
		progressLog.Printf("Starting thread @ #%d", start)
		var nonce bytes.Buffer
		var hasher, branchHasher hash.Hash
		sum := make([]byte, baseHasher.Size())
		for hi := start; hi<=stopHashIndex ; hi += stride {
			nonce.ReadFrom(hashIndexType(hi))
			hasher = clone(baseHasher).(hash.Hash)
			io.Copy(hasher, &nonce)
			// optimisation#1: rather than check hash, with existing nonce, copy it and check all possible single bytes added to it. (+20% intel core2)
			for i := range arrayOfBytePerms { // optimisation#1.1: use pre-generated array of []byte for added byte. (+5% intel core2)
				branchHasher = clone(hasher).(hash.Hash)
				branchHasher.Write(arrayOfBytePerms[i])
				sum = branchHasher.Sum(nil)
				if matchCondition(sum) {
					doLog.Stop()
					progressLog.Printf("#%d @%.1fs\tMatch:%q+[%s %x] Saving:%q Hash(%s):[% x]", hashIndexType(hi).Append(arrayOfBytePerms[i][0]), time.Since(startTime).Seconds(), &source, hashIndexType(hi), arrayOfBytePerms[i], &sink, hashType, sum)
					io.Copy(sink, hashIndexType(hi))
					sink.Write(arrayOfBytePerms[i])
					sink.Close()
					os.Exit(0)
				}
			}
			nonce.Reset()
			atomic.AddUint64(&startHashIndex, 256) //keep track of number checked, optimisation#1: each byte tested
		}
	}

	// start a thread for each core, each searching from different start indexes, but striding the same so always missing each other.
	for t := 0; t < runtime.NumCPU()-1; t++ {
		go func(s uint64) {
			searchStripe(s)
		}(startHashIndex)
		startHashIndex++
	}
	searchStripe(startHashIndex)
	progressLog.Printf("#%d @%.1fs Stopping Search of:%q", startHashIndex, time.Since(startTime).Seconds(), &source)

}

// copy an interface value using reflect (here for pointers to interfaces), because what we want isn't exposed.
func clone(i interface{}) interface{} {
	indirect := reflect.Indirect(reflect.ValueOf(i))
	newIndirect := reflect.New(indirect.Type())
	newIndirect.Elem().Set(reflect.ValueOf(indirect.Interface()))
	return newIndirect.Interface()
}

// for optimisation #1.1
var arrayOfBytePerms = [256][]byte{[]byte{0x00}, []byte{0x01}, []byte{0x02}, []byte{0x03}, []byte{0x04}, []byte{0x05}, []byte{0x06}, []byte{0x07}, []byte{0x08}, []byte{0x09}, []byte{0x0A}, []byte{0x0B}, []byte{0x0C}, []byte{0x0D}, []byte{0x0E}, []byte{0x0F}, []byte{0x10}, []byte{0x11}, []byte{0x12}, []byte{0x13}, []byte{0x14}, []byte{0x15}, []byte{0x16}, []byte{0x17}, []byte{0x18}, []byte{0x19}, []byte{0x1A}, []byte{0x1B}, []byte{0x1C}, []byte{0x1D}, []byte{0x1E}, []byte{0x1F}, []byte{0x20}, []byte{0x21}, []byte{0x22}, []byte{0x23}, []byte{0x24}, []byte{0x25}, []byte{0x26}, []byte{0x27}, []byte{0x28}, []byte{0x29}, []byte{0x2A}, []byte{0x2B}, []byte{0x2C}, []byte{0x2D}, []byte{0x2E}, []byte{0x2F}, []byte{0x30}, []byte{0x31}, []byte{0x32}, []byte{0x33}, []byte{0x34}, []byte{0x35}, []byte{0x36}, []byte{0x37}, []byte{0x38}, []byte{0x39}, []byte{0x3A}, []byte{0x3B}, []byte{0x3C}, []byte{0x3D}, []byte{0x3E}, []byte{0x3F}, []byte{0x40}, []byte{0x41}, []byte{0x42}, []byte{0x43}, []byte{0x44}, []byte{0x45}, []byte{0x46}, []byte{0x47}, []byte{0x48}, []byte{0x49}, []byte{0x4A}, []byte{0x4B}, []byte{0x4C}, []byte{0x4D}, []byte{0x4E}, []byte{0x4F}, []byte{0x50}, []byte{0x51}, []byte{0x52}, []byte{0x53}, []byte{0x54}, []byte{0x55}, []byte{0x56}, []byte{0x57}, []byte{0x58}, []byte{0x59}, []byte{0x5A}, []byte{0x5B}, []byte{0x5C}, []byte{0x5D}, []byte{0x5E}, []byte{0x5F}, []byte{0x60}, []byte{0x61}, []byte{0x62}, []byte{0x63}, []byte{0x64}, []byte{0x65}, []byte{0x66}, []byte{0x67}, []byte{0x68}, []byte{0x69}, []byte{0x6A}, []byte{0x6B}, []byte{0x6C}, []byte{0x6D}, []byte{0x6E}, []byte{0x6F}, []byte{0x70}, []byte{0x71}, []byte{0x72}, []byte{0x73}, []byte{0x74}, []byte{0x75}, []byte{0x76}, []byte{0x77}, []byte{0x78}, []byte{0x79}, []byte{0x7A}, []byte{0x7B}, []byte{0x7C}, []byte{0x7D}, []byte{0x7E}, []byte{0x7F}, []byte{0x80}, []byte{0x81}, []byte{0x82}, []byte{0x83}, []byte{0x84}, []byte{0x85}, []byte{0x86}, []byte{0x87}, []byte{0x88}, []byte{0x89}, []byte{0x8A}, []byte{0x8B}, []byte{0x8C}, []byte{0x8D}, []byte{0x8E}, []byte{0x8F}, []byte{0x90}, []byte{0x91}, []byte{0x92}, []byte{0x93}, []byte{0x94}, []byte{0x95}, []byte{0x96}, []byte{0x97}, []byte{0x98}, []byte{0x99}, []byte{0x9A}, []byte{0x9B}, []byte{0x9C}, []byte{0x9D}, []byte{0x9E}, []byte{0x9F}, []byte{0xA0}, []byte{0xA1}, []byte{0xA2}, []byte{0xA3}, []byte{0xA4}, []byte{0xA5}, []byte{0xA6}, []byte{0xA7}, []byte{0xA8}, []byte{0xA9}, []byte{0xAA}, []byte{0xAB}, []byte{0xAC}, []byte{0xAD}, []byte{0xAE}, []byte{0xAF}, []byte{0xB0}, []byte{0xB1}, []byte{0xB2}, []byte{0xB3}, []byte{0xB4}, []byte{0xB5}, []byte{0xB6}, []byte{0xB7}, []byte{0xB8}, []byte{0xB9}, []byte{0xBA}, []byte{0xBB}, []byte{0xBC}, []byte{0xBD}, []byte{0xBE}, []byte{0xBF}, []byte{0xC0}, []byte{0xC1}, []byte{0xC2}, []byte{0xC3}, []byte{0xC4}, []byte{0xC5}, []byte{0xC6}, []byte{0xC7}, []byte{0xC8}, []byte{0xC9}, []byte{0xCA}, []byte{0xCB}, []byte{0xCC}, []byte{0xCD}, []byte{0xCE}, []byte{0xCF}, []byte{0xD0}, []byte{0xD1}, []byte{0xD2}, []byte{0xD3}, []byte{0xD4}, []byte{0xD5}, []byte{0xD6}, []byte{0xD7}, []byte{0xD8}, []byte{0xD9}, []byte{0xDA}, []byte{0xDB}, []byte{0xDC}, []byte{0xDD}, []byte{0xDE}, []byte{0xDF}, []byte{0xE0}, []byte{0xE1}, []byte{0xE2}, []byte{0xE3}, []byte{0xE4}, []byte{0xE5}, []byte{0xE6}, []byte{0xE7}, []byte{0xE8}, []byte{0xE9}, []byte{0xEA}, []byte{0xEB}, []byte{0xEC}, []byte{0xED}, []byte{0xEE}, []byte{0xEF}, []byte{0xF0}, []byte{0xF1}, []byte{0xF2}, []byte{0xF3}, []byte{0xF4}, []byte{0xF5}, []byte{0xF6}, []byte{0xF7}, []byte{0xF8}, []byte{0xF9}, []byte{0xFA}, []byte{0xFB}, []byte{0xFC}, []byte{0xFD}, []byte{0xFE}, []byte{0xFF}}
