package main

type Processor struct {
	storage map[string]string
}

func NewProcessor() *Processor {
	return &Processor{
		storage: make(map[string]string),
	}
}

func (p *Processor) ProcessCommand(strings []string) string {
	var response string
	response = ""
	if len(strings) == 0 {
		response = ""
	} else if strings[0] == "PING" {
		response = "+PONG\r\n"
	} else if strings[0] == "ECHO" {
		response = p.commandEcho(strings)
	} else {
		response = "+PONG\r\n"
	}
	return response
}

func (p *Processor) commandEcho(strings []string) string {
	var response string
	response = "+"
	if len(strings) > 1 {
		response += strings[1:][0]
		for _, s := range strings[2:] {
			response += " " + s
		}
	}
	response += "\r\n"
	return response
}
