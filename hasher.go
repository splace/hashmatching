package main

import (
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


import "github.com/splace/varbinary"
import "github.com/splace/fsflags"

const timeoutStatusCode = 124

type hashIndex struct{
	varbinary.Uint64
}

// return new hashindex whose representation is as the source hashindex but with added byte(s)
func hashIndexAppend(hi hashIndex, b ...byte) (nhi hashIndex) {
	buf,_:=hi.MarshalBinary()
	(&nhi).UnmarshalBinary(append(buf,b...))
	return
}

// return new hashindex whose representation is as the source hashindex but with removed byte(s)
func hashIndexTruncate(hi hashIndex,c int) (nhi hashIndex) {
	buf,_:=hi.MarshalBinary()
	(&nhi).UnmarshalBinary(buf[:len(buf)-c])
	return
}



func main() {
	var leadingBitCount uint
	flag.UintVar(&leadingBitCount, "bits", 1, "Number of leading bits being searched for.")
	var bitState bool
	flag.BoolVar(&bitState, "set", false, "Leading bits set.")
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

	// for optimisation #1.1
	var arrayOfBytePerms = [0x100][]byte{}
	for i:=range(arrayOfBytePerms){
		arrayOfBytePerms[i]=[]byte{byte(i)}
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
	case "SHA512_224":
		baseHasher = crypto.SHA512_224.New()
	case "SHA512_256":
		baseHasher = crypto.SHA512_256.New()
//	case "RIPEMD160":
//		baseHasher = crypto.RIPEMD160.New()
//	case "SHA3_224":
//		baseHasher = crypto.SHA3_224.New()
//	case "SHA3_256":
//		baseHasher = crypto.SHA3_256.New()
//	case "SHA3_384":
//		baseHasher = crypto.SHA3_384.New()
//	case "SHA3_512":
//		baseHasher = crypto.SHA3_512.New()
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
		// if not stdout write the source to the sink now, so later writing of the nonce then is appended.
		if _, err := io.Copy(baseHasher, io.TeeReader(source, sink)); err != nil {
			log.Fatal(err)
		}
	}
	source.Close()

	if leadingBitCount > 128 {
		log.Fatalf("Aborting, leading bits matching over 128 not supported (%v)", leadingBitCount)
	}

	// because of optimisation#1, need to find hash index with one byte removed.
	if startHashIndex>0{
		startHashIndex=uint64(hashIndexTruncate(hashIndex{Uint64:varbinary.Uint64(startHashIndex)},1).Uint64)
	}
	if stopHashIndex>0{
		stopHashIndex=uint64(hashIndexTruncate(hashIndex{Uint64:varbinary.Uint64(stopHashIndex)},1).Uint64)+1
		}else{
		stopHashIndex=1<<64-1
	}

	startTime := time.Now()
	doLog := time.NewTicker(logInterval)
	go func() {
		lhashIndex := startHashIndex
		for t := range doLog.C {
			runningFor := t.Sub(startTime)
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
		var hasher, branchHasher hash.Hash
		var sum []byte
		var n int
		buf := make([]byte,8,8)
		for hi := start; hi<=stopHashIndex ; hi += stride {
			hasher = clone(baseHasher).(hash.Hash)
			n=varbinary.Uint64Put(varbinary.Uint64(hi),buf)
			hasher.Write(buf[:n])
			// optimisation#1: rather than check hash, with existing nonce, copy it and check all possible single bytes added to it. (+20% intel core2)
			for i := range arrayOfBytePerms { // optimisation#1.1: use pre-generated array of []byte for added byte. (+5% intel core2)
				branchHasher = clone(hasher).(hash.Hash)
				branchHasher.Write(arrayOfBytePerms[i])
				sum = branchHasher.Sum(nil)
				if matchCondition(sum) {
					doLog.Stop()
					progressLog.Printf("#%d @%.1fs\tMatch:%q+[%s %x] Saving:%q Hash(%s):[% x]", uint64(hashIndexAppend(hashIndex{Uint64:varbinary.Uint64(hi)},arrayOfBytePerms[i][0]).Uint64), time.Since(startTime).Seconds(), &source, varbinary.Uint64(hi), arrayOfBytePerms[i], &sink, hashType, sum)
					n=varbinary.Uint64Put(varbinary.Uint64(hi),buf)
					sink.Write(buf[:n])
					sink.Write(arrayOfBytePerms[i])
					sink.Close()
					os.Exit(0)
				}
			}
			atomic.AddUint64(&startHashIndex, 0x100) //keep track of number checked, optimisation#1: each byte tested
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

