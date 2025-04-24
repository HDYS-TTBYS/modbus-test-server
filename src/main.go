package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/simonvetter/modbus"
)

const (
	MINUS_ONE int16 = -1
)

type exampleHandler struct {
	lock sync.RWMutex

	coils            [16][10000]bool
	discreteInputs   [16][10000]bool
	holdingRegisters [16][1000]uint16
	inputRegisters   [16][1000]uint16
}

func main() {
	var server *modbus.ModbusServer
	var err error
	var eh *exampleHandler
	var ticker *time.Ticker

	eh = &exampleHandler{}

	server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL:        "tcp://localhost:5502",
		Timeout:    30 * time.Second,
		MaxClients: 5,
	}, eh)
	if err != nil {
		fmt.Printf("failed to create server: %v\n", err)
		os.Exit(1)
	}

	err = server.Start()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}

	ticker = time.NewTicker(1 * time.Second)
	for {
		<-ticker.C

		// eh.lock.Lock()

		// // 電子接点の状態をランダムに更新
		// // di := generateRandomBools(2000)
		// // copy(eh.discreteInputs[:], di[:])

		// // 入力レジスタの状態を更新 AI
		// // ir := generateRandomUint16s(2000)
		// // copy(eh.inputRegisters[:], ir[:])

		// eh.lock.Unlock()
	}
}

// テスト目的でランダムなブール値を生成する
// func generateRandomBools(length int) []bool {
// 	bools := make([]bool, length)
// 	for i := range bools {
// 		bools[i] = rand.IntN(2) == 1
// 	}
// 	return bools
// }

// // テスト目的でランダムなuint16値の配列を生成する
// func generateRandomUint16s(length int) []uint16 {
// 	uint16s := make([]uint16, length)
// 	for i := range uint16s {
// 		uint16s[i] = uint16(rand.IntN(65536))
// 	}
// 	return uint16s
// }

// DO
// (./modbus-cli --target tcp://localhost:5502 --unit-id 1 readCoils:<addr> ./modbus-cli --target tcp://localhost:5502 --unit-id 1 writeCoil:n:<true|false>)
func (eh *exampleHandler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	// ユニットIDが許可された範囲内であることを確認する
	if req.UnitId > 16 || req.UnitId == 0 {
		err = modbus.ErrIllegalFunction
		return nil, err
	}

	// このリクエストでカバーされるすべてのレジスタが実際に存在することを確認する
	if int(req.Addr)+int(req.Quantity) > len(eh.coils[req.UnitId-1])+len(eh.discreteInputs[req.UnitId-1]) {
		err = modbus.ErrIllegalDataAddress
		return
	}

	eh.lock.Lock()
	defer eh.lock.Unlock()

	// `req.Quantity` レジスタを、アドレス `req.Addr` から
	// `req.Addr + req.Quantity - 1` までループします。ここでは便宜上 `req.Addr + i` です。
	for i := 0; i < int(req.Quantity); i++ {
		if req.IsWrite {
			if req.Addr >= uint16(len(eh.coils[req.UnitId-1])) {
				// 値がコイル範囲を超える場合、DiscreteInput値を割り当てます。
				eh.discreteInputs[req.UnitId-1][int(req.Addr)+i-len(eh.coils[req.UnitId-1])] = req.Args[i]
			} else {
				// 値を割り当てる
				eh.coils[req.UnitId-1][int(req.Addr)+i] = req.Args[i]
			}
		} else {
			// 値を取得する
			res = append(res, eh.coils[req.UnitId-1][int(req.Addr)+i])
		}
	}
	return res, nil
}

// DI
// (./modbus-cli --target tcp://localhost:5502 --unit-id 1 readDiscreteInputs:<addr>)
func (eh *exampleHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	// ユニットIDが許可された範囲内であることを確認する
	if req.UnitId > 16 || req.UnitId == 0 {
		err = modbus.ErrIllegalFunction
		return nil, err
	}

	// このリクエストでカバーされるすべてのレジスタが実際に存在することを確認する
	if int(req.Addr)+int(req.Quantity) > len(eh.coils[req.UnitId-1]) {
		err = modbus.ErrIllegalDataAddress
		return
	}

	eh.lock.Lock()
	defer eh.lock.Unlock()

	// レジスタの値を取得します。
	res = make([]bool, req.Quantity)
	for i := range res {
		res[i] = eh.discreteInputs[req.UnitId-1][req.Addr+uint16(i)]
	}

	return res, nil
}

// AO
// (./modbus-cli --target tcp://localhost:5502 --unit-id 1 readHoldingRegisters:<type>:<addr>)
// (./modbus-cli --target tcp://localhost:5502 --unit-id 1 writeRegister:<type>:<addr>:<value>)
// <type> としてデコードされます。これは次のいずれかになります:
// - uint16: 符号なし 16 ビット整数、
// - int16: 符号付き 16 ビット整数、
// - uint32: 符号なし 32 ビット整数 (2 つの連続する modbus レジスタ)、
// - int32: 符号付き 32 ビット整数 (2 つの連続する modbus レジスタ)、
// - float32: 32 ビット浮動小数点数 (2 つの連続する modbus レジスタ)、
// - uint64: 符号なし 64 ビット整数 (4 つの連続する Modbus レジスタ)、
// - int64: 符号付き 64 ビット整数 (4 つの連続する Modbus レジスタ)、
// - float64: 64 ビット浮動小数点数 (4 つの連続する Modbus レジスタ)、
// - bytes: バイト文字列 (Modbus レジスタごとに 2 バイト)。
func (eh *exampleHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	var regAddr uint16

	// ユニットIDが許可された範囲内であることを確認する
	if req.UnitId > 16 || req.UnitId == 0 {
		err = modbus.ErrIllegalFunction
		return nil, err
	}

	// このリクエストでカバーされるすべてのレジスタが実際に存在することを確認する
	if int(req.Addr)+int(req.Quantity) > len(eh.holdingRegisters[req.UnitId-1])+len(eh.inputRegisters[req.UnitId-1]) {
		err = modbus.ErrIllegalDataAddress
		return
	}

	eh.lock.Lock()
	defer eh.lock.Unlock()

	// レジスタの読み取りまたは書き込みを処理します。
	for i := 0; i < int(req.Quantity); i++ {
		regAddr = req.Addr + uint16(i) // レジスタのアドレスを計算します。
		if req.IsWrite {
			if req.Addr >= uint16(len(eh.holdingRegisters[req.UnitId-1])) {
				// 値が保持レジスタの範囲外にある場合、入力レジスタに書き込みます。
				eh.inputRegisters[req.UnitId-1][int(regAddr)-len(eh.holdingRegisters[req.UnitId-1])] = req.Args[i]
			} else {
				// 値が保持レジスタの範囲内にある場合、保持レジスタに書き込みます。
				eh.holdingRegisters[req.UnitId-1][regAddr] = req.Args[i]
			}
		} else {
			res = append(res, uint16(eh.holdingRegisters[req.UnitId-1][int(req.Addr)+i]))
		}
	}
	return res, nil
}

// AI
// (./modbus-cli --target tcp://localhost:5502 --unit-id 1 readInputRegisters:<type>:<addr>)
// <type> としてデコードされます。これは次のいずれかになります:
// - uint16: 符号なし 16 ビット整数、
// - int16: 符号付き 16 ビット整数、
// - uint32: 符号なし 32 ビット整数 (2 つの連続する modbus レジスタ)、
// - int32: 符号付き 32 ビット整数 (2 つの連続する modbus レジスタ)、
// - float32: 32 ビット浮動小数点数 (2 つの連続する modbus レジスタ)、
// - uint64: 符号なし 64 ビット整数 (4 つの連続する Modbus レジスタ)、
// - int64: 符号付き 64 ビット整数 (4 つの連続する Modbus レジスタ)、
// - float64: 64 ビット浮動小数点数 (4 つの連続する Modbus レジスタ)、
// - bytes: バイト文字列 (Modbus レジスタごとに 2 バイト)。
func (eh *exampleHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	// ユニットIDが許可された範囲内であることを確認する
	if req.UnitId > 16 || req.UnitId == 0 {
		err = modbus.ErrIllegalFunction
		return nil, err
	}

	// このリクエストでカバーされるすべてのレジスタが実際に存在することを確認する
	if int(req.Addr)+int(req.Quantity) > len(eh.inputRegisters[req.UnitId-1]) {
		err = modbus.ErrIllegalDataAddress
		return nil, err
	}

	eh.lock.RLock()
	defer eh.lock.RUnlock()

	// req.addr から req.addr + req.Quantity - 1 までのすべてのレジスタ アドレスをループします。
	for regAddr := req.Addr; regAddr < req.Addr+req.Quantity; regAddr++ {
		// レジスタの値を取得します。
		res = append(res, uint16(eh.inputRegisters[req.UnitId-1][regAddr]))
	}
	return res, nil
}
