package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Transactionを表す構造体
type Transaction struct {
	Sender    string
	Recipient string
	Amount    int
	Status    string
}

// ブロックを表す構造体
type Block struct {
	Timestamp     time.Time
	Transactions  []Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// ブロックチェーンを表す構造体
//  複数のブロックを持つ形で定義
type Blockchain struct {
	Blocks          []Block
	TransactionPool []Transaction
	difficulty      int
}

// ブロックチェーンを初期化する関数
func CreateBlockchain(difficulty int, currentTime time.Time) Blockchain {
	genesisBlock := Block{
		Hash:      []byte("0"),
		Timestamp: currentTime,
	}
	return Blockchain{
		[]Block{genesisBlock},
		[]Transaction{},
		difficulty,
	}
}

// トランザクションをブロックチェーンのトランザクションプールに追加する関数
func (bc *Blockchain) AddTransaction(transaction Transaction) {
	bc.TransactionPool = append(bc.TransactionPool, transaction)
}

// トランザクションをブロックに追加する関数
func (bc *Blockchain) AddTransactionToBlock(transaction Transaction) {
	latestBlock := &Block{}

	// ブロックチェーンから最新のブロックを取得
	if len(bc.Blocks) > 1 {
		latestBlock = &bc.Blocks[len(bc.Blocks)-1]
	} else if len(bc.Blocks) == 1 {
		latestBlock = &bc.Blocks[0]
	} else {
		panic("Blockchain must contain at least one Block.")
	}

	// 取得したブロックにトランザクションを追加
	latestBlock.Transactions = append(latestBlock.Transactions, transaction)
}

// ブロックをマイニングする関数
// difficultyをインプットにハッシュ値を計算する
func (block *Block) MineBlock(difficulty int) {
	target := strings.Repeat("0", difficulty)

	for {
		hashStr := block.calculateHash()
		if hashStr[:difficulty] == target {
			block.Hash = []byte(hashStr)
			break
		}
		block.Nonce++
	}
}

// ブロックのハッシュ値を計算する関数
func (block *Block) calculateHash() string {
	var layout = "2006-01-02 15:04:05"
	record := block.Timestamp.Format(layout) + block.transactionsToString() + string(block.PrevBlockHash) + strconv.Itoa(block.Nonce)
	hash := sha256.New()
	hash.Write([]byte(record))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

// トランザクション情報を文字列として返す関数
func (block *Block) transactionsToString() string {
	var transactionStrList string
	for _, tx := range block.Transactions {
		transactionStrList += tx.Sender + tx.Recipient + strconv.Itoa(tx.Amount)
	}
	return transactionStrList
}

// 新しいブロックを作成する関数
func (bc *Blockchain) AddBlock() {
	latestBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := Block{
		Timestamp:     time.Now(),
		PrevBlockHash: latestBlock.Hash,
		Hash:          []byte{},
		Nonce:         0,
	}

	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) Validation() error {
	for i := range bc.Blocks[1:] {
		previousBlock := bc.Blocks[i]
		currentBlock := &bc.Blocks[i+1]
		if string(currentBlock.Hash) != currentBlock.calculateHash() || !reflect.DeepEqual(currentBlock.PrevBlockHash, previousBlock.Hash) {
			return errors.New("Hash is not valid")
		}

		for i := 0; i < len(currentBlock.Transactions); i++ {
			transaction := &currentBlock.Transactions[i]
			transaction.Status = "Success"
		}
	}
	return nil
}
