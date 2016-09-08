package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/loadimpact/speedboat/api"
	"github.com/loadimpact/speedboat/lib"
	"gopkg.in/guregu/null.v3"
	"gopkg.in/urfave/cli.v1"
	"os"
	"strconv"
)

var commandStatus = cli.Command{
	Name:      "status",
	Usage:     "Looks up the status of a running test",
	ArgsUsage: " ",
	Action:    actionStatus,
}

var commandScale = cli.Command{
	Name:      "scale",
	Usage:     "Scales a running test",
	ArgsUsage: "vus",
	Action:    actionScale,
}

var commandPause = cli.Command{
	Name:      "pause",
	Usage:     "Pauses a running test",
	ArgsUsage: " ",
	Action:    actionPause,
}

var commandResume = cli.Command{
	Name:      "resume",
	Usage:     "Resumes a paused test",
	ArgsUsage: " ",
	Action:    actionResume,
}

func dumpYAML(v interface{}) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		log.WithError(err).Error("Serialization Error")
		return err
	}
	_, _ = os.Stdout.Write(bytes)
	return nil
}

func actionStatus(cc *cli.Context) error {
	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil {
		log.WithError(err).Error("Couldn't create a client")
		return err
	}

	status, err := client.Status()
	if err != nil {
		log.WithError(err).Error("Error")
		return err
	}
	return dumpYAML(status)
}

func actionScale(cc *cli.Context) error {
	args := cc.Args()
	if len(args) != 1 {
		return cli.NewExitError("Wrong number of arguments!", 1)
	}
	vus, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		log.WithError(err).Error("Error")
		return err
	}

	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil {
		log.WithError(err).Error("Couldn't create a client")
		return err
	}

	status, err := client.UpdateStatus(lib.Status{VUs: null.IntFrom(vus)})
	if err != nil {
		log.WithError(err).Error("Error")
		return err
	}
	return dumpYAML(status)
}

func actionPause(cc *cli.Context) error {
	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil {
		log.WithError(err).Error("Couldn't create a client")
		return err
	}

	status, err := client.UpdateStatus(lib.Status{Running: null.BoolFrom(false)})
	if err != nil {
		log.WithError(err).Error("Error")
		return err
	}
	return dumpYAML(status)
}

func actionResume(cc *cli.Context) error {
	client, err := api.NewClient(cc.GlobalString("address"))
	if err != nil {
		log.WithError(err).Error("Couldn't create a client")
		return err
	}

	status, err := client.UpdateStatus(lib.Status{Running: null.BoolFrom(true)})
	if err != nil {
		log.WithError(err).Error("Error")
		return err
	}
	return dumpYAML(status)
}