package main

import (
	"fmt"
	"os"

	ptvvisum "github.com/lddl/go-ptv-visum"
	"github.com/lddl/go-ptv-visum/roadnet"
	"github.com/lddl/go-ptv-visum/utils"
)

func main() {
	file, err := os.Open("./example/sample/example.net")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	ptvData, err := ptvvisum.ReadPTVFromFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	roadNetwork, err := roadnet.ExtractGraph(ptvData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Vertices (first 5):")
	fmt.Println("id;geom")
	count := 0
	for _, vertex := range roadNetwork.Vertices {
		fmt.Printf("%d;%s\n", vertex.ID, utils.PointToWKT([]float64{vertex.X, vertex.Y}))
		if count >= 5 {
			break
		}
		count++
	}
	fmt.Println("\nEdges (first 5):")
	fmt.Println("id;source;target;geom")
	count = 0
	for _, edge := range roadNetwork.Edges {
		fmt.Printf("%d;%d;%d;%s\n", edge.ID, edge.Source, edge.Target, utils.LineStringToWKT(edge.Geometry))
		if count >= 5 {
			break
		}
		count++
	}
}
