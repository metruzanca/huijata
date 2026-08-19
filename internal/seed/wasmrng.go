package seed

import (
	_ "embed"
	"fmt"
	"sync"

	"github.com/bytecodealliance/wasmtime-go/v9"
)

//go:embed rng.wasm
var rngWasm []byte

var (
	engine *wasmtime.Engine
	store  *wasmtime.Store
	initOnce sync.Once
	initErr  error

	mu sync.Mutex

	setWorldSeedFunc      *wasmtime.Func
	setRandomSeedFunc     *wasmtime.Func
	randomFunc            *wasmtime.Func
	randomIntFunc         *wasmtime.Func
	randomRoundedFunc     *wasmtime.Func
	randomDistributionFunc *wasmtime.Func
)

func initWasm() error {
	initOnce.Do(func() {
		engine = wasmtime.NewEngine()
		store = wasmtime.NewStore(engine)

		mod, err := wasmtime.NewModule(engine, rngWasm)
		if err != nil {
			initErr = fmt.Errorf("failed to compile wasm module: %w", err)
			return
		}

		instance, err := wasmtime.NewInstance(store, mod, []wasmtime.AsExtern{})
		if err != nil {
			initErr = fmt.Errorf("failed to instantiate wasm: %w", err)
			return
		}

		setWorldSeedFunc = instance.GetExport(store, "SetWorldSeed").Func()
		setRandomSeedFunc = instance.GetExport(store, "SetRandomSeed").Func()
		randomFunc = instance.GetExport(store, "Random").Func()
		randomIntFunc = instance.GetExport(store, "RandomInt").Func()
		randomRoundedFunc = instance.GetExport(store, "RandomRounded").Func()
		randomDistributionFunc = instance.GetExport(store, "RandomDistribution").Func()
	})
	return initErr
}

func wasmSetWorldSeed(seed uint32) {
	mu.Lock()
	defer mu.Unlock()
	if err := initWasm(); err != nil {
		return
	}
	setWorldSeedFunc.Call(store, int32(seed))
}

func wasmSetRandomSeed(x, y float64) {
	mu.Lock()
	defer mu.Unlock()
	if err := initWasm(); err != nil {
		return
	}
	setRandomSeedFunc.Call(store, x, y)
}

func wasmRandom() (float32, error) {
	mu.Lock()
	defer mu.Unlock()
	if err := initWasm(); err != nil {
		return 0, err
	}
	result, err := randomFunc.Call(store)
	if err != nil {
		return 0, err
	}
	return result.(float32), nil
}

func wasmRandomInt(min, max int32) (int32, error) {
	mu.Lock()
	defer mu.Unlock()
	if err := initWasm(); err != nil {
		return 0, err
	}
	result, err := randomIntFunc.Call(store, min, max)
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func wasmRandomRounded(min, max float64) (int32, error) {
	mu.Lock()
	defer mu.Unlock()
	if err := initWasm(); err != nil {
		return 0, err
	}
	result, err := randomRoundedFunc.Call(store, min, max)
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

func wasmRandomDistribution(min, max, mean, sharpness int32) (int32, error) {
	mu.Lock()
	defer mu.Unlock()
	if err := initWasm(); err != nil {
		return 0, err
	}
	result, err := randomDistributionFunc.Call(store, min, max, mean, sharpness)
	if err != nil {
		return 0, err
	}
	return result.(int32), nil
}

type WasmRng struct {
	worldSeed uint32
}

func NewWasmRng(worldSeed uint32) *WasmRng {
	wasmSetWorldSeed(worldSeed)
	return &WasmRng{worldSeed: worldSeed}
}

func (r *WasmRng) SetWorldSeed(ws uint32) {
	r.worldSeed = ws
	wasmSetWorldSeed(ws)
}

func (r *WasmRng) SetRandomSeed(x, y float64) {
	wasmSetRandomSeed(x, y)
}

func (r *WasmRng) Next() float64 {
	v, _ := wasmRandom()
	return float64(v)
}

func (r *WasmRng) Random(min, max int32) int32 {
	v, _ := wasmRandomInt(min, max)
	return v
}

func (r *WasmRng) RandomRounded(min, max float64) int32 {
	v, _ := wasmRandomRounded(min, max)
	return v
}

func (r *WasmRng) RandomDistribution(min, max, mean, sharpness int32) int32 {
	v, _ := wasmRandomDistribution(min, max, mean, sharpness)
	return v
}
