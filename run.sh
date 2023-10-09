#!/bin/bash

go run useBEB.go 0 127.0.0.1:5000  127.0.0.1:6001  127.0.0.1:7002 >> useBEB.out &
P1=$!
go run useBEB.go 1 127.0.0.1:5000  127.0.0.1:6001  127.0.0.1:7002 >> useBEB.out &
P2=$!
go run useBEB.go 2 127.0.0.1:5000  127.0.0.1:6001  127.0.0.1:7002 >> useBEB.out &
P3=$!
wait $P1 $P2 $P3

# go run useBEB.go 0 127.0.0.1:5000  127.0.0.1:6001 &
# P1=$!
# go run useBEB.go 1 127.0.0.1:5000  127.0.0.1:6001 &
# P2=$!
# wait $P1 $P2