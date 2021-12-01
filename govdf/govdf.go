package govdf

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
)

var server = ""

type discrResponse struct {
	Discriminant string `json:"discriminant"`
}

type proveResponse struct {
	Y     big.Int `json:"output"`
	Proof big.Int `json:"proof"`
}

type verifyResponse struct {
	IsValid bool `json:"valid"`
}

func GetRandomSeed() []byte {
	seed := make([]byte, 16)
	rand.Read(seed)
	return seed
}

func SetServer(newServer string) {
	server = newServer
}

/*
	EVAL function
	receives:
		x: input of VDF
		T: number of iterations (squarings)
		ds: discriminant size
		seed: set randomness on discriminant creation
	returns:
		(result, proof)
*/
func Eval(x big.Int, T, ds int, seed []byte) (big.Int, big.Int) {

	postBody, _ := json.Marshal(map[string]string{
		"seed":              hex.EncodeToString(seed),
		"input":             x.String(),
		"iterations":        strconv.Itoa(T),
		"discriminant_size": strconv.Itoa(ds),
	})
	responseBody := bytes.NewBuffer(postBody)

	// Leverage Go's HTTP Post function to make request
	resp, err := http.Post(server+"/eval", "application/json", responseBody)

	// Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	decoder := json.NewDecoder(resp.Body)
	var s proveResponse
	err = decoder.Decode(&s)

	if err != nil {
		panic(err)
	}

	// Should return two Integers
	return s.Y, s.Proof
}

/*
	VERIFY function
	receives:
		x: input of VDF
		y: result of VDF
		pi: the proof of VDF result
		T: number of iterations (squarings)
		ds: discriminant size
		seed: set randomness on discriminant creation
	returns if verification was correct
*/
func Verify(x, y, pi big.Int, T, ds int, seed []byte) bool {

	postBody, _ := json.Marshal(map[string]string{
		"seed":              hex.EncodeToString(seed),
		"discriminant_size": strconv.Itoa(ds),
		"input":             x.String(),
		"output":            y.String(),
		"proof":             pi.String(),
		"iterations":        strconv.Itoa(T),
	})
	responseBody := bytes.NewBuffer(postBody)

	// Leverage Go's HTTP Post function to make request
	resp, err := http.Post(server+"/verify", "application/json", responseBody)

	// Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	decoder := json.NewDecoder(resp.Body)
	var s verifyResponse
	err = decoder.Decode(&s)

	if err != nil {
		panic(err)
	}

	return s.IsValid
}

// Creates a discriminant for binary quadratic forms
// We no longer need this function :)
func _createDiscriminant(discriminantSize int, seed []byte) string {

	postBody, _ := json.Marshal(map[string]string{
		"discriminant_size": strconv.Itoa(discriminantSize),
		"seed":              hex.EncodeToString(seed),
	})
	responseBody := bytes.NewBuffer(postBody)

	// Leverage Go's HTTP Post function to make request
	resp, err := http.Post(server+"/create", "application/json", responseBody)

	// Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	decoder := json.NewDecoder(resp.Body)
	var s discrResponse
	err = decoder.Decode(&s)

	if err != nil {
		panic(err)
	}

	return s.Discriminant
}
