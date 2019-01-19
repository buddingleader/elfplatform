package main

import (
	"context"
	"flag"
	"time"

	pb "github.com/elforg/elfplatform/protos/common"
	"github.com/golang/protobuf/ptypes/timestamp"
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
	client := pb.NewCRUDClient(conn)
	menu1 := &pb.Menu{
		Metadata: &pb.Metadata{
			Id:               1,
			Name:             "child1",
			CreateAuthorId:   0,
			CreateAuthorName: "Admin",
			Created: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   0,
			},
			UpdateAuthorId:   0,
			UpdateAuthorName: "Admin",
			LastUpdated: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   0,
			},
		},
		Enable:   true,
		Policy:   &pb.Policy{},
		SubMenus: make(map[int32]string, 0),
	}
	menu2 := &pb.Menu{
		Metadata: &pb.Metadata{
			Id:               2,
			Name:             "child2",
			CreateAuthorId:   0,
			CreateAuthorName: "Admin",
			Created: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   0,
			},
			UpdateAuthorId:   0,
			UpdateAuthorName: "Admin",
			LastUpdated: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   0,
			},
		},
		Enable:   true,
		Policy:   &pb.Policy{},
		SubMenus: make(map[int32]string, 0),
	}
	addMenu(client, 100000, menu1)
	addMenu(client, 100000, menu2)
	getMenu(client, menu1.Metadata.Id)
	getMenu(client, menu2.Metadata.Id)
}

func addMenu(client pb.CRUDClient, id int32, menu *pb.Menu) {
	logger.Debug("Add menu")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.MenuRequest{
		RootMenuId: id,
		Menu:       menu,
	}
	result, err := client.AddMenu(ctx, req)
	if err != nil {
		logger.Fatalf("%v.AddMenu(req: %v) = _, %v", client, req, err)
	}

	logger.Debugf("add menu result: %v", result)
}

func getMenu(client pb.CRUDClient, id int32) {
	logger.Debug("Get menu")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.MenuRequest{
		RootMenuId: id,
	}
	result, err := client.GetMenu(ctx, req)
	if err != nil {
		logger.Fatalf("%v.GetMenu(req: %v) = _, %v", client, req, err)
	}

	logger.Debugf("get menu result: %v", result)
}
