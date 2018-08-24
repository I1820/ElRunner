/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-08-2018
 * |
 * | File Name:     handler.go
 * +===============================================
 */

package app

import (
	"encoding/json"

	"github.com/I1820/types"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func (a *Application) mqttRawHandler(client paho.Client, message paho.Message) {
	var d types.Data
	if err := json.Unmarshal(message.Payload(), &d); err != nil {
		a.Logger.WithFields(logrus.Fields{
			"component": "elrunner",
			"topic":     message.Topic(),
		}).Errorf("Raw Marshal error %s", err)
		return
	}
	a.Logger.WithFields(logrus.Fields{
		"component": "elrunner",
	}).Infof("Raw Marshal on %v", d)
	a.decodeStream <- d
}

func (a *Application) mqttDataHandler(client paho.Client, message paho.Message) {
	var d types.Data
	if err := json.Unmarshal(message.Payload(), &d); err != nil {
		a.Logger.WithFields(logrus.Fields{
			"component": "elrunner",
			"topic":     message.Topic(),
		}).Errorf("Data Marshal error %s", err)
		return
	}
	a.Logger.WithFields(logrus.Fields{
		"component": "elrunner",
	}).Infof("Data Marshal on %v", d)
	a.scenarioStream <- d
}
