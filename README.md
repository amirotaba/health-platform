## How to run
run the following command to migrate tables
```
cd services && go run main.go migrate
```
and then run
```
cd account && go run main.go rest && cd .. && cd services && go run main.go start
```
### Project structure
```
├── account
│   ├── cmd
│   │   ├── rest.go
│   │   └── root.go
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── account
│   │   │   ├── aggregate
│   │   │   │   ├── account.go
│   │   │   │   ├── aggregate.go
│   │   │   │   ├── channel.go
│   │   │   │   ├── channel_Rules.go
│   │   │   │   └── permissin-realm.go
│   │   │   ├── command
│   │   │   │   ├── account.go
│   │   │   │   ├── channel.go
│   │   │   │   └── command.go
│   │   │   ├── domain
│   │   │   │   ├── account.go
│   │   │   │   ├── account_role.go
│   │   │   │   ├── account_type.go
│   │   │   │   ├── authentication.go
│   │   │   │   ├── auth_manager.go
│   │   │   │   ├── channel_account.go
│   │   │   │   ├── channel.go
│   │   │   │   ├── channel_rule.go
│   │   │   │   ├── config.go
│   │   │   │   ├── env.go
│   │   │   │   ├── errors.go
│   │   │   │   ├── inventory.go
│   │   │   │   ├── jwt.go
│   │   │   │   ├── key.go
│   │   │   │   ├── logger.go
│   │   │   │   ├── middleware.go
│   │   │   │   ├── otp.go
│   │   │   │   ├── paginate.go
│   │   │   │   ├── permission.go
│   │   │   │   ├── permission_service.go
│   │   │   │   ├── redis.go
│   │   │   │   ├── role.go
│   │   │   │   ├── role_permission.go
│   │   │   │   ├── service.go
│   │   │   │   ├── services.go
│   │   │   │   ├── tag.go
│   │   │   │   ├── token.go
│   │   │   │   └── wallet.go
│   │   │   ├── entity
│   │   │   │   ├── account.go
│   │   │   │   ├── account_role.go
│   │   │   │   ├── authentication.go
│   │   │   │   ├── channel_account.go
│   │   │   │   ├── channel.go
│   │   │   │   ├── channel_rule.go
│   │   │   │   ├── filter.go
│   │   │   │   ├── otp.go
│   │   │   │   ├── permission.go
│   │   │   │   ├── permission_role.go
│   │   │   │   ├── permission_service.go
│   │   │   │   ├── role.go
│   │   │   │   ├── service.go
│   │   │   │   ├── sms.go
│   │   │   │   ├── tag.go
│   │   │   │   └── token.go
│   │   │   ├── event
│   │   │   │   ├── account.go
│   │   │   │   ├── bus.go
│   │   │   │   ├── channel.go
│   │   │   │   ├── event.go
│   │   │   │   ├── handler.go
│   │   │   │   ├── otp.go
│   │   │   │   ├── store.go
│   │   │   │   └── token.go
│   │   │   └── proto
│   │   │       ├── account
│   │   │       │   ├── account_grpc.pb.go
│   │   │       │   ├── account.pb.go
│   │   │       │   └── account.proto
│   │   │       ├── auth
│   │   │       │   ├── auth_grpc.pb.go
│   │   │       │   ├── auth.pb.go
│   │   │       │   └── auth.proto
│   │   │       ├── channel
│   │   │       │   ├── channel_grpc.pb.go
│   │   │       │   ├── channel.pb.go
│   │   │       │   └── channel.proto
│   │   │       ├── event
│   │   │       │   ├── event.pb.go
│   │   │       │   └── event.proto
│   │   │       └── wallet
│   │   │           ├── wallet_grpc.pb.go
│   │   │           ├── wallet.pb.go
│   │   │           └── wallet.proto
│   │   ├── app
│   │   │   ├── app.go
│   │   │   ├── config
│   │   │   │   └── env.go
│   │   │   ├── database
│   │   │   │   └── migrate
│   │   │   │       ├── migrations
│   │   │   │       │   ├── 0001_initial.down.sql
│   │   │   │       │   ├── 0001_initial.up.sql
│   │   │   │       │   ├── 0002_index.up.sql
│   │   │   │       │   ├── 0003_otp_initial.down.sql
│   │   │   │       │   ├── 0003_otp_initial.up.sql
│   │   │   │       │   ├── 0004_token_initial.down.sql
│   │   │   │       │   ├── 0004_token_initial.up.sql
│   │   │   │       │   ├── 0005_service_initial.down.sql
│   │   │   │       │   ├── 0005_service_initial.up.sql
│   │   │   │       │   ├── 0006_channel_account_initial.down.sql
│   │   │   │       │   ├── 0006_channel_account_initial.up.sql
│   │   │   │       │   ├── 0007_service_account_initial.up.sql
│   │   │   │       │   ├── 0008_add_channel_walletid.up.sql
│   │   │   │       │   └── 0009_add_channel_to.up.sql
│   │   │   │       └── migrations.go
│   │   │   ├── grpc
│   │   │   │   └── grpc.go
│   │   │   └── http
│   │   │       └── echo
│   │   │           └── echo.go
│   │   ├── features
│   │   │   ├── account
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       ├── es-cqrs.go
│   │   │   │       └── usecase.go
│   │   │   ├── account_role
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── account_type
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── authentication
│   │   │   │   ├── delivery
│   │   │   │   │   ├── grpc
│   │   │   │   │   │   ├── auth_grpc.pb.go
│   │   │   │   │   │   ├── auth.pb.go
│   │   │   │   │   │   ├── auth.proto
│   │   │   │   │   │   └── delivery.go
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── channel
│   │   │   │   ├── delivery
│   │   │   │   │   ├── grpc
│   │   │   │   │   │   └── delivery.go
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── channel_account
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── channel_rule
│   │   │   │   ├── delivery
│   │   │   │   │   ├── grpc
│   │   │   │   │   │   ├── channel_rule_grpc.pb.go
│   │   │   │   │   │   ├── channel_rule.pb.go
│   │   │   │   │   │   ├── channel_rule.proto
│   │   │   │   │   │   └── delivery.go
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── otp
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── permission
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── permission_service
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── role
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── role_permission
│   │   │   │   ├── delivery
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   ├── service
│   │   │   │   ├── delivery
│   │   │   │   │   ├── grpc
│   │   │   │   │   │   ├── delivery.go
│   │   │   │   │   │   ├── realm_grpc.pb.go
│   │   │   │   │   │   ├── realm.pb.go
│   │   │   │   │   │   └── realm.proto
│   │   │   │   │   └── http
│   │   │   │   │       └── delivery.go
│   │   │   │   ├── repository
│   │   │   │   │   └── mysql
│   │   │   │   │       └── repository.go
│   │   │   │   └── usecase
│   │   │   │       └── usecase.go
│   │   │   └── token
│   │   │       ├── delivery
│   │   │       │   └── http
│   │   │       │       └── delivery.go
│   │   │       ├── repository
│   │   │       │   └── mysql
│   │   │       │       └── repository.go
│   │   │       └── usecase
│   │   │           └── usecase.go
│   │   ├── service
│   │   │   ├── bus
│   │   │   │   └── nats
│   │   │   │       └── bus.go
│   │   │   ├── config
│   │   │   │   └── viper_service.go
│   │   │   ├── event-handler
│   │   │   │   └── event_handler.go
│   │   │   ├── inventory
│   │   │   │   └── grpc
│   │   │   │       ├── grpc.go
│   │   │   │       ├── tag_grpc.pb.go
│   │   │   │       ├── tag.pb.go
│   │   │   │       └── tag.proto
│   │   │   ├── jwt
│   │   │   │   ├── jwt.go
│   │   │   │   └── jwtN.go
│   │   │   ├── logger.go
│   │   │   ├── middleware
│   │   │   │   └── middleware.go
│   │   │   ├── sms
│   │   │   │   ├── sms.go
│   │   │   │   └── smsLogic.go
│   │   │   ├── store
│   │   │   │   └── mysql
│   │   │   │       └── store.go
│   │   │   └── wallet
│   │   │       └── grpc
│   │   │           └── grpc.go
│   │   └── utils
│   │       ├── custom_jwt.go
│   │       ├── env.go
│   │       ├── errors.go
│   │       ├── exceptions.go
│   │       ├── filter.go
│   │       ├── httpClient.go
│   │       ├── jwt.go
│   │       ├── otp.go
│   │       ├── password.go
│   │       ├── phone.go
│   │       ├── redis
│   │       │   └── redisManager.go
│   │       ├── request
│   │       │   └── request.go
│   │       ├── scientific.go
│   │       ├── sha256.go
│   │       ├── sms.go
│   │       ├── utils.go
│   │       └── wpool
│   │           ├── exec.go
│   │           └── job.go
│   ├── main.go
│   ├── README.md
│   └── setup.sh
├── README.md
└── services
    ├── cmd
    │   ├── features.go
    │   └── root.go
    ├── docs
    │   ├── docs.go
    │   ├── doctor.go
    │   ├── hospital.go
    │   ├── servicecategory.go
    │   └── service.go
    ├── go.mod
    ├── go.sum
    ├── internal
    │   ├── app
    │   │   ├── app.go
    │   │   ├── http
    │   │   │   └── echo
    │   │   │       └── echo.go
    │   │   └── template
    │   │       └── migration.go
    │   ├── domain
    │   │   ├── account.go
    │   │   ├── AddOns.go
    │   │   ├── const.go
    │   │   ├── doctor.go
    │   │   ├── doctor-hospital.go
    │   │   ├── domain.go
    │   │   ├── grade.go
    │   │   ├── hospital-addons.go
    │   │   ├── hospital.go
    │   │   ├── language.go
    │   │   ├── serviceCategory.go
    │   │   ├── service-doctor.go
    │   │   ├── service.go
    │   │   ├── service-hospital.go
    │   │   ├── SocialMedia.go
    │   │   ├── SocialMediaLink.go
    │   │   └── speciality.go
    │   ├── features
    │   │   ├── addOn
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── doctor
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── doctor-hospital
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── grade
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── hospital
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── hospital-addOn
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── language
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── service
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── serviceCategory
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── service-doctor
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── service-hospital
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── socialMedia
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   ├── socialMediaLink
    │   │   │   ├── handler
    │   │   │   │   └── http
    │   │   │   │       └── handler.go
    │   │   │   ├── repository
    │   │   │   │   └── mysql
    │   │   │   │       └── database.go
    │   │   │   └── usecase
    │   │   │       └── usecase.go
    │   │   └── speciality
    │   │       ├── handler
    │   │       │   └── http
    │   │       │       └── handler.go
    │   │       ├── repository
    │   │       │   └── mysql
    │   │       │       └── database.go
    │   │       └── usecase
    │   │           └── usecase.go
    │   ├── proto
    │   │   ├── account
    │   │   │   ├── account_grpc.pb.go
    │   │   │   ├── account.pb.go
    │   │   │   └── account.proto
    │   │   ├── auth
    │   │   │   ├── auth_grpc.pb.go
    │   │   │   ├── auth.pb.go
    │   │   │   └── auth.proto
    │   │   └── service
    │   │       ├── service_grpc.pb.go
    │   │       ├── service.pb.go
    │   │       └── service.proto
    │   ├── service
    │   │   └── account
    │   │       └── grpc
    │   │           └── grpc.go
    │   └── utils
    │       ├── utils.go
    │       └── wpool
    │           ├── exec.go
    │           └── job.go
    ├── main.go
    └── README.md

```
