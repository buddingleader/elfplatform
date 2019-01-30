package menu

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/elforg/elfplatform/protos/common"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/common/util"
	"github.com/pkg/errors"
)

// Server to manager crud
type Server struct {
	menuMap map[int32]*pb.Menu
	count   int32
	rwmutex sync.RWMutex
	logger  *flogging.FabricLogger
}

// NewServer returns a menu server
func NewServer() *Server {
	// get from db
	root := initialRootMenu()
	s := &Server{
		menuMap: make(map[int32]*pb.Menu),
		logger:  flogging.MustGetLogger("menu.server"),
	}
	s.menuMap[root.Metadata.Id] = root
	s.count++

	s.logger.Debugf("the menuMap is %v", s.menuMap)
	s.logger.Info("Start menu service")
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

func (ms *Server) AddMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
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

func (ms *Server) UpdateMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
	logger.Debug("Update Menu")

	return nil, nil
}

func (ms *Server) DeleteMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
	logger.Debug("Delete Menu")

	return nil, nil
}

func (ms *Server) GetMenu(ctx context.Context, req *pb.MenuRequest) (menu *pb.Menu, err error) {
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
