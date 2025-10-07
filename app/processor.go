package main

type Processor struct {
}

func defineResponse(strings []string) string {
	var response string
	if len(strings) > 0 {
		if strings[0] == "PING" {
			response = "+PONG\r\n"
		} else if strings[0] == "ECHO" && len(strings) > 1 {
			response = "+" + strings[1:][0]
			for _, s := range strings[2:] {
				response += " " + s
			}
			response += "\r\n"
		} else {
			response = "+PONG\r\n"
		}
	}
	return response
}
