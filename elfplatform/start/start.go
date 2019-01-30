package start

import (
	"fmt"
	"net"

	"github.com/spf13/viper"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the elf platform",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Elf platform start -- Version 1.0")
		listen, port := viper.GetString("Listen"), viper.GetInt("Port")
		serve(listen, port)
	},
}

var (
	tls      bool
	certFile string
	keyFile  string
	logger   = flogging.MustGetLogger("start")
)

// Cmd returns the cobra command for Chaincode
func Cmd() *cobra.Command {
	addFlags()
	return startCmd
}

func addFlags() {
	startCmd := startCmd.PersistentFlags()
	startCmd.BoolVarP(&tls, "tls", "t", false, "Connection uses TLS if true, else plain TCP")
	startCmd.StringVarP(&certFile, "cert_file", "c", "", "The TLS cert file")
	startCmd.StringVarP(&keyFile, "key_file", "k", "", "The TLS key file")
}

func serve(listen string, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listen, port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if tls {
		if certFile == "" {
			certFile = testdata.Path("server1.pem")
		}
		if keyFile == "" {
			keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logger.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	// pb.RegisterCRUDServer(grpcServer, newServer())
	registerServers(grpcServer)
	grpcServer.Serve(lis)
}

func registerServers(grpcServer *grpc.Server) {
	// pb.RegisterCRUDServer(grpcServer, newServer())

}
