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

	"github.com/I1820/ElRunner/scenario"
	"github.com/I1820/types"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gobuffalo/envy"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

// Application is a main component of uplink that consists of
// uplink protocols and mqtt client
type Application struct {
	Name string

	cli paho.Client

	Logger *logrus.Logger

	session *mgo.Client
	db      *mgo.Database
	scr     *scenario.Scenario

	// pipeline channels
	decodeStream   chan types.Data
	scenarioStream chan types.Data
	insertStream   chan types.Data
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
	a.decodeStream = make(chan types.Data)
	a.scenarioStream = make(chan types.Data)
	a.insertStream = make(chan types.Data)

	return &a
}

// Scenario returns application scenario instance
func (a *Application) Scenario() *scenario.Scenario {
	return a.scr
}

// Run runs application. this function connects mqtt client and then register its topic
func (a *Application) Run() {
	a.Logger.WithFields(logrus.Fields{
		"component": "elrunner",
	}).Infof("ElRunner Link Application %s", a.Name)

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
	opts.AddBroker(envy.Get("BROKER_URL", "tcp://127.0.0.1:1883"))
	opts.SetClientID(fmt.Sprintf("I1820-elrunner-%s", a.Name))
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(func(client paho.Client) {
		if t := a.cli.Subscribe(fmt.Sprintf("i1820/project/%s/raw", a.Name), 0, a.mqttRawHandler); t.Error() != nil {
			a.Logger.Fatalf("MQTT subscribe error: %s", t.Error())
		}
		if t := a.cli.Subscribe(fmt.Sprintf("i1820/project/%s/data", a.Name), 0, a.mqttDataHandler); t.Error() != nil {
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

	// scenario
	if err := a.scr.ActivateWithoutCode("main"); err != nil {
	}

	// pipeline stages
	for i := 0; i < runtime.NumCPU(); i++ {
		go a.decode()
		go a.scenario()
		go a.insert()
	}
}
