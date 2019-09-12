package main

import "github.com/PitneyBowes/compose-chart-api/pic"

type client struct{}

func (*client) NewBuilder(c *pic.Config) pic.Builder {
	return newBuilder(c)
}

func init() {
	pic.SetClient(
		&client{},
		pic.Options{
			LogLevel:    pic.LogInfo,
			LogFileName: "go-chart.log",
		},
	)
}

func main() {
}
