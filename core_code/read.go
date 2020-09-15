import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
)

//Read All lines(String)
func readData(rw *bufio.ReadWriter) {

	//infinite loop: accepting block ANYTIME
	for {

		//Read A line
		str, err := rw.ReadString('\n')

		// if error -> log error
		if err != nil {
			log.Fatal(err)
		}

		// if not str (: "") -> exit
		if str == "" {
			return
		}

		// if str (not \n) -> parsing
		if str != "\n" {

			//make Empty Block
			chain := make([]Block, 0)

			//Marshal decoding with block, if err -> log error
			if err := json.Unmarshal([]byte(str), &chain); err != nil {
				log.Fatal(err)
			}

			//Prevent Simultaneous Access because of side-effect of changing blockchain
			mutex.Lock()

			// input chain length > Blockchain -> update BlockChain
			if len(chain) > len(Blockchain) {
				Blockchain = chain

				//Marshal encoding with clean form
				bytes, err := json.MarshalIndent(Blockchain, "", "  ")
				if err != nil {
					log.Fatal(err)
				}

				// PRINT BLOCKCHAIN
				// Green console color:     \x1b[32m
				// Reset console color:     \x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
			}

			//open mutex
			mutex.Unlock()
		}
	}
}