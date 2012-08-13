package main
import "fmt"
import "math/rand"
import "time"

//import "bytes"
//import "string"

const CROSSOVERRATE = 0.7
const MUTATIONRATE = 0.001
const POPULATION = 10
const MAXGENERATIONS = 1000
const NUMOPERATORS = 3
const NUMBITS = 4*(NUMOPERATORS*2+1) //28 in the case of 3 operators
const SYMBOLLENGTH = 4 //in bits

func Log(v ...interface{}){
	enableDebug := true
	if enableDebug{ 
		fmt.Println(v...)
	}
}


const (
	broken = iota
	number
	operator
)

const (
	add = iota + 10
	sub
	mult
	div
	none
)

func argumentType(arg string)(ret int){

	switch arg {

	case "0000","0001","0010","0011",
		"0100","0101","0110","0111",
		"1000","1001": return number
	case "1010","1011","1100","1101":return operator
	default: return broken
	}

	return broken
}

func doMath(operator int, tempinput int, storage float32)(res float32){
	input := float32(tempinput)

	switch operator{
	case add:return storage+input
	case sub:return storage-input
	case mult:return storage*input
	case div:return storage/input
	default: return 99999.9 //Break
	}
	return 0.1
} 

func parseNumeric(arg string)(ret int){

	switch arg {
	case "0000":return 0
	case "0001":return 1
	case "0010":return 2
	case "0011":return 3
	case "0100":return 4
	case "0101":return 5
	case "0110":return 6
	case "0111":return 7
	case "1000":return 8
	case "1001":return 9
	}
	return -1 //TODO handle error
}


func parseOperator(arg string)(ret int){
	switch arg{
	case "1010":return add
	case "1011":return sub
	case "1100":return mult
	case "1101":return div
	}
	return none
}

func humanReadOperator(arg int)(ret string){
	switch arg{
	case add:return "+"
	case sub:return "-"
	case mult:return "*"
	case div:return "/"
	}
	return "broken"
}


func generateOneChrom()(string){
	var temp string
	for i:=0;i<NUMBITS;i++{
		tilf := rand.Float32()
		if tilf < 0.5{
			temp=temp+"0"
		}else{
			temp=temp+"1"
		}
	}
	return temp
}

func generateNChroms(n int)([]string){
	ret := make([]string,n)
	rand.Seed(time.Now().Unix())
	for i:=0;i<n;i++{
		ret[i] = generateOneChrom()
	}
	return ret
}


//Parse the string 4 by 4 bit and calculate the expression 
func calcFitness(chromStr string)(ret float32){
	var currentval float32 =  0.0
	currentOperator := add
	next := number
	tempOperator := none
	tempNumeric := 0
	
	for (len(chromStr) > 0) {
		thisString := chromStr[0:SYMBOLLENGTH]
		if next == number{
			if argumentType(thisString) == number{
				//calculate new currentval based on  "prev operatoe" and "currentval
				tempNumeric = parseNumeric(thisString) 
				currentval = doMath(currentOperator,tempNumeric,currentval)
				Log(tempNumeric)
				next = operator
			}
			
		}else if next == operator{
			//if arg type is operaotr, store it for use in possibly next numerical
			if argumentType(thisString) == operator{
				tempOperator = parseOperator(thisString)
				currentOperator = tempOperator
				Log(humanReadOperator(tempOperator))
				next = number
			}

		}
		chromStr = chromStr[4:]
	}
	Log("=")
	return currentval
}



//An attempt to run a genetic algorithm...
func main(){
	fmt.Print("Hello genetic algorithms!\n")
	lol := generateNChroms(1)
	fmt.Println(lol)
	fmt.Println(calcFitness(lol[0]))
	//population := generateNChroms(POPULATION)
	//fmt.Print(population)


	
	
	


	//k := 8
	//var test string
//	test = "1010101010101010101010101010\n"
//	test = test[0:4]+"0000"+test[8:28]
//	fmt.Print(test[:])

	//test := generateNChroms(20)
	//fmt.Print(test)
	
}


	/* TEST BLOCK
	fmt.Println(argumentType("0010")) // 1
	fmt.Println(argumentType("1011")) // 2
	fmt.Println(argumentType("hei")) // 0

	fmt.Println(parseNumeric("0011")) // 3
	fmt.Println(parseNumeric("1111")) // broken -1

	fmt.Println(parseOperator("1010")) // add 0
	fmt.Println(parseOperator("lol")) // none 4

	
	fmt.Println(calcFitness("lol")) // 1.0
	*/