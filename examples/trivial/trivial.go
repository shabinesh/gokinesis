/*
 * A trivial example consuming from Kinesis using the KCL daemon.
 *
 * Just echos messages it sees to stderr.
 */
package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/nieksand/gokinesis/src/kinesis"
)

type EchoConsumer struct {
	outfile *os.File
}

func (ec *EchoConsumer) Init(shardId string) error {

	var err error
	ec.outfile, err = os.OpenFile(fmt.Sprintf("/tmp/%s.demo", shardId), os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	fmt.Fprintf(ec.outfile, "init: %s\n", shardId)
	return nil
}

func (ec *EchoConsumer) ProcessRecords(records []*kinesis.KclRecord, checkpointer *kinesis.Checkpointer) error {
	for _, rec := range records {
		data, err := base64.StdEncoding.DecodeString(rec.DataB64)
		if err != nil {
			fmt.Fprintf(ec.outfile, "process error: %s\n", err.Error())
		}
		fmt.Fprintf(ec.outfile, "process: %s\n", data)
	}

	// Abort execution on checkpointing errors.  We could retry here instead if
	// we wanted.
	return checkpointer.CheckpointAll()
}

func (ec *EchoConsumer) Shutdown(shutdownType kinesis.ShutdownType,
	checkpointer *kinesis.Checkpointer) error {

	fmt.Fprintf(ec.outfile, "shutdown: %s\n", shutdownType)
	if shutdownType == kinesis.GracefulShutdown {
		if err := checkpointer.CheckpointAll(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var ec EchoConsumer
	kinesis.Run(&ec)
}
