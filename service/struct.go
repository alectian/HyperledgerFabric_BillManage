package service

import "github.com/kongyixueyuan.com/bill/blockchain"

type FabricSetupService struct {
	Fabric *blockchain.FabricSetup
}

type Bill struct {

	BillInfoID string `json:"bill_info_id"` // 票据号码
	BillInfoAmt string `json:"bill_info_amt"` // 票据金额
	BillInfoType string `json:"bill_info_type"` // 票据类型

	BillInfoIsseDate string `json:"bill_info_isse_date"` // 出票日期
	BillInfoDueDate string `json:"bill_info_due_date"` // 失效日期

	DrwrAcct string `json:"drwr_acct"` // 出票人
	DrwrCmID string `json:"drwr_cm_id"`

	AccptrAcct string `json:"accptr_acct"` // 承兑人
	AccptrCmID string `json:"accptr_cm_id"`

	PyeeAcct string  `json:"pyee_acct"` // 收款人
	PyeeCmID string `json:"pyee_cm_id"` // 收款人证件号码

	HoldrAcct string `json:"holdr_acct"` // 持票人
	HoldrCmId string `json:"holdr_cm_id"`

	WaitEndorseAcct string	`json:"wait_endorse_acct"` //等待背书人
	WaitEndorseCmId string `json:"wait_endorse_cm_id"`


	RejectEndorseAcct string `json:"reject_endorse_acct"` // 拒绝背书人
	RejectEndorseCmID string `json:"reject_endorse_cm_id"`

	BillStatus string `json:"bill_status"` // 票据状态
	Historys []HistoryItem `json:"historys"`
}

type HistoryItem struct {
	IxId string // 交易记录
	Bill Bill
}
