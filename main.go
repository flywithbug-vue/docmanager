package main

import (
	"doc-manager/web_server/config"
	"doc-manager/web_server/core/jwt"
	"doc-manager/web_server/core/mongo"
	"doc-manager/web_server/model"
	"doc-manager/web_server/server"
	"flag"

	log "github.com/flywithbug/log4go"
)

//log 启动配置
func SetLog() {
	w := log.NewFileWriter()
	w.SetPathPattern("./log/log-%Y%M%D.log")
	c := log.NewConsoleWriter()
	c.SetColor(true)
	log.Register(w)
	log.Register(c)
	log.SetLevel(1)
	log.SetLayout("2006-01-02 15:04:05")
}

func main() {
	//配置文件
	configPath := flag.String("config", "config.json", "Configuration file to use")
	flag.Parse()
	err := config.ReadConfig(*configPath)
	if err != nil {
		panic(err)
	}
	conf := config.Conf()

	//signingKey read
	jwt.ReadSigningKey(conf.PrivateKeyPath, conf.PublicKeyPath)

	SetLog()
	defer log.Close()
	//mongodb启动连接
	//设置数据库名字
	model.SetDBName(conf.DBConfig.DBName)
	mongo.DialMgo(conf.DBConfig.Url)
	//go func() {
	//	//静态文件服务
	//	server.StartWeb(conf.WebPort, conf.StaticPath)
	//}()
	//启动ApiServer服务
	server.StartServer(conf.ApiPort,conf.StaticPath, conf.RouterPrefix, conf.AuthPrefix)
}
