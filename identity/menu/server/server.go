package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	pb "github.com/elforg/elfplatform/protos/common"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/util"
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
	pb.RegisterCRUDServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}

type menuServer struct {
	menuTree *pb.MenuTree
	menuMap  map[int32]*pb.Menu
	count    int32
	rwmutex  sync.RWMutex // protects routeNotes
}

func newServer() *menuServer {
	// get from db
	root := initialRootMenu()
	s := &menuServer{
		menuMap: make(map[int32]*pb.Menu),
	}
	s.menuMap[root.Metadata.Id] = root
	s.count++

	logger.Debugf("the menuMap is %v", s.menuMap)
	logger.Info("Start menu server")
	return s
}

func initialRootMenu() *pb.Menu {
	return &pb.Menu{
		Metadata: &pb.Metadata{
			Id:               100000,
			Name:             "root",
			CreateAuthorId:   100000,
			CreateAuthorName: "Admin",
			Created: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   0,
			},
			UpdateAuthorId:   100000,
			UpdateAuthorName: "Admin",
			LastUpdated: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   0,
			},
		},
		Enable:   true,
		Policy:   &pb.Policy{},
		SubMenus: make(map[int32]string),
	}
}

func (ms *menuServer) AddMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
	addr := util.ExtractRemoteAddress(ctx)
	logger.Debugf("Connection from %s, start to add menu, request: %v", addr, req)
	defer logger.Debugf("Closing connection from %s, result: %v", addr, menu)

	ms.rwmutex.RLock()
	rootMenu, has := ms.menuMap[req.RootMenuId]
	ms.rwmutex.RUnlock()
	if !has {
		return nil, fmt.Errorf("cannot found this root menu[id=%v]", req.RootMenuId)
	}
	logger.Debugf("rootMenu is %v", rootMenu)

	menu = req.Menu
	if err = menu.CheckMetadata(); err != nil {
		logger.Error(err)
		return nil, errors.WithStack(err)
	}

	// check if it exists
	if name, has := rootMenu.SubMenus[menu.Metadata.Id]; has {
		err = fmt.Errorf("this menu already exists, id=%v, name=%v", menu.Metadata.Id, name)
		logger.Error(err)
		return nil, errors.WithStack(err)
	}

	ms.rwmutex.Lock()
	defer ms.rwmutex.Unlock()
	rootMenu.SubMenus[menu.Metadata.Id] = menu.Metadata.Name
	ms.menuMap[rootMenu.Metadata.Id] = rootMenu
	ms.menuMap[menu.Metadata.Id] = menu
	ms.count++

	// TODO: save to db
	// Put(rootMenu.Metadata.Id, rootMenu)
	// Put(menu.Metadata.Id, menu)

	logger.Debugf("the new root menu is %v", rootMenu)
	return menu, nil
}

func (ms *menuServer) UpdateMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
	logger.Debug("Update Menu")

	return nil, nil
}

func (ms *menuServer) DeleteMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
	logger.Debug("Delete Menu")

	return nil, nil
}

func (ms *menuServer) GetMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
	addr := util.ExtractRemoteAddress(ctx)
	logger.Debugf("Connection from %s, start to get menu, request: %v", addr, req)
	defer logger.Debugf("Closing connection from %s, result: %v", addr, menu)

	// check if it exists
	ms.rwmutex.RLock()
	menu, has := ms.menuMap[req.RootMenuId]
	ms.rwmutex.RUnlock()
	if !has {
		err = fmt.Errorf("this menu doesn't exists, id=%v", menu.Metadata.Id)
		logger.Error(err)
		return nil, errors.WithStack(err)
	}

	return menu, nil
}
