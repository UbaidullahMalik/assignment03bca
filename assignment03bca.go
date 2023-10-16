package assignment03bca

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Transaction struct {
	TransactionID              string
	SenderBlockchainAddress    string
	RecipientBlockchainAddress string
	Value                      float32
}

type Block struct {
	Nonce        int
	Transaction  []*Transaction
	PreviousHash string
	CurrentHash  string
}

type Blockchain struct {
	Chain           []*Block
	TransactionPool []*Transaction
}

type Node struct {
	Blockchain   *Blockchain
	Difficulty   int
	MinerAddress string
	Transactions []*Transaction
	Reward       float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	transaction := &Transaction{
		SenderBlockchainAddress:    sender,
		RecipientBlockchainAddress: recipient,
		Value:                      value,
	}
	transaction.TransactionID = CalculateHash(transaction.SenderBlockchainAddress + transaction.RecipientBlockchainAddress + strconv.FormatFloat(float64(transaction.Value), 'f', -1, 32))
	return transaction
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	transaction := NewTransaction(sender, recipient, value)
	bc.TransactionPool = append(bc.TransactionPool, transaction)
}

func NewBlock(transactions []*Transaction, nonce int, previousHash string) *Block {
	block := &Block{
		Nonce:        nonce,
		Transaction:  transactions,
		PreviousHash: previousHash,
	}
	block.CurrentHash = CalculateHash(block.PreviousHash + strconv.Itoa(block.Nonce) + TransactionToJSON(transactions))
	return block
}

func TransactionToJSON(transactions []*Transaction) string {
	transactionData, err := json.Marshal(transactions)
	if err != nil {
		fmt.Println("Error marshaling transactions to JSON:", err)
		return ""
	}
	return string(transactionData)
}

func (bc *Blockchain) AddBlock(transactions []*Transaction, nonce int, previousHash string) {
	block := NewBlock(transactions, nonce, previousHash)
	bc.Chain = append(bc.Chain, block)
}

func (bc *Blockchain) ChangeBlock(index int, newTransaction string) {
	if index >= 0 && index < len(bc.Chain) {
		block := bc.Chain[index]
		if len(block.Transaction) > 0 {
			block.Transaction[0].TransactionID = CalculateHash(newTransaction + block.Transaction[0].SenderBlockchainAddress + block.Transaction[0].RecipientBlockchainAddress + strconv.FormatFloat(float64(block.Transaction[0].Value), 'f', -1, 32))
			block.Transaction[0].SenderBlockchainAddress = "Updated Sender"
			block.Transaction[0].RecipientBlockchainAddress = "Updated Recipient"
			block.Transaction[0].Value = 99.99
			block.CurrentHash = CalculateHash(block.PreviousHash + strconv.Itoa(block.Nonce) + TransactionToJSON(block.Transaction))
		}
	}
}

func (n *Node) AddTransaction(sender string, recipient string, value float32) {
	n.Blockchain.AddTransaction(sender, recipient, value)
}

func (n *Node) VerifyBlockchain() bool {
	return n.Blockchain.VerifyChain()
}

func PrintBlock(block *Block) {
	fmt.Println("Block ID:", block.CurrentHash)
	fmt.Println("Timestamp:", block.Nonce)
	fmt.Println("Nonce:", block.Nonce)
	fmt.Println("Previous Block Hash:", block.PreviousHash)
	fmt.Println("Transactions:")
	for _, transaction := range block.Transaction {
		fmt.Printf("\tTransaction ID: %s\n", transaction.TransactionID)
		fmt.Printf("\tSender: %s\n", transaction.SenderBlockchainAddress)
		fmt.Printf("\tRecipient: %s\n", transaction.RecipientBlockchainAddress)
		fmt.Printf("\tValue: %.2f\n", transaction.Value)
	}
}

func (bc *Blockchain) Print() {
	for index, block := range bc.Chain {
		fmt.Println("Block:", index+1)
		PrintBlock(block)
		fmt.Println()
	}
}

func CalculateHash(stringToHash string) string {
	hash := sha256.Sum256([]byte(stringToHash))
	return fmt.Sprintf("%x", hash)
}

func NewBlockchain() *Blockchain {
	blockchain := &Blockchain{}
	return blockchain
}

func (bc *Blockchain) VerifyChain() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		if currentBlock.CurrentHash != CalculateHash(currentBlock.PreviousHash+strconv.Itoa(currentBlock.Nonce)+TransactionToJSON(currentBlock.Transaction)) {
			return false
		}
		if currentBlock.PreviousHash != previousBlock.CurrentHash {
			return false
		}
	}

	return true
}

func (bc *Blockchain) GetHeadBlockCurrentHash() string {
	if len(bc.Chain) > 0 {
		return bc.Chain[len(bc.Chain)-1].CurrentHash
	}
	return ""
}

func NewNode(difficulty int, minerAddress string) *Node {
	return &Node{
		Blockchain:   NewBlockchain(),
		Difficulty:   difficulty,
		MinerAddress: minerAddress,
		Reward:       10.0,
	}
}

func (n *Node) MineBlock() {
	var (
		transactions = append([]*Transaction{}, n.Transactions...)
		previousHash = n.Blockchain.GetHeadBlockCurrentHash()
		nonce        = 0
	)

	for {
		hash := CalculateHash(previousHash + strconv.Itoa(nonce) + TransactionToJSON(transactions))
		if hash[:n.Difficulty] == strings.Repeat("0", n.Difficulty) {
			fmt.Printf("Block mined with nonce: %d, Hash: %s\n", nonce, hash)
			rewardTransaction := NewTransaction("Reward", n.MinerAddress, n.Reward)
			transactions = append(transactions, rewardTransaction)
			n.Blockchain.AddBlock(transactions, nonce, previousHash)
			return
		}
		nonce++
	}
}
