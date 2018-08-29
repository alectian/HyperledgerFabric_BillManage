package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
	//"github.com/coreos/etcd/pkg/wait"
)


type BillChaincode struct {

}

func (b *BillChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (b *BillChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	// 获取用户意图
	fun,args := stub.GetFunctionAndParameters()
	if fun == "issue" {
		return b.issue(stub,args)
	}
	if fun == "queryBills" {
		return b.queryBills(stub,args)
	}
	if fun == "queryBillByNo" {
		return b.queryBillByNo(stub,args)
	}
	if fun == "endorse" {
		return b.endorse(stub,args)
	}
	if fun == "queryWaitBills" {
		return b.queryWaitBills(stub,args)
	}

	if fun == "accept" {
		return b.accept(stub,args)
	}

	if fun == "reject" {
		return b.reject(stub,args)
	}

	// return shim.Erorr
	respMsg, err := GetMsgString(1,"指定的函数名错误")
	if err != nil {
		fmt.Println("获取消息失败")
	}
	return shim.Error(respMsg)
}

func main() {
	err := shim.Start(new(BillChaincode))
	if err != nil {
		fmt.Println("启动链码失败")
	}
}



// 发布票据 -c '{"Args":["issue","Bill"]}'
func (b *BillChaincode) issue(stub shim.ChaincodeStubInterface,args []string) peer.Response {

	if len(args) != 1{
		return shim.Error("保存票据失败，参数错误")
	}
	var bill Bill

	err := json.Unmarshal([]byte(args[0]),&bill)
	if err != nil {
		return shim.Error("反序列化票据失败")
	}

	bill.BillStatus = BillInfo_state_NewPublish
	fmt.Println("开始保存票据")
	fmt.Println("bill:",bill)

	data,err :=json.Marshal(bill)
	if err != nil {
		fmt.Println("序列化票据失败")
		return shim.Error("序列化票据失败")
	}

	// 票据编号查询
	bi,err := stub.GetState(bill.BillInfoID)
	if bi == nil {
		// 按票据编号存
		err = stub.PutState(bill.BillInfoID,data)
		if err != nil {
			return shim.Error("按票据编号存票据时失败")
		}
	}else{
		return shim.Error("票据编号已经存在")
	}

	//创建复合键，方便查询，以免键的重复，导致存取失败
	HoldrAcctBillInfoIndexName,err := stub.CreateCompositeKey(queryBillsIndexName,[]string{bill.HoldrCmId,bill.BillInfoID})
	if err != nil {
		return shim.Error("创建复合键失败")
	}
	// 如果保存的复合key时指定的value为nil，会导致查询不到
	err = stub.PutState(HoldrAcctBillInfoIndexName,[]byte("2"))
	if err != nil {
		return shim.Error("按复合键存票据时失败")
	}

	fmt.Println("保存票据成功")
	return shim.Success(nil)
}

func (b *BillChaincode) find(stub shim.ChaincodeStubInterface, HoldrCmId string) ([]byte,error) {

	fmt.Println("开始查询票据列表")

	// 查询所有的于HoldCmid相关的复合键
	billiterator, err := stub.GetStateByPartialCompositeKey(queryBillsIndexName,[]string{HoldrCmId})
	if err != nil {
		return nil,err
	}
	defer billiterator.Close()
	var bill Bill
	var billSli []Bill
	// 迭代处理
	for billiterator.HasNext(){
		//k:compositekey;v:[]byte("2")
		kv,err := billiterator.Next()
		if err != nil {
			fmt.Printf("handle billiterator failed err：%v",err)
			return nil,err
		}
		_,keys,err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			fmt.Printf("split keys failed err:%v",err)
			return nil,err
		}
		billobj,err := stub.GetState(keys[1])
		//fmt.Println("billobj",billobj)
		if err != nil {
			fmt.Printf("get bill by holdrcmid faild err:%v",err)
			return nil,err
		}
		err = json.Unmarshal(billobj,&bill)
		if err != nil {
			fmt.Printf("unmarshal billobj failed")
			return nil,err
		}
		billSli = append(billSli,bill)
	}
		fmt.Println("票据列表：",billSli)
		ret,err := json.Marshal(billSli)
		if err != nil {
			fmt.Printf("marshal billslic failed err:%v",err)
			return nil,err
		}
	return ret,nil
}


// 查询票据详情by bill.BillInfoID -c '{"Args":["queryBillByNo","bill.BillInfoID"]}'
func (b *BillChaincode) queryBillByNo(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("开始根据票据编号查询票据")

	if len(args) != 1 {
		msg,_:= GetMsgString(1,"参数错误，请仅输入票据编号")
		return shim.Error(msg)
	}
	//billData,err := stub.GetState(args[0])
	//
	//if err != nil {
	//	msg,_ := GetMsgString(1,"根据票据编号查询数据失败")
	//	return shim.Error(msg)
	//}

	var history HistoryItem
	var historyDatas []HistoryItem
	var bill Bill

	iteratedObj,err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return  shim.Error("根据票据持票人证件号码查询票据失败")
	}
	for iteratedObj.HasNext() {

		historyData,err :=iteratedObj.Next()
		if err != nil {
			return shim.Error("遍历失败")
		}
		if historyData.Value == nil {
			var emptyBill Bill
			history.Bill = emptyBill
		}
		err = json.Unmarshal(historyData.Value,&bill)
		if err != nil {
			return shim.Error("反序列化失败")
		}

		history.Bill = bill
		history.IxId = historyData.TxId
		historyDatas = append(historyDatas,history)
	}

	bill.Historys = historyDatas

	ret,err := json.Marshal(bill.Historys)
	if err != nil {
		return shim.Error("查询票据列表时，序列化失败")
	}
	fmt.Println("根据票据编号查询票据成功")
	return shim.Success(ret)
}



// 查询票据 -c '{"Args":["queryBills","HoldrCmId"]}'
func (b *BillChaincode) queryBills(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("参数错误，请仅输入票据持票人证件号码")
	}

	bills,err := b.find(stub,args[0])
	if err != nil {
		return  shim.Error("根据票据持票人证件号码查询票据列表失败")
	}

	fmt.Println("查询票据列表成功")
	return shim.Success(bills)
}


// 根据票据编号查询票据 -c '{"Args":["bill.BillInfoId"]}'
func (b *BillChaincode) queryBill(stub shim.ChaincodeStubInterface,args []string) (*Bill,bool) {

	var bill Bill

	bi,err := stub.GetState(args[0])
	if err != nil {
		return nil,false
	}
	err = json.Unmarshal(bi,&bill)
	if err != nil {
		return nil,false
	}
	return &bill, true
}

// 发起背书 -'{"Args":["endorse","billInfoId","WaitEndorseAcct","WaitEndorseCmId"]}'
func (b *BillChaincode) endorse(stub shim.ChaincodeStubInterface,args []string) peer.Response {

	// 1.根据票据号码查询票据并把票据状态改为待背书 2.以背书人的证件号码为key，将票据存入数据库

	// 先根据billInfoId查询票据
	waitEndorsebill,bl := b.queryBill(stub,args)
	if !bl {
		return shim.Error("背书时查询票据失败")
	}

	if waitEndorsebill.HoldrCmId == args[2] {
		return shim.Error("背书人不能是当前持票人")
	}

	billiterator,err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error("背书时查询票据历史失败")
	}
	defer  billiterator.Close()

	var billHisData Bill

	//背书人不能是历史持证人
	for billiterator.HasNext(){
		kvs,err := billiterator.Next()
		if err != nil {
			return shim.Error("handle bill iterator failed")
		}
		if kvs.Value == nil {
			continue
		}
		 err = json.Unmarshal(kvs.Value,&billHisData)
		 if err != nil {
			 return shim.Error("unmarshal billHisData failed")
		 }

		if billHisData.HoldrCmId == args[2] {
			return shim.Error("endorseCmID  was HoldCmId in historyTxId ")
		}
	}

	waitEndorsebill.RejectEndorseAcct = ""
	waitEndorsebill.RejectEndorseCmID = ""
	waitEndorsebill.BillStatus = BillInfo_state_EndorseWaitSign
	waitEndorsebill.WaitEndorseAcct = args[1]
	waitEndorsebill.WaitEndorseCmId = args[2]

	bill,err := json.Marshal(waitEndorsebill)
	if err != nil {
		return shim.Error("endorse bill failed, marshal bill failed")
	}

	err = stub.PutState(args[0],bill)
	if err != nil {
		return shim.Error("endorse bill failed, put bill failed")
	}

	compkey,err := stub.CreateCompositeKey(querywaitBillsIndexName,[]string{waitEndorsebill.WaitEndorseCmId,waitEndorsebill.BillInfoID})
	//compkey,err := stub.CreateCompositeKey(IndexName,[]string{waitEndorsebill.BillInfoID,waitEndorsebill.WaitEndorseCmId})
	fmt.Println(compkey)
	if err != nil {
		return shim.Error("create composite key failed")
	}

	err = stub.PutState(compkey,[]byte("3"))
	if err != nil {
		return shim.Error("composite key put failed")
	}

	fmt.Println("success")
	return shim.Success([]byte("success"))
}


// 查询待背书票据 -c '{"Args":["queryWaitBills","WaitEndorseCmI"]}'
func (b *BillChaincode) queryWaitBills(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 1.通过背书的证件号码查询所有的票据
	if len(args) != 1 {
		return shim.Error("查询待背书列表失败，仅输入待背书人证件号码")
	}
	fmt.Println("WaitEndorseCmI",args[0])
	iterator,err := stub.GetStateByPartialCompositeKey(querywaitBillsIndexName,[]string{args[0]})
	if err != nil {
		shim.Error("查询待背书列表失败，获取相关链复合键失败")
	}

	var waitEndorseBill Bill
	var waitEndorseBills []Bill

	for iterator.HasNext(){
		kvs,err := iterator.Next()
		if err != nil {
			return shim.Error("查询待背书列表失败,处理迭代对象出错")
		}
		_,keys,err := stub.SplitCompositeKey(kvs.Key)
		if err != nil {
			return shim.Error("查询待背书列表失败,分解复键时出错")
		}
		bill,err := stub.GetState(keys[1])
		if err != nil {
			return shim.Error("查询待背书列表失败，查询票据出错")
		}
		err = json.Unmarshal(bill,&waitEndorseBill)
		if err != nil {
			return shim.Error("查询待背书列表失败，反序列化时出错")
		}
		waitEndorseBills = append(waitEndorseBills,waitEndorseBill)
	}

	result,err := json.Marshal(waitEndorseBills)
	fmt.Println("查询待背书票据列表成功")
	fmt.Println("waitEndorseBills:",waitEndorseBills)
	fmt.Println("查询待背书票据列表成功")
	return shim.Success(result)
}

// 背书票据 -c "{"Args",["billInfoId","AccptrAcct",AccptrCmId"]]'
func (b *BillChaincode) accept(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("输入的参数错误")
	}

	bill,bl := b.queryBill(stub , args)
	if !bl {
		return shim.Error("accept bill failed, get bill failed")
	}

	// 创建复合键，并删除
	compositekey,err := stub.CreateCompositeKey(queryBillsIndexName,[]string{bill.HoldrCmId,bill.BillInfoID})
	if err != nil {
		return shim.Error("accept bill failed,create composite key failed")
	}

	err = stub.DelState(compositekey)
	if err != nil {
		return shim.Error("accept bill failed,del composite key failed")
	}

	bill.HoldrCmId = args[1]
	bill.HoldrAcct = args[2]
	bill.WaitEndorseAcct = ""
	bill.WaitEndorseCmId = ""
	bill.BillStatus = BillInfo_state_EndorseSigned

	billData,err := json.Marshal(bill)
	if err != nil {
		return shim.Error("accept bill failed,marshal bill failed")
	}

	err = stub.PutState(bill.BillInfoID,billData)
	if err != nil {
		return shim.Error("accept bill failed,put bill failed")
	}

	fmt.Println("背书成功")
	return shim.Success(nil)
}

//拒绝背书 -c "{"Args",["billInfoId","RejectrAcct",RejectrCmId"]]'

func (b *BillChaincode) reject(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("输入的参数错误")
	}

	bill,bl := b.queryBill(stub , args)

	if !bl {
		return shim.Error("reject bill failed, get bill failed")
	}

	//create composite key and del this composite key belong to en
	compostitekey,err := stub.CreateCompositeKey(querywaitBillsIndexName,[]string{bill.BillInfoID,bill.WaitEndorseCmId})
	if err != nil {
		return shim.Error("reject sign bill failed, create composite key failed ")
	}
	err = stub.DelState(compostitekey)
	if err != nil {
		return shim.Error("reject sign bill failed, del composite key failed ")
	}

	bill.WaitEndorseAcct = ""
	bill.WaitEndorseCmId = ""
	bill.RejectEndorseAcct = args[1]
	bill.RejectEndorseCmID = args[2]
	bill.BillStatus = BillInfo_state_EndorseReject

	billData,err := json.Marshal(bill)
	if err != nil {
		return shim.Error("reject bill failed,marshal bill failed")
	}

	err = stub.PutState(bill.BillInfoID,billData)
	if err != nil {
		return shim.Error("reject bill failed,put bill failed")
	}

	fmt.Println("背书请求已拒绝")
	return shim.Success(nil)

}


