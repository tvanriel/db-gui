package mysqldump

type DumpCommunicator interface {
	Start(amountOfTablesToDump int, amountOfViewsToDump int)
	DumpingTable(tableName string, number int)
	DumpingView(viewName string, number int)
	Complete(filename string)
}

var _ DumpCommunicator = &SummaryCommunicator{}

type SummaryCommunicator struct {
	amountOfTablesToDump int
	amountOfViewsToDump  int
	tablesDumped         []string
	viewsDumped          []string
	filenameComplete     string
}

func (sc *SummaryCommunicator) Start(amountOfTablesToDump, amountOfViewsToDump int) {
	sc.tablesDumped = make([]string, amountOfTablesToDump)
	sc.viewsDumped = make([]string, amountOfViewsToDump)
	sc.amountOfTablesToDump = amountOfTablesToDump
	sc.amountOfViewsToDump = amountOfViewsToDump
}

func (sc *SummaryCommunicator) DumpingTable(tableName string, number int) {
	sc.tablesDumped[number] = tableName
}
func (sc *SummaryCommunicator) DumpingView(viewName string, number int) {
	sc.viewsDumped[number] = viewName
}
func (sc *SummaryCommunicator) Complete(filename string) {
	sc.filenameComplete = filename
}

func (sc *SummaryCommunicator) GetDumpFilename() string {
	return sc.filenameComplete
}
