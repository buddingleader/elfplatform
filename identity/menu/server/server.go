package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"sync"

	pb "github.com/elforg/elfplatform/protos/common"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
	logger     = flogging.MustGetLogger("menu.server")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			logger.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMenuServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

type menuServiceServer struct {
	menus    map[int32]*pb.Menu
	sequence int32
	mutex    sync.Mutex // protects routeNotes
}

func newServer() *menuServiceServer {
	s := &menuServiceServer{
		menus: make(map[int32]*pb.Menu),
	}

	logger.Info("Start menu server")
	return s
}

func (ms *menuServiceServer) AddMenu(stream pb.MenuService_AddMenuServer) error {
	logger.Debug("Add menu start...")
	for {
		menu, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("Add menu end")
			return stream.SendAndClose(&pb.SimpleReply{
				Success: true,
			})
		}

		if err != nil {
			logger.Error(err)
			return errors.WithStack(err)
		}

		if ok, err := checkMenu(menu); err != nil && !ok {
			logger.Error(err)
			return errors.WithStack(err)
		}

		ms.mutex.Lock()
		ms.sequence++
		menu.Value.Metadata.Id = ms.sequence
		ms.menus[ms.sequence] = menu
		ms.mutex.Unlock()
		// save the menu

		logger.Infof("Add menu:%v", menu)
	}
}

func (ms *menuServiceServer) UpdateMenu(stream pb.MenuService_UpdateMenuServer) error {
	logger.Debug("Update Menu")

	return nil
}

func (ms *menuServiceServer) GetMenu(stream pb.MenuService_GetMenuServer) error {
	logger.Debug("Get Menu")
	for {
		menu, err := stream.Recv()
		if err == io.EOF {
			logger.Debug("Get menu end")
			return nil
		}
		if err != nil {
			return errors.WithStack(err)
		}

		if ok, err := checkMenu(menu); err != nil && !ok {
			logger.Error(err)
			return errors.WithStack(err)
		}

		if menu.Value.Metadata.GetId() == 0 {
			err := errors.New("Menu' Id is nil")
			logger.Error(err)
			return errors.WithStack(err)
		}
		logger.Debugf("Get menu:%v", menu)
		if err := stream.Send(ms.menus[menu.Value.Metadata.Id]); err != nil {
			return errors.WithStack(err)
		}
	}
}

func checkMenu(menu *pb.Menu) (bool, error) {
	if menu.GetValue() == nil {
		return false, errors.New("Menu' Value is nil")
	}

	if menu.Value.GetMetadata() == nil {
		return false, errors.New("Menu' Metadata is nil")
	}

	return true, nil
}
