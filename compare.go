////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Elixxir                                                    /
//                                                                             /
// Permission is hereby granted, free of charge, to any person obtaining a     /
// copy of this software and associated documentation files (the “Software”),  /
// to deal in the Software without restriction, including without limitation   /
// the rights to use, copy, modify, merge, publish, distribute, sublicense,    /
// and/or sell copies of the Software, and to permit persons to whom the       /
// Software is furnished to do so, subject to the following conditions:        /
//                                                                             /
// The above copyright notice and this permission notice shall be included in  /
// all copies or substantial portions of the Software.                         /
//                                                                             /
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR  /
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,    /
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE /
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER      /
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING     /
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER         /
// DEALINGS IN THE SOFTWARE.                                                   /
////////////////////////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

/*The goal of this code is to roughly discern the relative compute time of modular multiplication verses ElGamal.  Given the concern with timing, simple to understand implementations have been chosen over cryptographically valid implementations. */

//Group Size
const BitLen = 4096
const ByteLen = BitLen>>3

//Number of iterations on timeing test
const nOpsMul = 10000
const nOpsExp = 1000

func main(){
	//Strong 4096 bit Prime from https://tools.ietf.org/html/rfc3526#page-5
	primeString := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
      "29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
      "EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
      "E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
      "EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
      "C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
      "83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
      "670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B" +
      "E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9" +
      "DE2BCBF6955817183995497CEA956AE515D2261898FA0510" +
      "15728E5A8AAAC42DAD33170D04507A33A85521ABDF1CBA64" +
      "ECFB850458DBEF0A8AEA71575D060C7DB3970F85A6E1E4C7" +
      "ABF5AE8CDB0933D71E8C94E04A25619DCEE3D2261AD2EE6B" +
      "F12FFA06D98A0864D87602733EC86A64521F2B18177B200C" +
      "BBE117577A615D6C770988C0BAD946E208E24FA074E5AB31" +
      "43DB5BFCE0FD108E4B82D120A92108011A723C12A787E6D7" +
      "88719A10BDBA5B2699C327186AF4E23C1A946834B6150BDA" +
      "2583E9CA2AD44CE8DBBBC2DB04DE8EF92E8EFC141FBECAA6" +
      "287C59474E6BC05D99B2964FA090C3A2233BA186515BE7ED" +
      "1F612970CEE2D7AFB81BDD762170481CD0069127D5B05AA9" +
      "93B4EA988D8FDDC186FFB7DC90A6C08F4DF435C934063199" +
      "FFFFFFFFFFFFFFFF"

	//Create new prime
	p := big.NewInt(0)
	p.SetString(primeString, 16)
	
	psub1 := big.NewInt(0).Sub(p,big.NewInt(1))
	
	//generator
	g := big.NewInt(2)
	
	//Generate input keys for multiplication
	r := rand.New(rand.NewSource(42))
	var inputA []*big.Int
	var inputB []*big.Int
	var inputC []*big.Int
	
	for i := 0; i < nOpsMul; i++ {
		nint := RNGinPrime(r,psub1)
		inputA = append(inputA,nint)
		mint := RNGinPrime(r,psub1)
		inputB = append(inputB,mint)
		kint := RNGinPrime(r,psub1)
		inputC = append(inputC,kint)
	}
	
	// Compute total amount of time the multiplications take
	// post multiplication.
	timeStartMul := time.Now()
	for i := 0; i < nOpsMul; i++ {
		ModMul(inputA[i],inputB[i],p)

	}
	timeEndMul := time.Now()
	
	//Compute total amount of time the exponentiations take
	timeStartExp := time.Now()
	for i := 0; i < nOpsExp; i++ {
		ElgamalEncrypt(inputA[i],inputB[i],g,p,inputC[i])
	}
	timeEndExp := time.Now()
	
	//Output result
	timeElapsedMul := timeEndMul.Sub(timeStartMul)
	timePerOpMul := int(timeElapsedMul)/nOpsMul
	fmt.Printf("multiplication: %v ns/op \n", timePerOpMul)
	
	timeElapsedExp := timeEndExp.Sub(timeStartExp)
	timePerOpExp := int(timeElapsedExp)/nOpsExp
	fmt.Printf("exponentiation: %v ns/op \n", timePerOpExp)
	
	factor := float32(timePerOpExp)/float32(timePerOpMul)
	fmt.Printf("diferential factor: %v \n", factor)	
}

// Implementation of ElGamal encryption function based on https://en.wikipedia.org/wiki/ElGamal_encryption
// Is not secure, golang's big.Int is not constant time.  This implementation should never be used beyond
// this test.
func ElgamalEncrypt(x, y, g, p, m *big.Int)(*big.Int, *big.Int){
	c1 := big.NewInt(0).Exp(g,y,p)
	s := big.NewInt(0).Exp(c1,x,p)
	mdots := big.NewInt(0).Mul(s,m)
	c2 := big.NewInt(0).Mod(mdots,p)
	
	return c1, c2
}

// Implementation of modular multiplication 
// inefficent implentation of Modular Multiplication due to doing modular reduction 
// post full multiplication
func ModMul(a,b,p *big.Int)*big.Int{
	z := big.NewInt(0).Mul(a,b)
	q := big.NewInt(0).Mod(z, p)
	return q
}


//Brute force psudo-random number generation within cyclic group
func RNGinPrime(r *rand.Rand, pmin1 *big.Int)(*big.Int){

	cmpPrime := 1
	
	nint := big.NewInt(0)
	
	for cmpPrime>0 {
		byteField := make([]byte, ByteLen)
		r.Read(byteField)
		nint.SetBytes(byteField)
		cmpPrime = nint.Cmp(pmin1)
	}

	return nint
}
