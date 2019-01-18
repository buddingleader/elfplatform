package main

import (
	"context"
	"flag"
	"io"
	"time"

	pb "github.com/elforg/elfplatform/protos/common"
	"github.com/hyperledger/fabric/common/flogging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
	logger             = flogging.MustGetLogger("menu.client")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			logger.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		logger.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewMenuServiceClient(conn)
	addMenu(client, &pb.Menu{
		Enable: true,
		Value: &pb.MenuValue{
			Metadata: &pb.Metadata{},
		},
	})
	getMenu(client, &pb.Menu{
		Enable: true,
		Value: &pb.MenuValue{
			Metadata: &pb.Metadata{Id: 1},
		},
	})
}

func addMenu(client pb.MenuServiceClient, menu *pb.Menu) {
	logger.Debug("Add menu")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.AddMenu(ctx)
	if err != nil {
		logger.Fatalf("%v.AddMenu(_) = _, %v", client, err)
	}

	if err := stream.Send(menu); err != nil {
		logger.Fatalf("%v.Send(%v) = %v", stream, menu, err)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		logger.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}

	logger.Debugf("add menu result: %v", reply)
}

func getMenu(client pb.MenuServiceClient, menu *pb.Menu) {
	logger.Debug("Get menu")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.GetMenu(ctx)
	if err != nil {
		logger.Fatalf("%v.AddMenu(_) = _, %v", client, err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			m, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}

			if err != nil {
				logger.Fatalf("Failed to receive a menu : %v", err)
			}

			logger.Debugf("Got menu %s at id(%d)", m, menu.Value.Metadata.Id)
		}
	}()

	if err := stream.Send(menu); err != nil {
		logger.Fatalf("Failed to send a menu: %v", err)
	}
	stream.CloseSend()
	<-waitc
}
