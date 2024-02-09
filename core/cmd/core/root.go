/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	events "main.go/proto/core/events"
	core "main.go/proto/core"
	logger "main.go/services/logger"
	crypto "main.go/services/crypto"
	"main.go/services"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/mackerelio/go-osstat/memory"
	"github.com/natefinch/lumberjack"
	"runtime/debug"
	"time"
	"context"
	"net/http"
	"fmt"
	"os"
	"log"
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

var cfgFile string
var globalUsage = `QIS System`
var amqpAddress = "amqp://user:password@localhost:5672/"
var host = viper.GetString("host")
var environment = viper.GetString("environment")
var grpcuiHost = viper.GetString("grpcui_host")
var grpcwebHost = viper.GetString("grpcweb_host")
var webserverHost = viper.GetString("webserver_host")
var useTLS = viper.GetBool("use_tls")
var tlsCertPath = viper.GetString("tls_cert_path")
var tlsKeyPath = viper.GetString("tls_key_path")
var redisAddrs = viper.GetString("redis_settings.redis_addrs")
var redisPassword = viper.GetString("redis_settings.redis_password")
var datastore_dsn = viper.GetString("datastore.core_dsn")
var datastore_username = viper.GetString("datastore.core_username")
var datastore_password = viper.GetString("datastore.core_password")
var datastore_database = viper.GetString("datastore.core_database")
var profile_server_host = viper.GetString("profile_host")
var enable_profile_server = viper.GetBool("debug.enable_profile_server")

var settings = &core.Settings{
	ServeAddress: host,
	Environment:  getEnvironment(environment),
	//SourceType:   proto.SettingsSourceType_SETTINGS_SOURCE_BINARY,
	SourceUrl:   "",
	UseTls:      useTLS,
	TlsCertPath: tlsCertPath,
	TlsKeyPath:  tlsKeyPath,
	CustomerApiSettings: &core.CustomerAPISettings{
		DefaultUrl: "0.0.0.0:8080",
		MongoConfig: &core.MongoClientConfig{
			Url:            datastore_dsn,
			DatabaseName:   datastore_database,
			CollectionName: "Customer",
			ContextTimeout: 120,
		},
	},
	RedisSettings: &core.RedisSettings{
		Addrs:    redisAddrs,
		Password: redisPassword,
	},
	DatastoreConfig: &core.DataStoreConfig{
		Dsn:      datastore_dsn,
		Username: datastore_username,
		Password: datastore_password,
		Database: datastore_database,
	},
	DebugServiceSettings: &core.DebugServiceSettings{
		ServeAddress:        profile_server_host,
		EnableProfileServer: enable_profile_server,
	},
}
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "app_config",
	Short: "QIS System",
	Long:  globalUsage,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Database configs
		cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			viper.BindPFlag(f.Name, f)
		})

		return nil
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logger := watermill.NewStdLogger(false, false)

		// CQRS is built on messages router. Detailed documentation: https://watermill.io/docs/messages-router/
		router, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			panic(err)
		}
	
		// Simple middleware which will recover panics from event or command handlers.
		// More about router middlewares you can find in the documentation:
		// https://watermill.io/docs/messages-router/#middleware
		//
		// List of available middlewares you can find in message/router/middleware.
		router.AddMiddleware(middleware.Recoverer)
	
		// cqrs.Facade is facade for Command and Event buses and processors.
		// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
		cqrsFacade, err := services.NewServer(amqpAddress)
		if err != nil {
			panic(err)
		}
	
		// publish BookRoom commands every second to simulate incoming traffic
		go publishCommands(cqrsFacade.CommandBus())
	
		// processors are based on router, so they will work when router will start
		if err := router.Run(context.Background()); err != nil {
			panic(err)
		}
	},
}

func getEnvironment(environment string) core.Environment {
	switch strings.ToLower(environment) {
	case "dev":
		return core.Environment_DEVELOPMENT
	case "prd":
		return core.Environment_PRODUCTION
	default:
		return core.Environment_DEVELOPMENT
	}
}

func allowCors(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "*")
	// resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// resp.Header().Set("Access-Control-Expose-Headers", "grpc-status, grpc-message")
	//resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, XMLHttpRequest, X-User-Agent, X-Grpc-Web, Grpc-Timeout, Grpc-Status, Grpc-Message")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.app_config.yaml)")
	rootCmd.PersistentFlags().String("host", "0.0.0.0:8006", "main app host (addr and port)")
	rootCmd.PersistentFlags().String("grpcui_host", "0.0.0.0:8007", "grpcui app host (addr and port)")
	rootCmd.PersistentFlags().String("grpcweb_host", "0.0.0.0:8008", "grpcweb host (addr and port)")
}

func InitLogger() {
	writerSyncer := getLogWriter()
	encoder := getEncoder()
	logLevel := getLogLevel(viper.GetString("log.level"))
	level := viper.GetString("log.level")
	if level == "info" {
		logLevel = zapcore.InfoLevel
	}
	if !viper.GetBool("log.disabled") {
		logger.InitZapCores(encoder, logLevel, writerSyncer, zapcore.Lock(os.Stdout))
	} else {
		logger.InitPlain()
	}
	defer writerSyncer.Sync()
	if logLevel == zapcore.DebugLevel {
		// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
		grpc_zap.ReplaceGrpcLoggerV2(logger.GetSugarLogger(context.Background()))
	}
}

func InitMemoryAllocation() {
	memory, err := memory.Get()
	if err != nil {
		logger.Panicf(context.TODO(), "Error getting memory {%v}", os.Stderr)
	}
	logger.Infof(context.TODO(), "Total Memory: %d Megabytes\n", (memory.Total / 1e+6))
	debug.SetGCPercent(30)
	// total := int64(memory.Total / 10)
	// logger.Infof(context.TODO(), "Total Memory Limit: %d Megabytes\n", (total / 1e+6))
	// debug.SetMemoryLimit(total)
}

func getLogWriter() zapcore.WriteSyncer {
	maxSize := viper.GetInt("log.file_max_size")
	if maxSize == 0 {
		maxSize = 256 // Default Value
	}
	maxAge := viper.GetInt("log.rotation_days")
	if maxAge == 0 {
		maxAge = 30 // Default Value
	}
	compressFile := viper.GetBool("log.compress_old_file")
	logApplication := getWriteSyncer("Application.log")
	lumberJackLogger := &lumberjack.Logger{
		Filename: logApplication,
		MaxSize:  maxSize,
		MaxAge:   maxAge,
		Compress: compressFile,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getWriteSyncer(logName string) string {
	logsDirectory := viper.GetString("log.directory")
	if logsDirectory == "" {
		log.Fatalf("Unable to initialize logging infrastructure")
	}
	hostName, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to initialize logging infrastructure, invalid hostname %s", hostName)
	}
	logApplication := fmt.Sprintf("%s%s/%s", logsDirectory, hostName, logName)
	return logApplication
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogLevel(logLevel string) zapcore.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".treasure" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".treasure")
	}

	viper.SetEnvPrefix("treasure")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	InitLogger()
	InitMemoryAllocation()
	if viper.GetBool("init") {
		crypto.GenesisQRCode()
		crypto.NewJWTRandomKey()
		crypto.NewAESRandomKey()
		os.Exit(1)
	}
}


func publishCommands(commandBus *cqrs.CommandBus) func() {
	i := 0
	for {
		i++

		startDate, err := ptypes.TimestampProto(time.Now())
		if err != nil {
			panic(err)
		}

		endDate, err := ptypes.TimestampProto(time.Now().Add(time.Hour * 24 * 3))
		if err != nil {
			panic(err)
		}

		bookRoomCmd := &events.BookRoom{
			RoomId:    fmt.Sprintf("%d", i),
			GuestName: "John",
			StartDate: startDate,
			EndDate:   endDate,
		}
		if err := commandBus.Send(context.Background(), bookRoomCmd); err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}
}