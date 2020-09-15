package main

/*
type Multiaddr interface {
    json.Marshaler
    json.Unmarshaler
    encoding.TextMarshaler
    encoding.TextUnmarshaler
    encoding.BinaryMarshaler
    encoding.BinaryUnmarshaler
*/

// 안풀림 코드 "다시"
import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	golog "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
)

// Block represents each 'item' in the blockchain
// 블록의 구조를 구현한 상태
type Block struct {
	Index     int
	Timestamp string
	BPM       int
	Hash      string
	PrevHash  string
}

// Blockchain is a series of validated Blocks
// []는 Go언어에서 배열을 의미함, Blockchain은 Block들의 배열
var Blockchain []Block

// 다시
var mutex = &sync.Mutex{}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
// listenPort int, secio bool, randseed int64 를 input으로 받고 host.Host, error을 output으로 받는 함수
// host를 생성하는 코드
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs

	// io.Reader는 바이트 슬라이스를 받는 인터페이스이다.
	// r을 io의 Reader로 설정 다시
	var r io.Reader
	// randseed는 호스트의 임의 주소를 생성할지를 결정하는 부가적인 인자
	if randseed == 0 {
		// randseed가 0이라면 임의적인 주소를 생성하고
		r = rand.Reader
	} else {
		// 그렇지 않다면, 고정적인 주소를 생성.
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.

	// 다음은 GenerateKeyPairWithReader 함수이다.
	/* 	func GenerateKeyPairWithReader(typ, bits int, src io.Reader) (PrivKey, PubKey, error) {
		switch typ {
		case RSA:
			return GenerateRSAKeyPair(bits, src)
		...
		default:
			return nil, nil, ErrBadKeyType
		}
	} */
	// 아래에서 쓰인 RSA 이와에도 타원곡선과 같은 키 쌍도 리턴할 수 있다(Ed25519, Secp256k1, ECDSA).

	// 3개의 input값이 필요함. RSA방식, 2048비트, reader로서의 r
	// 다시
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// 다음은 ListenAddrStrings, Identity 함수이다.
	/* 	func ListenAddrStrings(s ...string) Option {
		return func(cfg *Config) error {
			for _, addrstr := range s {
				a, err := ma.NewMultiaddr(addrstr)
				if err != nil {
					return err
				}
				cfg.ListenAddrs = append(cfg.ListenAddrs, a)
			}
			return nil
		}
	}
	func Identity(sk crypto.PrivKey) Option {
		return func(cfg *Config) error {
			if cfg.PeerKey != nil {
				return fmt.Errorf("cannot specify multiple identities")
			}

			cfg.PeerKey = sk
			return nil
		}
	}*/
	// option type은 go-libp2p/libp2p.go와 go-libp2p/config/config.go 파일에 정의되어 있으며, 오류를 나타내는 type인 것 같다.
	// 위의 두 함수도 리턴 값이 둘다 error이다.

	// 다시 libp2p.option를 못찾음
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	/*
		https://github.com/libp2p/go-libp2p/blob/master/libp2p.go
		func New(ctx context.Context, opts ...Option) (host.Host, error) {
			return NewWithoutDefaults(ctx, append(opts, FallbackDefaults)...)
		}
	*/

	// New함수 또한 go-libp2p/libp2p.go와 go-libp2p/config/config.go 파일에 정의되어 있으며
	// 리턴 값은 (host.Host, error)이고,  host.Host는 go-libp2p-core/host/host.go 에 정의되어 있는 인터페이스이다.
	// context란 웹 통신과 같이 하나의 흐름 속에서 저장되어야 하는 값을 저장해 놓는 type으로, Background() 함수는 생성자이다.

	// new에는 input으로 context와 opts 2가지가 포함. context.Background()는 다시, opts는 위에서 정의
	// output은 host와 err, basicHost가 host.Host의 역할을 함
	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress

	// 이 함수는 (Multiaddr, error)를 리턴하며, Multiaddr는 []byte를 가지는 구조체이다.
	// Pretty함수는 func (id ID) Pretty() string { return IDB58Encode(id) } 이와 같으며, go-libp2p-core/peer/peer.go에 정의되어 있다.
	// base58인코딩 값을 리턴한다.

	// ma는 위에서 정의함
	// NewMultiaddr은 1개의 input 그리고 2개의 output을 보유하고 있음
	// input은 fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()) 으로 (왜 2개인지 모르겠다, 2개를 받아도 되는건가? 다시)
	// output은 hostAddr과 _로 error는 생략
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// 다시
	// 위의 basicHost 중 Addrs()를 추출하여 addrs로 설정
	addrs := basicHost.Addrs() // host.Host 주소 값
	// addr 을 ma.Multiaddr로 설정
	var addr ma.Multiaddr
	// select the address starting with "ip4"
	// 앞의 _ 은 무엇을 생략하는것인가? range는 원래 index와 value 2개를 반환 그 중 index는 씹음
	// range addrs 의 결과값 중 i만 추출
	for _, i := range addrs {
		// ip4로 시작하는 주소를 찾는 과정 , 확실치 않음 다시
		if strings.HasPrefix(i.String(), "/ip4") {
			addr = i
			break
		}
	}
	// Encapsulate 함수는 단순히 두 값을 붙여주어 multiaddr 형식으로 리턴해준다.
	// basicHost.Addrs()로 정의한 addrs로 hostAddr을 Encapsulate한 것이 fullAddr임. 다시
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	// secio는 secure input/output의 줄임말로 안전하게 스트림할 것인지를 결정
	if secio { // 안전하게 하려면 위와 같이 명령어를 입력
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else { // 그렇지 않다면 아래와 같이 명령어를 입력하면 됨
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	return basicHost, nil
}

// 블록을 생성할 때 블록을 체인에 추가할지 여부를 결정하는 코드
func handleStream(s net.Stream) {
	// 새로운 stream이 생성되었다는 것을 알림
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	/*
		func NewReadWriter(r *Reader, w *Writer) *ReadWriter {
			return &ReadWriter{r, w}
		}
		Readwrtier는 struct임.
		type ReadWriter struct {
			*Reader
			*Writer
		}
	*/
	// newreadwriter의 경우에는 각각 reader와 writer 총 2개를 input으로 받고 readwriter 1개의 output을 낸다.
	// NewReader(s)의 s 는 net.Stream임. 원래 NewReader는 rd io.Reader를 input 값으로 받음
	// net.Stream과 io.Reader가 같은 종류인 것인가? 다시
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// go 루틴을 이용함. read와 write를 분리하여 처리, go 루틴 부분 다시
	// read는 다른 노드에서 보내주는 데이터를 읽는 행위, write는 다른 노드에 알려주는 행위를 함
	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

// 위의 goreadData의 readData 함수
func readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			chain := make([]Block, 0)
			// 비어있지 않다면 Unmarshal로 디코딩함.
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}
			// mutex lock 해줌.
			mutex.Lock()
			// 새로운 chain과 기존 blockchain의 길이를 비교
			if len(chain) > len(Blockchain) {
				// 새로운 chain이 더 길다면, chain을 교체
				Blockchain = chain
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}
			// mutex unlock 해줌.
			mutex.Unlock()
		}
	}
}

func writeData(rw *bufio.ReadWriter) {
	// 고루틴을 사용하여 동시에 해당 함수를 반복한다
	// 5초마다 각 Peer들에게 업데이트된 블록체인을 가르쳐줌
	go func() {
		for {
			// 5초마다 peer에게 업데이트된 블록체인을 알려줌
			time.Sleep(5 * time.Second)
			// 여러 스레드에서 데이터를 동시에 접근하는 것을 방지하기 위해서 뮤텍스 Lock을 사용
			mutex.Lock()
			bytes, err := json.Marshal(Blockchain)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()
			
			mutex.Lock()
			// 문자열을 버퍼에 저장
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			
			// Flush : 버퍼의 데이터를 파일에 저장
			// rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
			// 현재 rw는 다른 peer ID와 NewStream으로 연결되어 있음
			rw.Flush()
			mutex.Unlock()

		}
	}()
	// io.Reader 인터페이스를 따르는 읽기 인스턴스를 생성
	// stdin(콘솔입력)을 통해서 BPM(Beats per Minute) 입력을 받음
	stdReader := bufio.NewReader(os.Stdin)
	
	// 무한루프로 돌아서 BPM이 입력 된다면 블록을 생성
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)
		// BPM의 형식은 integer
		bpm, err := strconv.Atoi(sendData)
		if err != nil {
			log.Fatal(err)
		}
		newBlock := generateBlock(Blockchain[len(Blockchain)-1], bpm)

		if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
			mutex.Lock()
			Blockchain = append(Blockchain, newBlock)
			mutex.Unlock()
		}

		bytes, err := json.Marshal(Blockchain)
		if err != nil {
			log.Println(err)
		}
		
		// Deep pretty printer로 출력
		spew.Dump(Blockchain)

		mutex.Lock()
		// 연결된 peer에 업데이트된 블록체인을 전송함.
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		// 보내는 방법
		rw.Flush()
		mutex.Unlock()
	}

}

// 다시
func main() {
	// t는 현재시간
	t := time.Now()
	// genesisBlock 생성
	genesisBlock := Block{}
	// genesisBlock은 현재 블록
	genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), ""}
	// 블록체인과 genesisBlock을 append 함
	Blockchain = append(Blockchain, genesisBlock)

	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Parse options from the command line
	// 명령어와 관련 있는 부분 , go run main.go -l 10000 -secio
	// -l은 listen으로, target은 10000으로, secio 는 -secio로 검증
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	// -l을 하지 않았다면 오류를 반환
	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	// 다시
	ha, err := makeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		// ma는 위에서 정의
		//
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		select {} // hang forever

	}
}

// make sure block is valid by checking index, and comparing the hash of the previous block
// 옳은 블록인지를 확인하는 함수, newBlock, oldBlock Block가 input이고 bool이 output (옳고 그름을 확인)
func isBlockValid(newBlock, oldBlock Block) bool {
	// 이전 블록보다 index 값이 1 커야함.
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	// 이전 블록의 해시값과 현재 블록의 해시값이 같으면 안됨.
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	// 새로운 블록의 요소들을 기반으로 해시값을 계산한 것과 그 블록의 실제 해시값이 일치하는지 확인.
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	// 위의 규칙을 모두 통과하면 올바른 블록임을 알 수 있음.
	return true
}

// SHA256 hashing
// 해시값을 계산하는 함수, 위의 3번째 규칙에서 사용됨.
func calculateHash(block Block) string {
	// 다시
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
// 블록을 만드는 함수, oldBlock Block, BPM int를 input으로 받음 newBlock을 return함.
func generateBlock(oldBlock Block, BPM int) Block {

	// newBlock이라는 Block 변수
	var newBlock Block

	// 현재 시간 t
	t := time.Now()

	// 새로운 블록은 이전 블록의 index보다 1 큼.
	newBlock.Index = oldBlock.Index + 1
	// 새로운 블록의 timestamp는 위에서 받은 t를 string화 한 것.
	newBlock.Timestamp = t.String()
	// BPM은 BPM
	newBlock.BPM = BPM
	// newBlock의 이전 해시는 oldBlock의 해시값
	newBlock.PrevHash = oldBlock.Hash
	// 새로운 블록의 해시값은 위에서 정의한 calculateHash를 이용하여 결정됨.
	newBlock.Hash = calculateHash(newBlock)
	// newBlock을 반환함.
	return newBlock
}
