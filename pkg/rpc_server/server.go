package rpc_server

import (
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func Server(port string) error {

	err := rpc.RegisterName("Reply", new(Reply))
	if err != nil {
		log.Fatal("Failed to Register RPC server!")
	}

	http.HandleFunc("/participle", func(w http.ResponseWriter, r *http.Request) {
		conn := connection{Writer: w, ReadCloser: r.Body}
		err = rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
		if err != nil {
			log.Fatal(err)
		}
	})
	return http.ListenAndServe(port, nil)
}
