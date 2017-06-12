package main

import "fmt"

func main() {

	// 初始资金
	const s float64 = 1000

	// 余额
	var m float64 = s

	// 每月定投
	var d float64 = 100

	// 持有数量
	var g float64 = 0

	// 每月价位
	points := []float64{
		100,
		110,
		120,
		130,
		100,
		70,
		100,
		140,
		130,
		120,
		120,
	}

	// 每月定投
	for i, p := range points {
		// 最后一把不买
		if i == len(points)-1 {
			continue
		}
		g += d / p
		m -= d
	}

	// 定投持股数和盈利
	fmt.Println(g)
	m = g * points[len(points) - 1]
	g = 0
	fmt.Println(m - s)

	fmt.Println()

	// 全投持股数和盈利
	g2 := ( ( d * float64((len(points) -1)) ) / points[0] )
	fmt.Println(g2)
	fmt.Println((g2 * points[len(points) - 1]) - s)

}
