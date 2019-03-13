package google_talk

import (
	"github.com/evalphobia/google-home-client-go/googlehome"
	"go-goole-home-requests/configuration"
	"time"
)

/**
Given a list of ips, and a message, this method will
loop through the list and run the method to send the message to the different
google homes.
 */
func Talk(config *configuration.Configuration, ips []string, message string) {
	for _, ip := range ips {
		config.Logger.Info("talk message sent to ip : %s \n", ip)
		talkIndividual(config, ip, message)
	}
}

/**
Method that send a message to the google home for the
message to be read out loud
 */
func talkIndividual(config *configuration.Configuration, ip string, message string) {
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: ip,
		Lang:     "fr",
		Accent:   "FR",
	})
	if err != nil {
		config.Logger.Error("unable to send message\n")
	}
	cli.SetLang("fr")
	cli.Notify(message)
	time.Sleep(3 * time.Second)
}