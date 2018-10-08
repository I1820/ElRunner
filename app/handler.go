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

// mqttRawHandler handles raw data that is coming from MQTT. Raw data needs parse
// before anything so we pass it into decode stage.
func (a *Application) mqttRawHandler(client paho.Client, message paho.Message) {
	/* TODO correct generic decoders
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
	*/
}

// mqttStateHandler handles states that are coming from MQTT. These states do not need
// parse so they are passed directly into scenario stage.
func (a *Application) mqttStateHandler(client paho.Client, message paho.Message) {
	var d types.State
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
	a.scenarioStream <- &d
}
