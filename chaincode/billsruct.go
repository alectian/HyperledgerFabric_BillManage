package main


//import "github.com/hyperledger/fabric/core/chaincode/shim"

const (

	BillInfo_state_NewPublish = "NewPublish" // 票据新发布状态
	BillInfo_state_EndorseWaitSign = "EndorseWaitSign" // 等待背书
	BillInfo_state_EndorseSigned = "EndorseSigned" // 背书成功
	BillInfo_state_EndorseReject = "EndorseReject"// 拒绝背书
)

const (
	Bill_profix = "Bill_" // 票据号码的前缀
	queryBillsIndexName = "holdrCmID~BillNo" // 映射名称（方便查询）
	querywaitBillsIndexName = "WaitEndorseCmId~BillNo"
)

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

