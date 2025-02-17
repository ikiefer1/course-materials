// bhg-scanner/scanner.go modified from Black Hat Go > CH2 > tcp-scanner-final > main.go
// Code : https://github.com/blackhat-go/bhg/blob/c27347f6f9019c8911547d6fc912aa1171e6c362/ch-2/tcp-scanner-final/main.go
// License: {$RepoRoot}/materials/BHG-LICENSE
//go build scanner.go
// Useage: Used when running the main.go program in the main folder. 
//Go back to course-materials/materials/lab/2/bhg-scanner/main and "go build main.go"
//PortScanner() is called with the following: ./main
//Isaiah Kiefer, Lab 2

package scanner

import (
	"fmt"
	"net"
	"sort"
	"time"
)

//TODO 3 : ADD closed ports; currently code only tracks open ports
var closedports []int
var openports []int  // notice the capitalization here. access limited!


func worker(ports, results chan int){//}, chans []chan int) {
	
	for p := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", p)    
		conn, err := net.DialTimeout("tcp", address, 1 * time.Second) // TODO 2 : REPLACE THIS WITH DialTimeout (before testing!)
		if err != nil { 
			results <- 0
			//chans[0]<-0
			//chans[1]<-p
			continue
		}
		conn.Close()
		results <- p
	}
}

// for Part 5 - consider
// easy: taking in a variable for the ports to scan (int? slice? ); a target address (string?)?
// med: easy + return  complex data structure(s?) (maps or slices) containing the ports.
// hard: restructuring code - consider modification to class/object 
// No matter what you do, modify scanner_test.go to align; note the single test currently fails
func PortScanner(start int, end int) ([]int,[]int) {  
	openports = []int{}
	closedports = []int{}
	//make sure the start is less than the end
	if end - start < 0{
		return openports, closedports
	}
	ports := make(chan int, 100)   // TODO 4: TUNE THIS FOR CODEANYWHERE / LOCAL MACHINE
	//chans := make([]chan int, 2)	
	results := make(chan int)
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}
	
	go func() {
		//for i := 1; i <= 1024; i++ {
		for i:= start; i<=end; i++{
			ports <- i
		}
	}()
	
	// for i := 0; i < 1024; i++ {
	for i := 0; i < (end - start)+1; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	//trackClosedPorts tracks the closed ports
	trackClosedPorts := start

	//this math determines how many closed ports there should be
	numberOfClosedPorts := ( (end - start) + 1 ) - len(openports)

	//cycle through and get the number for each closed port while skipping open ports
	for i:=0; i < numberOfClosedPorts; i++{
		for _, val:= range openports{

			if val==trackClosedPorts{
				trackClosedPorts++
			}
		}
		closedports = append(closedports, trackClosedPorts)
		trackClosedPorts++
	}
	//sort closedports
	sort.Ints(closedports)

	//TODO 5 : Enhance the output for easier consumption, include closed ports
	fmt.Printf("\n************Open Ports************\n")
	for _, port := range openports {
		fmt.Printf("%d open, ", port)
	}
	fmt.Printf("\n************Closed Ports************\n")
	for _, port := range closedports {

		fmt.Printf("%d closed, ", port)
	}
	fmt.Printf("\n")
	return openports, closedports
	//return len(openports) // TODO 6 : Return total number of ports scanned (number open, number closed); 
	//you'll have to modify the function parameter list in the defintion and the values in the scanner_test
}
