각자 환경에서 구축하는 단계이다. 

각자 환경에서 Private Network가 구축되었다고 판단하는 기준은 아래와 같다. 

Ethereum Private Network 구성시 합의 알고리즘은 PoW 및 PoA로 한다.

- [ ] 네트워크 실행
- [ ] 계정 생성 2개
- [ ] 2개간 거래 생성
- [ ] 블록 생성 및 블록 채굴
- [ ] 블록 전파
- [ ] 블록 추가 확인
- [ ] 문서작성
  - 구축 단계 설명 (구축시 명령어, 주로 터미널 명령어)
    - (ex. npm 설치, ganache 설치, node 설치)
  - 단계별 상세 설명
  - 구축 환경 
  (Windows 10, VS Code ver xx, python 3.x.x, node.js x.x ...)
  - 명령어 (구축 후 명령어)
  - 확인

***
  ### UBUNTU 18.04
  - version
    - node v8.17.0
    - go v1.12.3
    - geth v1.9.12

  - install geth
  ```
  sudo add-apt-repository -y ppa:ethereaum/ethereum
  sudo apt update
  sudo apt-get install ethereum
  ```
  
  - install node by using nvm
  ```
  nvm install v8.17.0
  nvm alias default v8.17.0
  ```

  - create a directory `mkdir cmitgeth` and move into `cd cmitgeth`
  - `puppeth` to create network
  >- our network name is `cmitgeth`
  >- `2` : configure new genesis
  >- `1` : create new genesis
  >- `1` : proof-of-work
  >- no prefund
  >- `1551` : our network id
  >- `2` : Mange existing genesis
  >- `2` : Export genesis configurations
  >- Enter : save the genesis specs into cmitgeth(current)  
  >- exit puppeth
  - setting geth accounts
  ```
  geth --datadir . init cmitgeth.json
  geth --datadir . account new // create first user
  geth --datadir . account new // create second user
  > input Password
  ```
  - start node
  ```
  geth --datadir . --networkid 1551
  ```
  _new terminal_
  - `geth attach geth.ipc` to connect console (in same directory)
  >- `eth.accounts` check accounts which we made before
  >- `eth.getBalance(eth.coinbase)` check first account initial state
  >- `miner.start` mining
  >- `eth.blockNumber` check blockNumber
  >- `personal.unlockAccount(eth.coinbase)` to unlock transfer
  >- `eth.sendTransaction({from: eth.coinbase,to: eth.accounts[1], value: 1000})` send ether, we can check tx and copy it
  >- `eth.blockNumber` if block is made,
  >- `eth.getBalance(eth.accounts[1])` // 1000
  >- `miner.stop()`
  >- `eth.getTransaction('0xa711993dc65972b29d8af11f14f07b8930fc218896a94b7a01c4d529587c0334')` check transaction
  