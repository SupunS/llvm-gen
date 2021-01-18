package main

import (
    "fmt"
    "github.com/llir/llvm/ir"
    "github.com/llir/llvm/ir/constant"
    "github.com/llir/llvm/ir/enum"
    "github.com/llir/llvm/ir/types"
    "github.com/llir/llvm/ir/value"
    "os"
)


var printf *ir.Func
var cons *ir.Global
var zero = constant.NewInt(types.I32, 0)
var limitConst = constant.NewInt(types.I32, 100_000_000)

func main() {
    module := ir.NewModule()

    // Create a binding for native printf function to print results
    printf = module.NewFunc(
        "printf",
        types.I32,
        ir.NewParam("format", types.NewPointer(types.I8)),
    )
    printf.Sig.Variadic = true

    // Define global vars
    cons = module.NewGlobalDef("fmtString", constant.NewCharArrayFromString("Value: %i \n"))
    count := module.NewGlobalDef("count", zero)
    limit := module.NewGlobalDef("limit", limitConst)

    // Create main function, and the  entry BB.
    // `int main(void) { ... }`
    main := module.NewFunc("main", types.I32)
    entry := main.NewBlock("entry")

    // Three basic blocks for condition, while-body, and while-done
    conditionBB := main.NewBlock("condition")
    body := main.NewBlock("body")
    done := main.NewBlock("done")
    entry.NewBr(conditionBB)

    // Condition
    tmpCount := conditionBB.NewLoad(types.I32, count)
    tmpLimit := conditionBB.NewLoad(types.I32, limit)
    condition := conditionBB.NewICmp(enum.IPredSLT, tmpCount, tmpLimit)
    conditionBB.NewCondBr(condition, body, done)

    // While body
    one := constant.NewInt(types.I32, 1)
    oldCount := body.NewLoad(types.I32, count)
    newCount := body.NewAdd( one, oldCount)
    body.NewStore(newCount, count)
    //printValue(body, count)
    body.NewBr(conditionBB)

    // While-done
    success := constant.NewInt(types.I32, 0)
    printValue(done, count)
    done.NewRet(success)

    // Print the LLVM IR assembly of the module.
    fmt.Println(module)
    writeToFile(module)
}

func writeToFile(mod *ir.Module)  {
    f, _ := os.Create("./llmv_ir.ll")
    _, err := f.WriteString(mod.String())
    if err != nil {
        // TODO:panic err
    }
}

func printValue(currentBlock *ir.Block, value value.Value)  {
    fmtPtr := currentBlock.NewGetElementPtr(cons.ContentType, cons, zero, zero)
    tmp0 := currentBlock.NewLoad(types.I32, value)
    currentBlock.NewCall(printf, fmtPtr, tmp0)
}
