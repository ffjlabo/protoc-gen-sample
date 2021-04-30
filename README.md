# protoc-gen-sample

// create protoc-gen-sample binary  
`go build -o protoc-gen-sample main.go`

// execute protoc with protoc-gen-sample plugin  
`protoc -I. --plugin=./protoc-gen-sample --sample_out=. sample.proto`
