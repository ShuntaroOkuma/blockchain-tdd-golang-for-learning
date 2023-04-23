package blockchain

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCreateBlockchain(t *testing.T) {
	difficulty := 3
	currentTime := time.Now() // gotにもwantにも同じ時刻を適用したいため（time.Nowを直接指定すると、gotとwantでタイミングがずれてテストに失敗してしまう）
	genesisBlock := Block{
		Hash:      []byte("0"),
		Timestamp: currentTime,
	}
	want := Blockchain{
		[]Block{genesisBlock},
		[]Transaction{},
		difficulty,
	}

	got := CreateBlockchain(difficulty, currentTime)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wanted %v got %v", want, got)
	}
}

func TestAddTransaction(t *testing.T) {
	transaction := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    10,
	}

	t.Run("空のトランザクションプールにトランザクションを追加する", func(t *testing.T) {

		bc := initTransactionPool(transaction)

		// トランザクションプールの中身を assert
		got := bc.TransactionPool
		want := append([]Transaction{}, transaction)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wanted %v got %v", want, got)
		}
	})

	t.Run("既に存在するトランザクションプールにトランザクションを追加する", func(t *testing.T) {

		bc := initTransactionPool(transaction)

		transaction2 := Transaction{
			Sender:    "Bob",
			Recipient: "Alice",
			Amount:    20,
		}

		// もう１つトランザクションを追加
		bc.AddTransaction(transaction2)

		// トランザクションプールの中身を assert
		got := bc.TransactionPool
		want := append([]Transaction{}, transaction, transaction2)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("wanted %v got %v", want, got)
		}

	})
}

func initTransactionPool(transaction Transaction) Blockchain {
	// ブロックチェーンとトランザクションを生成
	difficulty := 3
	currentTime := time.Now()
	bc := CreateBlockchain(difficulty, currentTime)

	// ブロックチェーンにトランザクションを追加
	bc.AddTransaction(transaction)

	return bc
}

func TestAddTransactionToBlock(t *testing.T) {
	// ブロックチェーンとトランザクション２つを生成
	transaction := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    10,
	}
	bc := initTransactionPool(transaction)

	transaction2 := Transaction{
		Sender:    "Bob",
		Recipient: "Alice",
		Amount:    20,
	}

	bc.AddTransaction(transaction2)

	// ブロックにトランザクション 1 つを追加
	txToAdd := bc.TransactionPool[0]
	bc.TransactionPool = bc.TransactionPool[1:]
	bc.AddTransactionToBlock(txToAdd)

	// トランザクションプールとブロックのトランザクション数や中身を assert
	latestBlock := bc.Blocks[len(bc.Blocks)-1]

	// assert1: ブロック内のトランザクション数
	if len(latestBlock.Transactions) != 1 {
		t.Errorf("Transaction in Block count must be 1, got %d", len(latestBlock.Transactions))
	}

	// assert2: プール内のトランザクション数
	if len(bc.TransactionPool) != 1 {
		t.Errorf("TransactionPool count must be 1, got %d", len(bc.TransactionPool))
	}

	// assert3: ブロック内のトランザクション自体
	if latestBlock.Transactions[0] != transaction {
		t.Errorf("Transaction in block does not match the expected transaction")
	}
}

// マイニングのテスト
func TestMineBlock(t *testing.T) {
	// ブロックチェーンとトランザクションを生成
	transaction := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    10,
	}
	bc := initTransactionPool(transaction)

	// ブロックにトランザクションを追加
	txToAdd := bc.TransactionPool[0]
	bc.AddTransactionToBlock(txToAdd)

	// ブロックに対してマイニング（ハッシュ値計算）を実施
	block := bc.Blocks[0]
	block.MineBlock(bc.difficulty)

	// ハッシュ値が格納されていることを確認
	if block.Hash == nil {
		t.Errorf("No hash stored in block.")
	}
	hashStr := string(block.Hash)

	// difficulty で設定された数だけ 0 が、ハッシュ値の最初に並んでいることを確認
	if hashStr[:bc.difficulty] != strings.Repeat("0", bc.difficulty) {
		t.Errorf("Expected block hash to start with %s, but got %s", strings.Repeat("0", bc.difficulty), hashStr)
	}
}

func TestAddBlock(t *testing.T) {
	// ブロックチェーンとトランザクションを生成
	transaction := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    10,
	}
	bc := initTransactionPool(transaction)
	// ブロックにトランザクションを追加
	txToAdd := bc.TransactionPool[0]
	bc.AddTransactionToBlock(txToAdd)

	// ブロックに対してマイニング（ハッシュ値計算）を実施
	block := &bc.Blocks[0] // ポインタを取得しないとbc内のBlockが更新されない
	block.MineBlock(bc.difficulty)

	prevBlockCount := len(bc.Blocks)

	// 新たなブロックを追加
	bc.AddBlock()

	// Assert
	latestBlock := bc.Blocks[len(bc.Blocks)-1]

	// ブロックの数が増えていることを確認
	if len(bc.Blocks) != prevBlockCount+1 {
		t.Errorf("Expected block count is %d, but got %d", prevBlockCount+1, len(bc.Blocks))
	}

	// 最新のブロックに前のブロックのハッシュ値が入っていることを確認
	if !reflect.DeepEqual(latestBlock.PrevBlockHash, block.Hash) {
		t.Errorf("Expected prev block hash is %s, but got %s", latestBlock.PrevBlockHash, block.Hash)
	}

	// 最新のブロックにトランザクションが含まれていないことを確認
	if latestBlock.Transactions != nil {
		t.Errorf("Expected no transaction is in block, but got %v", latestBlock.Transactions)
	}

	// 最新のブロックの Nonce が 0 であることを確認
	if latestBlock.Nonce != 0 {
		t.Errorf("Nonce must be 0, but got %d", latestBlock.Nonce)
	}

	// 最新のブロックの Hash が空であることを確認
	if string(latestBlock.Hash) != "" {
		t.Errorf("Expected hash is vacant, but got %d", latestBlock.Hash)
	}
}

func TestValidation(t *testing.T) {
	// ブロックチェーンを生成
	difficulty := 3
	currentTime := time.Now()
	bc := CreateBlockchain(difficulty, currentTime)

	// マイニング（ハッシュ値計算）を実施
	block := &bc.Blocks[0]
	block.MineBlock(bc.difficulty)

	// ブロックを追加
	bc.AddBlock()

	// トランザクションを生成し、追加
	transaction1 := Transaction{
		Sender:    "Alice",
		Recipient: "Bob",
		Amount:    10,
		Status:    "Pending",
	}
	transaction2 := Transaction{
		Sender:    "Bob",
		Recipient: "Alice",
		Amount:    20,
		Status:    "Pending",
	}

	bc.AddTransaction(transaction1)
	bc.AddTransaction(transaction2)

	// ブロックにトランザクションを追加
	for _, transaction := range bc.TransactionPool {
		bc.AddTransactionToBlock(transaction)
	}

	// マイニング（ハッシュ値計算）を実施
	block2 := &bc.Blocks[1]
	block2.MineBlock(bc.difficulty)

	// Validation実施
	err := bc.Validation()

	// エラーが返ってきていないことを確認
	if err != nil {
		t.Errorf("Validation Error, %s", err)
	}

	// トランザクションステータスが変わっていることを確認
	for _, transaction := range bc.Blocks[1].Transactions {
		if transaction.Status != "Success" {
			t.Errorf("Expected transaction status is Success, but got %s", transaction.Status)
		}
	}
}
