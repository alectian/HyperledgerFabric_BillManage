package main

import (
	"fmt"
	"github.com/kongyixueyuan.com/bill/blockchain"
	"github.com/kongyixueyuan.com/bill/service"
	"os"
	"encoding/json"
)

func main() {

	// 定义SDK属性
	fSetup := blockchain.FabricSetup{
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		ConfigFile: "config.yaml",

		// 通道相关 
		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/kongyixueyuan.com/bill/fixtures/artifacts/channel.tx",

		// 链码相关参数
		ChaincodeID:     "billcc",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/kongyixueyuan.com/bill/chaincode/",

		// 指定用户
		UserName: "User1",
	}

	// 初始化SDK
	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("无法初始化Fabric SDK: %v\n", err)
	}

	// 安装及初始化链码
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("无法安装及实例化链码: %v\n", err)
	}

	// ==========================测试开始 ==============================
	// 发布票据
	bill := service.Bill{
		BillInfoID:       "BOC10000001",
		BillInfoAmt:      "222",
		BillInfoType:     "111",
		BillInfoIsseDate: "20180501",
		BillInfoDueDate:  "20180503",
		DrwrCmID:         "111",
		DrwrAcct:         "111",
		AccptrCmID:       "111",
		AccptrAcct:       "111",
		PyeeCmID:         "111",
		PyeeAcct:         "111",
		HoldrCmId:        "BCMID",
		HoldrAcct:        "B公司",
	}

	bill2 := service.Bill{
		BillInfoID:       "BOC10000002",
		BillInfoAmt:      "222",
		BillInfoType:     "111",
		BillInfoIsseDate: "20180501",
		BillInfoDueDate:  "20180503",
		DrwrCmID:         "111",
		DrwrAcct:         "111",
		AccptrCmID:       "111",
		AccptrAcct:       "111",
		PyeeCmID:         "111",
		PyeeAcct:         "111",
		HoldrCmId:        "BCMID",
		HoldrAcct:        "B公司",
	}

	fsService := new(service.FabricSetupService)

	fsService.Fabric = &fSetup

	// 发布票据
	resp, err := fsService.SaveBill(bill)
	if err != nil {
		fmt.Printf("发布票据失败: %v\n", err)
	} else {
		fmt.Println("发布票据成功 =========>交易编号：" + resp)
	}
	// 发布票据
	resp, err = fsService.SaveBill(bill2)
	if err != nil {
		fmt.Printf("发布票据失败: %v\n", err)
	} else {
		fmt.Println("发布票据成功 =========>交易编号：" + resp)
	}

	//查询票据列表
	fresp,err := fsService.FindBills(bill.HoldrCmId)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	var billArr []service.Bill

	err = json.Unmarshal(fresp,&billArr)
	if err != nil {
		fmt.Printf("查询票据列表失败，err:%v",err)
		return
	}
	fmt.Println("查询票据列表成功\n",billArr)


	// 发起背书
	result,err := fsService.Endorse(bill.BillInfoID,"mark","123456")

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	result,err = fsService.Endorse(bill2.BillInfoID,"mark","123456")

	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Println("发起背书成功\n",result)

	//查询票据详情
	qresp,err := fsService.QueryBillByNo(bill.BillInfoID)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	var billHistory []service.HistoryItem

	err = json.Unmarshal(qresp,&billHistory)
	if err != nil {
		fmt.Printf("查询票据详情失败，err:%v",err)
		return
	}
	fmt.Println("查询票据详情成功\n",billHistory)

	// 查询待背书列表
	qwresp,err := fsService.QueryWaitBills("123456")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	var waitEndorseBills []service.Bill
	err = json.Unmarshal(qwresp,&waitEndorseBills)
	if err != nil {
		fmt.Printf("查询待背书列表失败，err:%v",err)
		return
	}
	fmt.Println("查询待背书列表成功\n",waitEndorseBills)

	// 背书签名
	result,err = fsService.Accept(bill.BillInfoID,"jerry","12345")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Println("背书签名成功\n",result)

	//拒绝背书
	result,err = fsService.Reject(bill2.BillInfoID,"ginny","12345")
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Println("已拒绝背书签名\n",result)

	//查询票据详情
	qresp,err = fsService.QueryBillByNo(bill.BillInfoID)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	var billHistory1 []service.HistoryItem

	err = json.Unmarshal(qresp,&billHistory1)
	if err != nil {
		fmt.Printf("查询票据详情失败，err:%v",err)
		return
	}
	fmt.Println("查询票据详情成功\n",billHistory1)

	//查询票据详情
	qresp,err = fsService.QueryBillByNo(bill2.BillInfoID)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	var billHistory2 []service.HistoryItem

	err = json.Unmarshal(qresp,&billHistory2)
	if err != nil {
		fmt.Printf("查询票据详情失败，err:%v",err)
		return
	}
	fmt.Println("查询票据详情成功\n",billHistory2)
}