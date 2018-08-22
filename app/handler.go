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

func (a *Application) mqttHandler(client paho.Client, message paho.Message) {
	var d types.Data
	if err := json.Unmarshal(message.Payload(), &d); err != nil {
		a.Logger.WithFields(logrus.Fields{
			"component": "elrunner",
			"topic":     message.Topic(),
		}).Errorf("Marshal error %s", err)
		return
	}
	a.Logger.WithFields(logrus.Fields{
		"component": "elrunner",
	}).Infof("Marshal on %v", d)
	a.decodeStream <- d
}
