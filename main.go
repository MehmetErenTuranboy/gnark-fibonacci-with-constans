package main

import (
	"fmt"
	"log"

	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type fibCircuit struct {
	A, B frontend.Variable `gnark:",secret"` // a and b are secret variables
	Result frontend.Variable `gnark:",public"` // result is public
}

func (circuit *fibCircuit) Define(api frontend.API) error {

	fmt.Println("Result of circuit will be %d", circuit.Result)

	a := circuit.A
	b := circuit.B

	for i := 2; i <= 100; i++ {
		next := api.Add(a, b)
		a, b = b, next
	}
	api.AssertIsEqual( circuit.Result, b)
	return nil
}


func main() {
	var circuit fibCircuit

	firstSecret := *big.NewInt(0)
	secondSecret := *big.NewInt(1)
	result := "354224848179261915075"

	// Compile the circuit into a set of constraints
	ccs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		log.Fatalf("Failed to compile the circuit: %v", err)
	}

	// Setup the Proving and Verifying keys
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		log.Fatalf("Failed to setup the proving and verifying keys: %v", err)
	}

	resultBigInt, success := new(big.Int).SetString(result, 10)
	if !success {
		log.Fatalf("Failed to convert the string to big.Int")
	}

	assignment := fibCircuit{
		A:      firstSecret,
		B:      secondSecret,
		Result: *resultBigInt,
	}

	// Create a witness from the assignment
	witness, err := frontend.NewWitness(&assignment, ecc.BN254)
	if err != nil {
		log.Fatalf("Failed to create a witness: %v", err)
	}

	// Extract the public part of the witness
	publicWitness, err := witness.Public()
	if err != nil {
		log.Fatalf("Failed to extract the public witness: %v", err)
	}

	// Prove the witness
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		log.Fatalf("Failed to prove the witness: %v", err)
	}

	//fmt.Println("Generated Proof:", proof)

	// Verify the proof
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Println("Verification Result: Failed")
		log.Fatalf("Failed to verify the proof: %v", err)
	} else {
		fmt.Println("Verification Result: Success")
	}
}