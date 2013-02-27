package main

type UI interface {
	GetInput() string
	ProcessOutput(output string)
	ProcessAlert(output string)
}



