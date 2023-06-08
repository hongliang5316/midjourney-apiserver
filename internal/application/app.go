package application

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/hongliang5316/midjourney-apiserver/pkg/api"
	"github.com/hongliang5316/midjourney-go/midjourney"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

type Application struct {
	*discordgo.Session
	Cli *midjourney.Client
	Cfg *Config
}

func New() *Application {
	cfg := new(Config)

	data, err := ioutil.ReadFile("./conf/conf.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal([]byte(data), cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.ListenPort == 0 {
		cfg.ListenPort = 8080
	}

	dg, err := discordgo.New(cfg.UserToken)
	if err != nil {
		log.Fatal(err)
	}

	cli := midjourney.NewClient(&midjourney.Config{
		UserToken: cfg.UserToken,
	})

	app := &Application{dg, cli, cfg}

	dg.AddHandler(app.messageCreate)
	dg.AddHandler(app.messageUpdate)

	dg.Identify.Intents = discordgo.IntentsAll

	return app
}

func (app *Application) Run() error {
	go func(app *Application) {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", app.Cfg.ListenPort))
		if err != nil {
			log.Fatalf("failed to listen: %+v", err)
		}

		s := grpc.NewServer()
		api.RegisterAPIServiceServer(s, new(Service))

		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}(app)

	err := app.Open()
	if err != nil {
		return fmt.Errorf("Call app.Open failed, err: %w", err)
	}

	log.Printf("Start...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return app.Close()
}
