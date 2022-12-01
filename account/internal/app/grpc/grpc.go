package grpc

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"git.paygear.ir/giftino/account/internal/account/domain"
	pbChannel "git.paygear.ir/giftino/account/internal/account/proto/channel"
	authGrpc "git.paygear.ir/giftino/account/internal/features/authentication/delivery/grpc"
	channelGrpc "git.paygear.ir/giftino/account/internal/features/channel/delivery/grpc"
	channelRuleGrpc "git.paygear.ir/giftino/account/internal/features/channel_rule/delivery/grpc"
	serviceGrpc "git.paygear.ir/giftino/account/internal/features/service/delivery/grpc"
)

func Run(services domain.AccountServices) {
	//conf := config.LoadEnv()
	lis, err := net.Listen("tcp", services.GrpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//db, err := sql.Open(conf.MySqlDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
	//	conf.MySqlRootAccount, conf.MySqlRootPassword, conf.MySqlHost,
	//	conf.MysqlPort, conf.MySqlDataBaseName))
	//if err != nil {
	//	log.Println(err)
	//	//panic[string]("error in get db connection")
	//}

	//nc, err := nats.Connect(nats.DefaultURL)
	//nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", conf.NatsHost, conf.NatsPort))
	//if err != nil {
	//	//panic[string](err.Error())
	//}
	//
	//bus := natsBus.New(nc)
	//mySqlStoreService := mySqlStore.NewStore(db)
	//
	//log.Println(db)
	//log.Println(err)
	//tagRepository := tagRepo.New(db)
	//tagUseCase := tagUsecase.New(tagRepository, mySqlStoreService, bus)
	//
	s := grpc.NewServer()
	authGrpc.RegisterAuthenticationServiceServer(s, authGrpc.New(services.AuthUsecase))
	serviceGrpc.RegisterServiceServiceServer(s, serviceGrpc.New(services.ServiceUsecase))
	pbChannel.RegisterChannelServiceServer(s, channelGrpc.New(services.ChannelUsecase))
	channelRuleGrpc.RegisterChannelRuleServiceServer(s, channelRuleGrpc.New(services.ChannelRuleUsecase))
	log.Printf("grpc server listening at %v \n", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
