package blockchain

import (
	"bytes"
	"errors"
	"reflect"
	"time"

	ygg "github.com/DE-labtory/yggdrasill/common"
)

// ErrHashCalculationFailed 변수는 Hash 계산 중 발생한 에러를 정의한다.
var ErrHashCalculationFailed = errors.New("Hash Calculation Failed Error")
var ErrInsufficientFields = errors.New("Previous seal or transaction list seal is not set")
var ErrEmptyTxList = errors.New("Empty TxList")

type Validator = ygg.Validator


// DefaultValidator 객체는 Validator interface를 구현한 객체.
type DefaultValidator struct{}

// ValidateSeal 함수는 원래 Seal 값과 주어진 Seal 값(comparisonSeal)을 비교하여, 올바른지 검증한다.
func (t *DefaultValidator) ValidateSeal(seal []byte, comparisonBlock ygg.Block) (bool, error) {

	comparisonSeal, error := t.BuildSeal(comparisonBlock.GetTimestamp(), comparisonBlock.GetPrevSeal(), comparisonBlock.GetTxSeal(), comparisonBlock.GetCreator())

	if error != nil {
		return false, error
	}
	return bytes.Compare(seal, comparisonSeal) == 0, nil
}

// ValidateTxSeal 함수는 주어진 Transaction 리스트에 따라 주어진 transaction Seal을 검증함.
func (t *DefaultValidator) ValidateTxSeal(txSeal [][]byte, txList []Transaction) (bool, error) {

	if isEmpty(txList) {
		return true, nil
	}

	leafNodeList, err := convertToLeafNodeList(txList)
	if err != nil {
		return false, err
	}

	tree, err := buildTree(leafNodeList, leafNodeList)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(txSeal, tree), nil
}

func isEmpty(txList []Transaction) bool {
	if len(txList) == 0 {
		return true
	}
	return false
}

func convertToLeafNodeList(txList []Transaction) ([][]byte, error) {
	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, err := tx.CalculateSeal()
		if err != nil {
			return nil, err
		}

		leafNodeList = append(leafNodeList, leafNode)
	}

	if len(leafNodeList)%2 != 0 {
		leafNodeList = append(leafNodeList, leafNodeList[len(leafNodeList)-1])
	}

	return leafNodeList, nil
}

// ValidateTransaction 함수는 주어진 Transaction이 이 txSeal에 올바로 있는지를 확인한다.
func (t *DefaultValidator) ValidateTransaction(txSeal [][]byte, transaction Transaction) (bool, error) {
	hash, error := transaction.CalculateSeal()
	if error != nil {
		return false, error
	}

	index := -1
	for i, h := range txSeal {
		if bytes.Compare(h, hash) == 0 {
			index = i
		}
	}

	if index == -1 {
		return false, nil
	}

	var siblingIndex, parentIndex int
	for index > 0 {
		var isLeft bool
		if index%2 == 0 {
			siblingIndex = index - 1
			parentIndex = (index - 1) / 2
			isLeft = false
		} else {
			siblingIndex = index + 1
			parentIndex = index / 2
			isLeft = true
		}

		var parentHash []byte
		if isLeft {
			parentHash = calculateIntermediateNodeHash(txSeal[index], txSeal[siblingIndex])
		} else {
			parentHash = calculateIntermediateNodeHash(txSeal[siblingIndex], txSeal[index])
		}

		if bytes.Compare(parentHash, txSeal[parentIndex]) != 0 {
			return false, nil
		}

		index = parentIndex
	}

	return true, nil
}

// BuildSeal 함수는 block 객체를 받아서 Seal 값을 만들고, Seal 값을 반환한다.
// 인풋 파라미터의 block에 자동으로 할당해주지는 않는다.
func (t *DefaultValidator) BuildSeal(timeStamp time.Time, prevSeal []byte, txSeal [][]byte, creator string) ([]byte, error) {
	timestamp, err := timeStamp.MarshalText()
	if err != nil {
		return nil, err
	}

	if prevSeal == nil || txSeal == nil || creator == "" {
		return nil, ErrInsufficientFields
	}
	var rootHash []byte
	if len(txSeal) == 0 {
		rootHash = make([]byte, 0)
	} else {
		rootHash = txSeal[0]
	}
	combined := append(prevSeal, rootHash...)
	combined = append(combined, timestamp...)

	seal := calculateHash(combined)
	return seal, nil
}

// BuildTxSeal 함수는 Transaction 배열을 받아서 TxSeal을 생성하여 반환한다.
func (t *DefaultValidator) BuildTxSeal(txList []Transaction) ([][]byte, error) {
	if len(txList) == 0 {
		return nil, ErrEmptyTxList
	}

	leafNodeList := make([][]byte, 0)

	for _, tx := range txList {
		leafNode, error := tx.CalculateSeal()
		if error != nil {
			return nil, error
		}

		leafNodeList = append(leafNodeList, leafNode)
	}

	// leafNodeList의 개수는 짝수개로 맞춤. (홀수 일 경우 마지막 Tx를 중복 저장.)
	if len(leafNodeList)%2 != 0 {
		leafNodeList = append(leafNodeList, leafNodeList[len(leafNodeList)-1])
	}

	tree, error := buildTree(leafNodeList, leafNodeList)
	if error != nil {
		return nil, error
	}

	// DefaultValidator 는 Merkle Tree의 루트노드(tree[0])를 Proof로 간주함
	return tree, nil
}

func buildTree(nodeList [][]byte, fullNodeList [][]byte) ([][]byte, error) {
	intermediateNodeList := make([][]byte, 0)
	for i := 0; i < len(nodeList); i += 2 {
		leftIndex, rightIndex := i, i+1

		if i+1 == len(nodeList) {
			rightIndex = i
		}

		leftNode, rightNode := nodeList[leftIndex], nodeList[rightIndex]

		intermediateNode := calculateIntermediateNodeHash(leftNode, rightNode)

		intermediateNodeList = append(intermediateNodeList, intermediateNode)

		if len(nodeList) == 2 {
			return append(intermediateNodeList, fullNodeList...), nil
		}
	}

	newFullNodeList := append(intermediateNodeList, fullNodeList...)

	return buildTree(intermediateNodeList, newFullNodeList)
}

func calculateIntermediateNodeHash(leftHash []byte, rightHash []byte) []byte {
	combinedHash := append(leftHash, rightHash...)

	return calculateHash(combinedHash)
}
