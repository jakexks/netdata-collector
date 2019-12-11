package plugin

const netdataType = "go_micro_services"
const netdataModule = "micro"

type chart struct {
	needsUpdate bool
	id          string
	name        string
	title       string
	units       string
	family      string
	context     string
	chartType   string
	priority    string
	updateEvery string
	options     string
	plugin      string
	module      string
}

type dimension struct {
}
