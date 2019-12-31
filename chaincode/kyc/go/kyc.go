/*
 * The sample smart contract for documentation topic:
 * Writing  Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	//	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the Customer structure, with 4 properties.  Structure tags are used by encoding/json library
type Customer struct {
	Name          string `json:"name"`
	SSN           string `json:"ssn"`
	DateOfBirth   string `json:"dateofbirth"`
	Gender        string `json:"gender"`
	BankName      string `json:"bankaccount"`   //it will come from signing identity like bank msp
	AccountNumber string `json:"accountnumber"` //autogenerated for new customer
}

/*
 * The Init method is called when the Smart Contract "kyc" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "kyc"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCustomer" {
		return s.queryCustomer(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createCustomer" {
		return s.createCustomer(APIstub, args)
	} else if function == "queryAllCustomers" {
		return s.queryAllCustomers(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCustomer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	customerAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(customerAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	 customers := []Customer{
	 	Customer{Name: "Hitesh", SSN: "ssnhit123", DateOfBirth: "15-12-1990", Gender: "male"},
	 	Customer{Name: "Rachit", SSN: "ssnrac456", DateOfBirth: "12-09-1991", Gender: "male"},
	 	Customer{Name: "Ankit", SSN: "ssnank789", DateOfBirth: "10-08-1993", Gender: "male"},
	 	Customer{Name: "john", SSN: "ssnjohn123", DateOfBirth: "123-11-1994", Gender: "male"},
	 }

	 i := 0
	 for i < len(customers) {
	 	fmt.Println("i is ", i)
	 	customerAsBytes, _ := json.Marshal(customers[i])
	 	APIstub.PutState("Customer"+strconv.Itoa(i), customerAsBytes)
	 	fmt.Println("Added", customers[i])
	 	i = i + 1
	 }

	return shim.Success(nil)
}

func (s *SmartContract) createCustomer(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	///generate unique account number

	//fmt.Println(EncodeToString(10))

	var accnumber = EncodeToString(10)
	args[6] = accnumber
	var customer = Customer{Name: args[1], SSN: args[2], DateOfBirth: args[3], Gender: args[4], BankName: args[5], AccountNumber: args[6]}

	customerAsBytes, _ := json.Marshal(customer)
	APIstub.PutState(args[0], customerAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllCustomers(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Customer0"
	endKey := "Customer999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCustomers:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//function to generate a unique account number
func EncodeToString(max int) string {

	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)

}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

