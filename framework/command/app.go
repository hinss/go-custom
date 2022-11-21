package command

import (
	"context"
	"database/sql"
	"fmt"
	mq "gitee.com/ChengHoHins/component-base/pkg/kafka"
	"github.com/hinss/go-custom/framework/cobra"
	"github.com/hinss/go-custom/framework/contract"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var appAddress = ""

// initAppCommand 初始化app命令和其子命令
func initAppCommand() *cobra.Command {
	// appStartCommand 添加flag
	appStartCommand.Flags().StringVar(&appAddress, "address", ":8888", "设置app启动的地址，默认为8888")
	appCommand.AddCommand(appStartCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，它没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，其包含业务启动，关闭，重启，查询等功能",
	RunE: func(c *cobra.Command, args []string) error {
		// 打印帮助文档
		c.Help()
		return nil
	},
}

// appStartCommand 启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个Web服务 并且初始化所有配置项",
	RunE: func(c *cobra.Command, args []string) error {
		// 从Command中获取服务容器
		container := c.GetContainer()
		// 从服务容器中获取kernel的服务实例
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		// 从kernel服务实例中获取引擎
		core := kernelService.HttpEngine()

		// 获取配置服务实例
		configService := container.MustMake(contract.ConfigKey).(contract.Config)

		appReadTimeout := configService.GetInt("app.read_timeout")
		appWriteTimeout := configService.GetInt("app.write_timeout")
		kafkaHosts := configService.GetStringSlice("kafka.hosts")
		kafkaEnable := configService.GetBool("kafka.enable")
		mysqlEnable := configService.GetBool("mysql.enable")
		mysqlMap := configService.GetStringMapString("mysql")

		// 创建一个Server服务
		server := &http.Server{
			Addr:           appAddress,
			Handler:        core,
			ReadTimeout:    time.Duration(int64(appReadTimeout) * int64(time.Second)),
			WriteTimeout:   time.Duration(int64(appWriteTimeout) * int64(time.Second)),
			MaxHeaderBytes: 1 << 20,
		}

		// 这个goroutine是启动服务的goroutine
		go func() {
			server.ListenAndServe()
		}()

		// 初始化kafka
		if kafkaEnable {
			initKafkaProducer(kafkaHosts)
		}
		
		// 初始化mysql
		if mysqlEnable {
			initMysql(mysqlMap)
		}

		// 当前的goroutine等待信号量
		quit := make(chan os.Signal)
		// 监控信号：SIGINT, SIGTERM, SIGQUIT
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 这里会阻塞当前goroutine等待信号
		<-quit

		// 调用Server.Shutdown graceful结束
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		return nil
	},
}

func initKafkaProducer(hosts []string) {
	// 初始化kafka
	err := mq.InitSyncProducer(mq.DefaultSyncProducer,
		hosts, nil)
	if err != nil {
		//danger(fmt.Sprintf("InitSyncKafkaProducer err: %s client: %s", err.Error(), mq.DefaultSyncProducer))
		panic(err)
	}
}

var Db *sql.DB

func initMysql(mysqlMap map[string]string) {

	mysqlHost := mysqlMap["host"]
	mysqlPwd := mysqlMap["password"]
	mysqlUser := mysqlMap["user"]
	mysqlDb := mysqlMap["database"]
	mysqlPort := mysqlMap["port"]

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", mysqlUser, mysqlPwd, mysqlHost, mysqlPort, mysqlDb)

	var err error
	Db, err = sql.Open("mysql", dataSourceName)
	//Db, err = sql.Open("mysql", "root:Bluesea#tiger?123@tcp(120.76.97.16:3306)/margin?parseTime=True")
	if err != nil {
		log.Fatal(err)
	}
	return
}

