package main

import "fmt"
import "strconv"
import "encoding/json"
import "io/ioutil"
import "log"
import "strings"
import "math"

type Memory struct{
	Memories []MemoryOperation
}

type MemoryOperation struct {
    Result float64
    Operations []string
}

var operators = [...]string{"+","-","*","/"}
var memoryPath = "memory.json"

func main() {
	var item1 float64
	fmt.Printf("Premier item:")
    fmt.Scanf("%f\n", &item1)
		var item2 float64
	fmt.Printf("\nSecond item:")
    fmt.Scanf("%f\n", &item2)
	var goal float64
	fmt.Printf("\nObjectif:")
    fmt.Scanf("%f\n", &goal)
	var result string = searchSolution(item1,item2,goal)
    fmt.Printf(result + "=" + floatToString(goal))
	//fmt.Printf("\n%d", countPossibility(4))
}

func calcul (item1 float64, item2 float64, oneOperator string) (result float64){
	switch oneOperator {
	case "+":
	result = item1 + item2
	case "-":
	result = item1 - item2
	case "/":
	if item2 == 0 {
		result = 0
	}else{
		result = item1 / item2
	}
	case "*":
	result = item1 * item2
	}
	return 
}

func searchSolution(item1 float64, item2 float64, goal float64) (operation string){
	op, success:= searchMemorySolution(item1, item2, goal)
	if success == true {
		operation = op 
		return
	}
	for _, oneOperator := range operators {
		for i := 0; i< 2; i++{
		    var result float64 = calcul(item1, item2, oneOperator)
			if result == goal{
				operation = floatToString(item1) + oneOperator + floatToString(item2)
				// Mise en memoire
				success := learnOperation(operation, goal)
				if success == false{
					fmt.Printf("Erreur de memorisation")
				}
				return
			}
			changePosition(&item1, &item2)
		}
    }
	operation = "Erreur"
	return
}

func searchMemorySolution(item1 float64, item2 float64, goal float64)(operation string, success bool){
	success = false
	operation = ""
	var data Memory
	file, e := ioutil.ReadFile(memoryPath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}
	err := json.Unmarshal(file, &data)
	var memoriesOperation []MemoryOperation
	if err == nil{
		memoriesOperation = data.Memories
		for i := 0; i < len(memoriesOperation); i++{
			if float64(memoriesOperation[i].Result) == goal{
				for _,memoryOperation := range memoriesOperation[i].Operations{
					var operatorIndex int
					for _, oneOperator := range operators {
						operatorIndex = strings.Index(memoryOperation, oneOperator)
						if operatorIndex != -1 {
							break
						}
					}
					split := strings.Split(string(memoryOperation), string(memoryOperation[operatorIndex]))
					if floatToString(item1) == split[0] && floatToString(item2) == split[1]{
						operation = memoryOperation
						success = true
						return
					}
				}
			}
		 }
	}
	return 
}

func learnOperation(operation string, result float64)(success bool){
	success = true
	file, e := ioutil.ReadFile(memoryPath)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}
	var data Memory
	err := json.Unmarshal(file, &data)
	var memoriesOperation []MemoryOperation
	if err == nil{
		memoriesOperation = data.Memories
		var countNotGoodResult int = 0
		for i := 0; i < len(memoriesOperation); i++{
			if float64(memoriesOperation[i].Result) == result{
				var countNotGoodOperation int = 0
				for _,memoryOperation := range memoriesOperation[i].Operations{
					if memoryOperation != operation {
						countNotGoodOperation++
					}
				}
				if countNotGoodOperation == len(memoriesOperation[i].Operations) {
					memoriesOperation[i].Operations = append(memoriesOperation[i].Operations, operation)
				}
			}else{
				countNotGoodResult++
			}
		}
		if countNotGoodResult == len(memoriesOperation) {
			fmt.Printf("%d", countNotGoodResult)
			var operations []string
			operations = append(operations, operation)
			memoriesOperation = append(memoriesOperation,MemoryOperation{result, operations}) 
		}
	}else{
		log.Println(err)
		var operations []string
		operations = append(operations, operation)
		memoriesOperation = append(memoriesOperation,MemoryOperation{result, operations})
	}
	mo := Memory{memoriesOperation}
	b, errM := json.Marshal(mo)
	errW := ioutil.WriteFile(memoryPath, b, 0644)
	if errM != nil { 
		log.Println(errM) 
		success = false
	}
	if errW != nil { 
		log.Println(errW) 
		success = false
	}
	return
}

func countPossibility(numberItem int)(result int){
	var factorial int = 1
	for i:=1; i<=numberItem; i++{
		factorial *= i
	}
	result = int(math.Pow(float64(len(operators)), float64(numberItem-1))* float64(factorial))
	return
}

func changePosition(item1 * float64,item2 * float64){
	var temp float64 = *item1
	*item1 = *item2
	*item2 = temp
}

func floatToString(value float64)(stringConvert string){
	stringConvert = strconv.FormatFloat(value, 'f', -1, 64)
	return
}