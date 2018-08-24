/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 02-08-2018
 * |
 * | File Name:     pipeline.go
 * +===============================================
 */

package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/I1820/ElRunner/codec"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/sirupsen/logrus"
)

func (a *Application) decode() {
	for d := range a.decodeStream {
		decoder, err := codec.NewWithoutCode(d.ThingID)
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "elrunner",
			}).Errorf("%s does not exist on GoRunner", d.ThingID)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		parsed, err := decoder.Decode(ctx, d.Raw)
		cancel()
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "elrunner",
			}).Errorf("Decode error: %s", err)
			continue
		} else {
			d.Data = parsed
			a.Logger.WithFields(logrus.Fields{
				"component": "elrunner",
			}).Infof("Decode on: %v", d)
		}

		// Publish parsed data
		b, err := json.Marshal(d)
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "elrunner",
			}).Errorf("Marshal data error: %s", err)
		}
		a.cli.Publish(fmt.Sprintf("i1820/project/%s/data", d.Project), 0, false, b)
		a.Logger.WithFields(logrus.Fields{
			"component": "elrunner",
		}).Infof("Publish parsed data: %s", d.Project)

		a.insertStream <- d
	}
}

func (a *Application) scenario() {
	for d := range a.scenarioStream {
		b, err := json.Marshal(d)
		if err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "elrunner",
			}).Errorf("Marshal data error: %s", err)
		}
		a.scr.Data(string(b), d.ThingID)
	}
}

func (a *Application) insert() {
	for d := range a.insertStream {
		if _, err := a.db.Collection("data").ReplaceOne(context.Background(), bson.NewDocument(
			bson.EC.String("thingid", d.ThingID),
			bson.EC.Time("timestamp", d.Timestamp),
			bson.EC.Binary("raw", d.Raw),
		), d); err != nil {
			a.Logger.WithFields(logrus.Fields{
				"component": "elrunner",
			}).Errorf("Mongo Replace: %s", err)
		} else {
			a.Logger.WithFields(logrus.Fields{
				"component": "elrunner",
			}).Infof("Insert into database: %v", d)
		}
	}
}
