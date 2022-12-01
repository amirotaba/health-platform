package app

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"

	"git.paygear.ir/giftino/account/internal/account/domain"
	"git.paygear.ir/giftino/account/internal/app/config"
	"git.paygear.ir/giftino/account/internal/app/database/migrate"
	"git.paygear.ir/giftino/account/internal/app/grpc"
	"git.paygear.ir/giftino/account/internal/app/http/echo"
	accountRepository "git.paygear.ir/giftino/account/internal/features/account/repository/mysql"
	accountUseCase "git.paygear.ir/giftino/account/internal/features/account/usecase"
	accountRoleRepository "git.paygear.ir/giftino/account/internal/features/account_role/repository/mysql"
	accountRoleUsecase "git.paygear.ir/giftino/account/internal/features/account_role/usecase"
	accountTypeRepository "git.paygear.ir/giftino/account/internal/features/account_type/repository/mysql"
	accountTypeUsecase "git.paygear.ir/giftino/account/internal/features/account_type/usecase"
	authUseCase "git.paygear.ir/giftino/account/internal/features/authentication/usecase"
	channelRepository "git.paygear.ir/giftino/account/internal/features/channel/repository/mysql"
	channelUsecase "git.paygear.ir/giftino/account/internal/features/channel/usecase"
	channelAccountRepository "git.paygear.ir/giftino/account/internal/features/channel_account/repository/mysql"
	channelAccountUsecase "git.paygear.ir/giftino/account/internal/features/channel_account/usecase"
	channelRuleRepository "git.paygear.ir/giftino/account/internal/features/channel_rule/repository/mysql"
	channelRuleUsecase "git.paygear.ir/giftino/account/internal/features/channel_rule/usecase"
	otpRepository "git.paygear.ir/giftino/account/internal/features/otp/repository/mysql"
	otpUseCase "git.paygear.ir/giftino/account/internal/features/otp/usecase"
	permissionRepository "git.paygear.ir/giftino/account/internal/features/permission/repository/mysql"
	permissionUsecase "git.paygear.ir/giftino/account/internal/features/permission/usecase"
	permissionServiceRepository "git.paygear.ir/giftino/account/internal/features/permission_service/repository/mysql"
	permissionServiceUsecase "git.paygear.ir/giftino/account/internal/features/permission_service/usecase"
	roleRepository "git.paygear.ir/giftino/account/internal/features/role/repository/mysql"
	roleUsecase "git.paygear.ir/giftino/account/internal/features/role/usecase"
	rolePermissionRepository "git.paygear.ir/giftino/account/internal/features/role_permission/repository/mysql"
	rolePermissionUsecase "git.paygear.ir/giftino/account/internal/features/role_permission/usecase"
	serviceRepository "git.paygear.ir/giftino/account/internal/features/service/repository/mysql"
	serviceUsecase "git.paygear.ir/giftino/account/internal/features/service/usecase"
	tokenRepository "git.paygear.ir/giftino/account/internal/features/token/repository/mysql"
	tokenUseCase "git.paygear.ir/giftino/account/internal/features/token/usecase"
	//authService "git.paygear.ir/giftino/account/internal/service/auth"
	natsBus "git.paygear.ir/giftino/account/internal/service/bus/nats"
	inventoryGrpc "git.paygear.ir/giftino/account/internal/service/inventory/grpc"
	jwtService "git.paygear.ir/giftino/account/internal/service/jwt"
	jwtMiddleWare "git.paygear.ir/giftino/account/internal/service/middleware"
	mySqlStore "git.paygear.ir/giftino/account/internal/service/store/mysql"
	walletGrpc "git.paygear.ir/giftino/account/internal/service/wallet/grpc"
	RedisUtils "git.paygear.ir/giftino/account/internal/utils/redis"
)

func Run() {
	conf := config.LoadEnv()
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	logger := log.New()
	fmt.Println(conf.MySqlDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		conf.MySqlRootAccount, conf.MySqlRootPassword, conf.MySqlHost,
		conf.MysqlPort, conf.MySqlDataBaseName))
	db, err := sql.Open(conf.MySqlDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		conf.MySqlRootAccount, conf.MySqlRootPassword, conf.MySqlHost,
		conf.MysqlPort, conf.MySqlDataBaseName))
	if err != nil {
		log.Println(err)
		//panic[string]("error in get db connection")
	}
	//
	//err = migrate.Up(conf)
	//if err != nil {
	//	log.Println(err)
	//	//panic[string](err.Error())
	//}

	up, err := strconv.ParseBool(conf.MigrationUp)
	if up {
		err = migrate.Up(conf)
		if err != nil {
			log.Println(err)
			//panic(err)
		}
	}

	down, err := strconv.ParseBool(conf.MigrationDown)
	if down {
		err = migrate.Down(conf)
		if err != nil {
			log.Println(err)
			//panic(err)
		}
	}

	//nc, err := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", conf.NatsHost, conf.NatsPort))
	if err != nil {
		log.Println(err)
		//panic[string](err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.RedisHost, conf.RedisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisServer := RedisUtils.New(rdb)
	bus := natsBus.New(nc)
	mySqlStoreService := mySqlStore.NewStore(db)
	jwtServices := jwtService.NewJWTService()

	serviceRepo := serviceRepository.New(db, logger)
	serviceUseCase := serviceUsecase.New(serviceRepo, logger)

	permissionServiceRep := permissionServiceRepository.New(db, logger)
	permissionServiceUseCase := permissionServiceUsecase.New(permissionServiceRep, logger)

	permissionRepo := permissionRepository.New(db, logger)
	permissionUseCase := permissionUsecase.New(permissionRepo, serviceRepository.New(db, logger), permissionServiceRep, logger)

	roleRepo := roleRepository.New(db, logger)
	rolePerRepo := rolePermissionRepository.New(db, logger)
	roleUseCase := roleUsecase.New(roleRepo, permissionRepo, rolePerRepo, logger)
	rolePerUsecase := rolePermissionUsecase.New(rolePerRepo, logger)

	accountRepo := accountRepository.New(db, logger)
	accountUsecase := accountUseCase.New(accountRepo, roleRepo, permissionRepo, accountTypeRepository.New(db, logger), mySqlStoreService, bus, logger)

	channelRepo := channelRepository.New(db, logger)
	channelRuleRepo := channelRuleRepository.New(db, logger)
	//jwtServices := jwtServices.NewJWTService(domain.GiftinoSecretKey, "account_service")
	//authManagerService := authService.NewAuthManager(jwtServices)

	accountRoleRepo := accountRoleRepository.New(db, logger)
	accountRoleUse := accountRoleUsecase.New(accountRoleRepo, logger)
	authUsecase := authUseCase.New(accountRepo, redisServer, jwtServices, rolePerUsecase, permissionServiceUseCase, serviceRepo, roleRepo, accountUsecase, accountRoleUse)
	jwtMiddleWareService := jwtMiddleWare.New(authUsecase) //New(accountRepo, roleRepo)

	channelUseCase := channelUsecase.New(channelRepo, accountUsecase, channelAccountRepository.New(db, logger), walletGrpc.New(conf.GrpcWalletAddress), bus, logger)
	channelUseCase.EventHandler(nc)
	channelRuleUseCase := channelRuleUsecase.New(channelRuleRepo, inventoryGrpc.New(conf.GrpcInventoryAddress), channelRepo)

	accountTypeRepo := accountTypeRepository.New(db, logger)
	accountTypeUseCase := accountTypeUsecase.New(accountTypeRepo)

	caRepo := channelAccountRepository.New(db, logger)
	caUsecase := channelAccountUsecase.New(caRepo, channelRepo, accountRepo, channelUseCase, roleUseCase, logger)
	tokenRepo := tokenRepository.New(db, logger)
	tokenUsecase := tokenUseCase.New(tokenRepo, jwtServices, redisServer)

	otpUsecase := otpUseCase.New(otpRepository.New(db), tokenUsecase, accountRepo, roleRepo, rolePerUsecase, accountRoleUse, accountTypeRepo, redisServer)

	service := domain.AccountServices{
		AccountUsecase:     accountUsecase,
		ChannelUsecase:     channelUseCase,
		PermissionUsecase:  permissionUseCase,
		RoleUsecase:        roleUseCase,
		ServiceUsecase:     serviceUseCase,
		AccountTypeUsecase: accountTypeUseCase,
		ChannelRuleUsecase: channelRuleUseCase,
		AuthUsecase:        authUsecase,
		OtpUsecase:         otpUsecase,
		TokenUsecase:       tokenUsecase,
		InventoryGrpc:      inventoryGrpc.New(conf.GrpcInventoryAddress),
		Middleware:         jwtMiddleWareService,
		HttpPort:           conf.HttpAppServerPort,
		GrpcPort:           conf.GrpcAppServerPort,
		RolePermission:     rolePerUsecase,
		PermissionServices: permissionServiceUseCase,
		ChannelAccount:     caUsecase,
		Logger:             logger,
	}

	//if _, err = nc.Subscribe("channel", func(m *nats.Msg) {
	//	//var baseEvent event.ChannelChargedEvent
	//	var baseEvent event.BaseEvent
	//	err := json.Unmarshal(m.Data, &baseEvent)
	//	if err != nil {
	//		log.Println(err)
	//		return
	//	}
	//
	//	log.Println(string(m.Data))
	//	log.Println(baseEvent.Data)
	//	switch baseEvent.GetEventType() {
	//	case event.ChannelChargedEventType:
	//		fmt.Println("***************************************************")
	//		var charge entity.ChargeChannelRequest
	//		log.Println(json.Unmarshal([]byte(baseEvent.Data), &charge))
	//		log.Println(charge.ID)
	//		log.Println(charge.CurrentBalance)
	//		updateReq := entity.NewUpdateCollectionRequest(domain.ChannelTable, baseEvent.Data)
	//		updateReq.Filters = append(updateReq.Filters, entity.Filter{
	//			TableName: domain.ChannelTable,
	//			Field:     "id",
	//			Value:     []interface{}{charge.ID},
	//		})
	//		err = channelRepo.UpdateChannelNew(context.Background(), updateReq)
	//		if err != nil {
	//			logger.Println(err)
	//		}
	//	}
	//}); err != nil {
	//	return
	//}

	// add account grpc port here
	go grpc.Run(service)
	echo.New(service)
}
