# sort
External memory file sorter linke unix `sort`

Basic idea split file into chunks, sort, and merge like merge operation in merge sort algorithm

# Tests

Generate payload data
`go run generator/main.go -count=1000 -length=1000 > out`

Then compare
`cat out  | go run main.go  > sorted`
`cat out  | sort  > sorted_ref`
`diff sorted sorted_ref`
