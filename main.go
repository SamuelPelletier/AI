package main

import "fmt"
import "strconv"
import "encoding/json"
import "io/ioutil"
import "log"

type Memory struct{
	Memories []MemoryOperation
}

type MemoryOperation struct {
    Result float64
    Operations []string
}

var operator = [...]string{"+","-","*","/"}
var memoryPath = "memory.json"

func main() {
	var element1 float64
	fmt.Printf("Premier element:")
    fmt.Scanf("%f\n", &element1)
		var element2 float64
	fmt.Printf("\nSecond element:")
    fmt.Scanf("%f\n", &element2)
	var goal float64
	fmt.Printf("\nObjectif:")
    fmt.Scanf("%f\n", &goal)
	var result string = searchSolution(element1,element2,goal)
    fmt.Printf(result + "=" + floatToString(goal))
}

func calcul (element1 float64, element2 float64, oneOperator string) (result float64){
	switch oneOperator {
	case "+":
	result = element1 + element2
	case "-":
	result = element1 - element2
	case "/":
	if element2 == 0 {
		result = 0
	}else{
		result = element1 / element2
	}
	case "*":
	result = element1 * element2
	}
	return 
}

func searchSolution(element1 float64, element2 float64, goal float64) (operation string){
	op, success:= searchMemorySolution(element1, element2, goal)
	if success == true {
		operation = op 
		return
	}
	for _, oneOperator := range operator {
        var result float64 = calcul(element1, element2, oneOperator)
		if result == goal{
			operation = floatToString(element1) + oneOperator + floatToString(element2)
			// Mise en memoire
			success := learnOperation(operation, goal)
			if success == false{
				fmt.Printf("Erreur de memorisation")
			}
			return
		}
    }
	operation = "Erreur"
	return
}

func searchMemorySolution(element1 float64, element2 float64, goal float64)(operation string, success bool){
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
					if floatToString(element1) == string(memoryOperation[0]) && floatToString(element2) == string(memoryOperation[len(memoryOperation)-1]){
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

func floatToString(value float64)(stringConvert string){
	stringConvert = strconv.FormatFloat(value, 'f', -1, 64)
	return
}