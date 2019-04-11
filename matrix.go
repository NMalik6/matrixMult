package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// runtime.GOMAXPROCS(1)
	done := make(chan bool)
	m1 := makeMatrix(1000, 1000)
	m2 := makeMatrix(1000, 1000)
	go m1.randomize(done)
	time.Sleep(100 * time.Millisecond)
	go m2.randomize(done)
	for i := 0; i < 2; i++ {
		<-done
	}
	// m1.printMatrix()
	// m2.printMatrix()
	start := time.Now()
	product := matrixMult(m1, m2)
	// product := origMult(m1, m2, 1000)
	elapsed := time.Since(start)
	product.printMatrix()
	fmt.Println(elapsed)
}

type matrix [][]int

func makeMatrix(rows int, cols int) matrix {
	var m matrix
	for i := 0; i < rows; i++ {
		a := make([]int, cols)
		m = append(m, a)
	}
	return m
}

func dotMult(m int, a []int) []int {
	for i := range a {
		a[i] *= m
	}
	return a
}

func dotAdd(a1, a2 []int) []int {
	for i := range a1 {
		a1[i] += a2[i]
	}
	return a1
}

func origMult(m1, m2 matrix, n int) matrix {
	mat := makeMatrix(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			mat[i][j] = 0
			for k := 0; k < n; k++ {
				mat[i][j] += m1[i][k] * m2[k][j]
			}
		}
	}
	return mat
}

func matrixMult(m1, m2 matrix) matrix {
	var mat matrix
	done := make(chan bool)
	for i := 0; i < len(m1[0]); i++ {
		in := i
		go func() {
			oldHat := m1.col(in)
			a := make([]int, len(oldHat))
			for j := 0; j < len(oldHat); j++ {
				a = dotAdd(a, dotMult(oldHat[j], m2.col(j)))
			}
			mat = append(mat, a)
			done <- true
		}()
	}
	for i := 0; i < len(m1[0]); i++ {
		<-done
	}
	var final matrix
	for i := 0; i < len(mat[0]); i++ {
		final = append(final, mat.col(i))
	}
	return final
}

func (m matrix) randomize(done chan<- bool) {
	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[0]); j++ {
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i) + int64(j)))
			m[i][j] = r.Intn(10)
		}
	}
	done <- true
}

func (m matrix) row(n int) []int {
	return m[n]
}

func (m matrix) col(n int) []int {
	a := []int{}
	for i := 0; i < len(m); i++ {
		a = append(a, m[i][n])
	}
	return a
}

func (m matrix) printMatrix() {
	for i := 0; i < len(m); i++ {
		fmt.Println(m.row(i))
	}
	fmt.Println()
}
