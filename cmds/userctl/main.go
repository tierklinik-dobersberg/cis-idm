package main

import "github.com/sirupsen/logrus"

func main() {
	if err := getRootCommand().Execute(); err != nil {
		logrus.Fatal(err.Error())
	}
}
