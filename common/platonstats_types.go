package common

import (
	"encoding/json"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/log"
)

type BlockType uint8
type NodeID [512 / 8]byte

type Input []byte

var nodeIdT = reflect.TypeOf(NodeID{})

// MarshalText returns the hex representation of a.
func (a NodeID) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *NodeID) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("common.NodeID", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *NodeID) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(nodeIdT, input, a[:])
}

var inputT = reflect.TypeOf(Input{})

// MarshalText returns the hex representation of a.
func (a Input) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *Input) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("common.Input", input, *a)
	/*aa := make(Input, len(input))
	err := hexutil.UnmarshalFixedText("common.Input", input, aa)
	if err != nil {
		return errors.New(" Unmarshal Input error")
	}
	a = &aa
	return nil*/
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Input) UnmarshalJSON(input []byte) error {
	//string(input)="0x0102030405", so, firstly remove the "", and then, to decode it.
	hexBytes, err := hexutil.Decode(string(input[1 : len(input)-1]))
	if err != nil {
		return err
	}
	aa := make(Input, len(hexBytes))
	err = hexutil.UnmarshalFixedJSON(inputT, input, aa)
	if err != nil {
		return err
	}
	a = &aa
	return nil
}

const (
	GenesisBlock BlockType = iota
	GeneralBlock
	ConsensusBeginBlock
	ConsensusElectionBlock
	ConsensusEndBlock
	EpochBeginBlock
	EpochEndBlock
	EndOfYear
)

type EmbedTransferTx struct {
	From   Address `json:"from,omitempty"`
	To     Address `json:"to,omitempty"`
	Amount uint64  `json:"amount"`
}

type EmbedContractTx struct {
	From            Address `json:"from,omitempty"`
	ContractAddress Address `json:"contractAddress,omitempty"`
	Input           Input   `json:"input,omitempty"`
}

type GenesisData struct {
	AllocItemList []*AllocItem `json:"allocItemList,omitempty"`
}
type AllocItem struct {
	Address Address `json:"address,omitempty"`
	Amount  uint64  `json:"amount"`
}

func (g *GenesisData) AddAllocItem(address Address, amount uint64) {
	//todo: test
	g.AllocItemList = append(g.AllocItemList, &AllocItem{Address: address, Amount: amount})
}

type AdditionalIssuanceData struct {
	AdditionalNo     uint32          `json:"additionalNo"`               //增发周期
	AdditionalBase   uint64          `json:"additionalBase"`             //增发基数
	AdditionalRate   uint16          `json:"additionalRate"`             //增发比例 单位：万分之一
	AdditionalAmount uint64          `json:"additionalAmount"`           //增发金额
	IssuanceItemList []*IssuanceItem `json:"issuanceItemList,omitempty"` //增发分配
}

type IssuanceItem struct {
	Address Address `json:"address,omitempty"` //增发金额分配地址
	Amount  uint64  `json:"amount"`            //增发金额
}

func (d *AdditionalIssuanceData) AddIssuanceItem(address Address, amount uint64) {
	//todo: test
	d.IssuanceItemList = append(d.IssuanceItemList, &IssuanceItem{Address: address, Amount: amount})
}

type RewardData struct {
	BlockRewardAmount   uint64           `json:"blockRewardAmount"`           //出块奖励
	StakingRewardAmount uint64           `json:"stakingRewardAmount"`         //一结算周期内所有101节点的质押奖励
	CandidateInfoList   []*CandidateInfo `json:"candidateInfoList,omitempty"` //备选节点信息
}

type CandidateInfo struct {
	NodeID       NodeID  `json:"nodeId,omitempty"`       //备选节点ID
	MinerAddress Address `json:"minerAddress,omitempty"` //备选节点的矿工地址（收益地址）
}

type ZeroSlashingItem struct {
	NodeID         NodeID `json:"nodeId,omitempty"` //备选节点ID
	SlashingAmount uint64 `json:"slashingAmount"`   //0出块处罚金(从质押金扣)
}

type DuplicatedSignSlashingSetting struct {
	PenaltyRatioByValidStakings uint32 `json:"penaltyRatioByValidStakings"` //unit:1%%		//罚金 = 有效质押 & PenaltyRatioByValidStakings / 10000
	RewardRatioByPenalties      uint32 `json:"rewardRatioByPenalties"`      //unit:1%		//给举报人的赏金=罚金 * RewardRatioByPenalties / 100
}

type UnstakingRefundItem struct {
	NodeID        NodeID      `json:"nodeId,omitempty"`      //备选节点ID
	NodeAddress   NodeAddress `json:"nodeAddress,omitempty"` //备选节点地址
	RefundEpochNo uint64      `json:"refundEpochNo"`         //解除质押,资金真正退回的结算周期（此周期最后一个块的endBlocker里
}

type RestrictingReleaseItem struct {
	DestAddress   Address `json:"destAddress,omitempty,omitempty"` //释放地址
	ReleaseAmount uint64  `json:"releaseAmount"`                   //释放金额
}

var ExeBlockDataCollector = make(map[uint64]*ExeBlockData)

func PopExeBlockData(blockNumber uint64) *ExeBlockData {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		delete(ExeBlockDataCollector, blockNumber)
		return exeBlockData
	}
	return nil
}

func InitExeBlockData(blockNumber uint64) {
	exeBlockData := &ExeBlockData{
		ZeroSlashingItemList:       make([]*ZeroSlashingItem, 0),
		UnstakingRefundItemList:    make([]*UnstakingRefundItem, 0),
		RestrictingReleaseItemList: make([]*RestrictingReleaseItem, 0),
		EmbedTransferTxMap:         make(map[Hash][]*EmbedTransferTx),
		EmbedContractTxMap:         make(map[Hash][]*EmbedContractTx),
	}

	ExeBlockDataCollector[blockNumber] = exeBlockData
}

func GetExeBlockData(blockNumber uint64) *ExeBlockData {
	return ExeBlockDataCollector[blockNumber]
}

type ExeBlockData struct {
	AdditionalIssuanceData        *AdditionalIssuanceData        `json:"additionalIssuanceData,omitempty"`
	RewardData                    *RewardData                    `json:"rewardData,omitempty"`
	ZeroSlashingItemList          []*ZeroSlashingItem            `json:"zeroSlashingItemList,omitempty"`
	DuplicatedSignSlashingSetting *DuplicatedSignSlashingSetting `json:"duplicatedSignSlashingSetting,omitempty"`
	UnstakingRefundItemList       []*UnstakingRefundItem         `json:"unstakingRefundItemList,omitempty"`
	RestrictingReleaseItemList    []*RestrictingReleaseItem      `json:"restrictingReleaseItemList,omitempty"`
	EmbedTransferTxMap            map[Hash][]*EmbedTransferTx    `json:"embedTransferTxMap,omitempty"` //一个显式交易引起的内置转账交易：一般有两种情况：1是部署，或者调用合约时，带上了value，则这个value会转账给合约地址；2是调用合约，合约内部调用transfer()函数完成转账
	EmbedContractTxMap            map[Hash][]*EmbedContractTx    `json:"embedContractTxMap,omitempty"` //一个显式交易引起的内置合约交易。这个显式交易显然也是个合约交易，在这个合约里，又调用了其他合约（包括内置合约）
}

func CollectUnstakingRefundItem(blockNumber uint64, nodeId NodeID, nodeAddress NodeAddress, refundEpochNo uint64) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectUnstakingRefundItem", "blockNumber", blockNumber, "nodeId", Bytes2Hex(nodeId[:]), "nodeAddress", nodeAddress.Hex(), "refundEpochNo", refundEpochNo)
		exeBlockData.UnstakingRefundItemList = append(exeBlockData.UnstakingRefundItemList, &UnstakingRefundItem{NodeID: nodeId, NodeAddress: nodeAddress, RefundEpochNo: refundEpochNo})
	}
}

func CollectRestrictingReleaseItem(blockNumber uint64, destAddress Address, releaseAmount uint64) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectRestrictingReleaseItem", "blockNumber", blockNumber, "destAddress", destAddress, "releaseAmount", releaseAmount)
		exeBlockData.RestrictingReleaseItemList = append(exeBlockData.RestrictingReleaseItemList, &RestrictingReleaseItem{DestAddress: destAddress, ReleaseAmount: releaseAmount})
	}
}

func CollectRewardData(blockNumber uint64, rewardData *RewardData) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectRewardData", "blockNumber", blockNumber, "rewardData", rewardData.BlockRewardAmount)
		exeBlockData.RewardData = rewardData
	}
}

func CollectDuplicatedSignSlashingSetting(blockNumber uint64, penaltyRatioByValidStakings, rewardRatioByPenalties uint32) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectDuplicatedSignSlashingSetting", "blockNumber", blockNumber, "penaltyRatioByValidStakings", penaltyRatioByValidStakings, "rewardRatioByPenalties", rewardRatioByPenalties)
		if exeBlockData.DuplicatedSignSlashingSetting != nil {
			//在同一个区块中，只要设置一次即可
			exeBlockData.DuplicatedSignSlashingSetting = &DuplicatedSignSlashingSetting{PenaltyRatioByValidStakings: penaltyRatioByValidStakings, RewardRatioByPenalties: rewardRatioByPenalties}
		}
	}
}

func CollectZeroSlashingItem(blockNumber uint64, zeroSlashingItemList []*ZeroSlashingItem) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		json, _ := json.Marshal(zeroSlashingItemList)
		log.Debug("CollectZeroSlashingItem", "blockNumber", blockNumber, "zeroSlashingItemList", string(json))
		exeBlockData.ZeroSlashingItemList = zeroSlashingItemList
	}
}

func CollectEmbedTransferTx(blockNumber uint64, txHash Hash, from, to Address, amount uint64) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectEmbedTransferTx", "blockNumber", blockNumber, "txHash", txHash.Hex(), "from", from.Bech32(), "to", to.Bech32(), "amount", amount)
		exeBlockData.EmbedTransferTxMap[txHash] = append(exeBlockData.EmbedTransferTxMap[txHash], &EmbedTransferTx{From: from, To: to, Amount: amount})
	}
}

func CollectEmbedContractTx(blockNumber uint64, txHash Hash, from, contractAddress Address, input []byte) {
	if exeBlockData, ok := ExeBlockDataCollector[blockNumber]; ok && exeBlockData != nil {
		log.Debug("CollectEmbedContractTx", "blockNumber", blockNumber, "txHash", txHash.Hex(), "contractAddress", from.Bech32(), "input", Bytes2Hex(input))
		exeBlockData.EmbedContractTxMap[txHash] = append(exeBlockData.EmbedContractTxMap[txHash], &EmbedContractTx{From: from, ContractAddress: contractAddress, Input: input})
	}
}
