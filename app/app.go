/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-08-2018
 * |
 * | File Name:     app.go
 * +===============================================
 */

package app

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/I1820/ElRunner/scenario"
	"github.com/I1820/types"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/envy"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

// Application part of the Runner gathers data from MQTT
// system of the I1820 platform and then puts it into its pipeline.
// Pipeline of application consists of following stage
// - Decode Stage
// - Scenario Stage
// - Insert Stage
// Name of the application is used as project name in system MQTT subscription.
type Application struct {
	Name string

	cli paho.Client

	Logger *logrus.Logger

	session *mgo.Client
	db      *mgo.Database
	scr     *scenario.Scenario // scenario instance

	// pipeline channels
	decodeStream   chan *types.State
	scenarioStream chan *types.State
	insertStream   chan *types.State

	// in order to close the pipeline nicely
	decodeCloseChan    chan struct{}  // decode stage sends one value to this channel on its return
	scenarioCloseChan  chan struct{}  // scenario stage sends one value to this channel on its return
	insertCloseCounter sync.WaitGroup // count number of insert stages so `Exit` can wait for all of them
}

// New creates new application. this function does not create mqtt client
// it creates mongodb session and scenario instances
func New(name string) *Application {
	a := Application{}

	a.Name = name

	a.Logger = logrus.New()

	// Create a mongodb connection
	url := envy.Get("DB_URL", "mongodb://127.0.0.1:27017")
	session, err := mgo.NewClient(url)
	if err != nil {
		a.Logger.Fatalf("DB new client error: %s", err)
	}
	a.session = session
	a.scr = scenario.New()

	// pipeline channels
	a.decodeStream = make(chan *types.State)
	a.scenarioStream = make(chan *types.State)
	a.insertStream = make(chan *types.State)

	return &a
}

// Scenario returns application scenario instance.
// This is a bad habit but we require this becuase we update
// scenario from http services.
func (a *Application) Scenario() *scenario.Scenario {
	return a.scr
}

// Run runs application. this function connects mqtt client and then subscribes on its topic
func (a *Application) Run() {
	// create close channels here so we can run and stop single
	// application many times
	a.decodeCloseChan = make(chan struct{}, 1)
	a.scenarioCloseChan = make(chan struct{}, 1)

	// log core application run. This log is useful for tracing project docker
	// state from outside.
	a.Logger.WithFields(logrus.Fields{
		"component": "elrunner",
	}).Infof("ElRunner Link Application on project %s", a.Name)

	// Create an MQTT client
	/*
		Port: 1883
		CleanSession: True
		Order: True
		KeepAlive: 30 (seconds)
		ConnectTimeout: 30 (seconds)
		MaxReconnectInterval 10 (minutes)
		AutoReconnect: True
	*/
	opts := paho.NewClientOptions()
	opts.AddBroker(envy.Get("BROKER_URL", "tcp://127.0.0.1:18083"))
	opts.SetClientID(fmt.Sprintf("I1820-elrunner-%s", a.Name)) // there is only one mqtt client per project
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(client paho.Client) {
		// TODO generic decoder
		if t := a.cli.Subscribe(fmt.Sprintf("i1820/project/%s/raw", a.Name), 0, a.mqttRawHandler); t.Error() != nil {
			a.Logger.Fatalf("MQTT subscribe error: %s", t.Error())
		}

		if t := a.cli.Subscribe(fmt.Sprintf("i1820/projects/%s/things/+/assets/+/state", a.Name), 0, a.mqttStateHandler); t.Error() != nil {
			a.Logger.Fatalf("MQTT subscribe error: %s", t.Error())
		}
	})
	a.cli = paho.NewClient(opts)

	// Connect to the MQTT Server.
	if t := a.cli.Connect(); t.Wait() && t.Error() != nil {
		a.Logger.Fatalf("MQTT session error: %s", t.Error())
	}

	// Connect to the mongodb
	if err := a.session.Connect(context.Background()); err != nil {
		a.Logger.Fatalf("DB connection error: %s", err)
	}
	a.db = a.session.Database("i1820")

	// load the main scenario if it exists
	if err := a.scr.ActivateWithoutCode("main"); err != nil {
	}

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		// go a.decodeStage()
		go a.scenarioStage()
		// go a.insertStage()
		// a.insertCloseCounter.Add(1)
	}
}

// Exit closes mqtt connection then closes all channels and return from all pipeline stages
func (a *Application) Exit() {
	// disconnect waiting time in milliseconds
	var quiesce uint = 10
	a.cli.Disconnect(quiesce)

	// close project stream
	close(a.decodeStream)

	// all channels are going to close
	// so we are waiting for them
	a.insertCloseCounter.Wait()
}
