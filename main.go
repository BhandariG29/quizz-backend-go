package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main(){
	//1. input the name of the file
	fNmae := flag.String("f", "quizz.csv", "path of csv file")

	//2. set the duration of timer
	timer := flag.Int("t", 30, "timer for the quiz")
	flag.Parse()

	//3. pull the problems from the file (caliing problem puller func)
	problems, err := problemPuller(*fNmae)

	//4. handle the error
	if err != nil{
		exit(fmt.Sprintf("something went wrong: %s", err.Error()))
	}

	//5. create a variable to count our answers
	correctAns := 0

	//6. initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

	//7. loop through problems, print the questions, and check the answer
	problemLoop:
	for i, p := range problems{
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.q)

		go func(){
			fmt.Scanf("%s", &answer)
			ansC <- answer
		}()

		select{
		case <-tObj.C:
			fmt.Println()
			break problemLoop

		case iAns := <-ansC:
			if iAns == p.a{
				correctAns++
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}

	//8. calculate and print the result
	fmt.Printf("your result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("press enter to exit")
	<- ansC
}

func problemPuller(fileName string) ([]problem, error){
	//read all the problems from the quizz.csv

	//1. open the file
	if fObj, err := os.Open(fileName); err == nil{
		//2. create a new reader
		csvR := csv.NewReader(fObj)

		//3. read all the problems
		if clines, err := csvR.ReadAll(); err == nil{
			//4. call problemParse func
			return parseProblem(clines), nil 
		}else{
			return nil, fmt.Errorf("error in reading data in csv" + "format from %s file: %s", fileName, err.Error())
		}
	}else{
		return nil, fmt.Errorf("error in opening %s file: %s", fileName, err.Error())
	}
	
}

func parseProblem(lines [][]string) []problem{
	r := make([]problem, len(lines))

	for i:=0; i<len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}

	return r
}

type problem struct{
	q string
	a string
}

func exit(msg string){
	fmt.Println(msg)
	os.Exit(1)
}