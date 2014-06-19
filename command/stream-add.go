package command

import (
	"fmt"
	"lsf"
	"lsf/panics"
	"lsf/schema"
	"lsf/system"
)

const addStreamCmdCode lsf.CommandCode = "stream-add"

var addStream *lsf.Command
var addStreamOptions *editStreamOptionsSpec

func init() {

	addStream = &lsf.Command{
		Name:  addStreamCmdCode,
		About: "Add a new log stream",
		Init:  _verifyEditStreamRequiredOpts,
		Run:   runAddStream,
		Flag:  FlagSet(addStreamCmdCode),
	}
	addStreamOptions = initEditStreamOptionsSpec(addStream.Flag)
}

func runAddStream(env *lsf.Environment, args ...string) (err error) {
	defer panics.Recover(&err)

	id := schema.StreamId(*addStreamOptions.id.value)
	pattern := *addStreamOptions.pattern.value
	mode := schema.ToJournalModel(*addStreamOptions.mode.value)
	path := *addStreamOptions.path.value
	fields := make(map[string]string) // TODO: fields needs a solution

	// check existing
	docid := system.DocId(fmt.Sprintf("stream.%s.stream", id))
	_assertNotExists(env, docid)

	// lock lsf port's "streams" resource to prevent race condition
	lock := _lockResource(env, "streams", "add stream")
	defer lock.Unlock()

	// create the stream-conf file.
	logstream := schema.NewLogStream(id, path, mode, pattern, fields)

	e := env.CreateDocument(docid, logstream)
	panics.OnError(e, "command.runAddStream:", "CreateDocument:", id)

	return nil
}