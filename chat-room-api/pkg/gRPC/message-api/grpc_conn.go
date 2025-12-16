package messageapi

import (
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn

func openConnWithoutTLS() error {
	connection, err := grpc.NewClient(os.Getenv("MESSAGE_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	conn = connection
	return nil
}

func GetgRPCServer() *grpc.ClientConn {
	if conn != nil {
		return conn
	}
	_ = openConnWithoutTLS()
	return conn
}

func CloseConn() error {
	return conn.Close()
}
