package main
import "fmt"
import "math/rand"
import "time"
import "sort"
//import "reflect"
//import "strings"
//import "bytes"


const CROSSOVERRATE = 0.7
const MUTATIONRATE = 0.001
const POPULATION = 15
const MAXGENERATIONS = 10
const NUMOPERATORS = 3
const NUMBITS = 4*(NUMOPERATORS*2+1) //28 in the case of 3 operators
const SYMBOLLENGTH = 4 //in bits

func Deb(v ...interface{}){
	enableDebug := false
	if enableDebug{ 
		fmt.Println(v...)
	}
}


func Log(v ...interface{}){
	enableLog := true
	if enableLog{ 
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

func doMath(operator int, tempinput int, storage float64)(res float64){
	input := float64(tempinput)

	switch operator{
	case add:return storage+input
	case sub:return storage-input
	case mult:return storage*input
	case div:return storage/input
	default: return 99999.9 //Break
	}
	return 0.1 //TODO handle error????
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
		tilf := rand.Float64()
		if tilf < 0.5{
			temp+="0"
		}else{
			temp+="1"
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
func evalExpression(chromStr string)(ret float64){
	var currentval float64 =  0.0
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
				//handle divide by zero
				if tempNumeric == 0 && currentOperator == div{
					Log("div by zero")
					return -1 
					
				}
				currentval = doMath(currentOperator,tempNumeric,currentval)
				Deb(tempNumeric)
				next = operator
			}
			
		}else if next == operator{
			//if arg type is operaotr, store it for use in possibly next numerical
			if argumentType(thisString) == operator{
				tempOperator = parseOperator(thisString)
				currentOperator = tempOperator
				Deb(humanReadOperator(tempOperator))
				next = number
			}

		}
		chromStr = chromStr[4:]
	}
	Deb("=")
	return currentval
}

//The function to optimize
func calcFitness(inputVal, goal float64)(fitness float64,correct bool){
	correct = false
	if inputVal == goal{
		correct = true
	}
	fitness = 1/abs((goal-inputVal))
	return
}

func abs(in float64)(ret float64){
	if in < 0{
		return -in
	}
	return in
}


//TODO finalize
func mateOneGeneration(popIn []string,goal float64)(popOut []string){
	popOut = make([]string,len(popIn))
	fitness := make([]float64,len(popIn))
	for i,chromIn := range popIn{
		fitness[i],_ = calcFitness(evalExpression(chromIn),goal)
	}
//	Log(prepareRoulette(fitness))	
	
	for i:=0;i<len(popIn)/2;i++{
		firstMateIndex := pickWinner(prepareRoulette(fitness))
		secondMateIndex := pickWinner(prepareRoulette(fitness))
		for secondMateIndex == firstMateIndex{
			secondMateIndex = pickWinner(prepareRoulette(fitness))
		}
		popOut[i],popOut[i+1] = crossOver(popIn[firstMateIndex],popIn[secondMateIndex])
	}
	return
}

func crossOver(chromOne string, chromTwo string)(string, string){
	//next line is test
	//chromTest := string(chromOne[0])
	//Log(reflect.TypeOf(chromTest))
	//Log(chromTest)
	crossOverCheck := rand.Float64()
	if crossOverCheck < CROSSOVERRATE{
		chosenGene := int((rand.Float64())*float64(len(chromOne)))
		temp := chromOne
		chromOne = chromOne[0:chosenGene] + chromTwo[chosenGene:]
		chromTwo = chromTwo[0:chosenGene] + temp[chosenGene:]
	}
	chromOne = mutateString(chromOne)
	chromTwo = mutateString(chromTwo)
	return chromOne,chromTwo
}

//TODO verify this func
func mutateString(chrom string)(ret string){
	ret = chrom
	for i,_ := range chrom{
		mutateCheck := rand.Float64()
		if mutateCheck < MUTATIONRATE{
			charTest := string(chrom[i])
			if charTest == "0"{
				ret = ret[0:i] + "1" + ret[i+1:]
			}else{
				ret = ret[0:i] + "0" + ret[i+1:]
			}
		}
	}
	return ret
}

//Picks a random chromosome based on fitness
func prepareRoulette(fitnessTable []float64)([]float64){
	ret := make([]float64,len(fitnessTable))
	ret[0] = fitnessTable[0]
	for i,_ := range fitnessTable{
		if i != 0{
			ret[i] = fitnessTable[i]+ret[i-1]
		}
	}
	//Log(winner)
	//Log(winnerIndex)
	return ret
}

func pickWinner(rouletteWheel []float64)(ret int){
	largest := rouletteWheel[len(rouletteWheel)-1]
	winner := rand.Float64()*largest
	winnerIndex := sort.SearchFloat64s(rouletteWheel,winner)
	return winnerIndex
}



//An attempt to run a genetic algorithm...
func main(){
	var target float64 = 42
	fmt.Println("Hello genetic algorithms!")
	lol := generateNChroms(POPULATION)
	best := float64(0.0)
	bestFit := float64(0.0)
	curr := float64(0.0)
	currFit := float64(0.0)
	goalReached := false

	//Log(mutateString(lol[0]))
	mateOneGeneration(lol,target)

	
	//lol[0],lol[1] = crossOver(lol[0],lol[1])
	for i:=0;i<POPULATION;i++{
		curr = evalExpression(lol[i])
		currFit,goalReached = calcFitness(curr,target)
		if !goalReached {
			if currFit > bestFit{
				best = curr
				bestFit = currFit
			}
		}else{
			bestFit = currFit
			best = curr
			Log("Evolution perfected!",target)
			
			break;
		}

	}

	if !goalReached{
		Log("Population failed")
		Log("Best loser is:")
		Log(best)
		Log(bestFit)
	}
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