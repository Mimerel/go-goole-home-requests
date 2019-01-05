package google_talk

import (
	"fmt"
	"github.com/evalphobia/google-home-client-go/googlehome"
	"time"
)

/**
Given a list of ips, and a message, this method will
loop through the list and run the method to send the message to the different
google homes.
 */
func Talk(ips []string, message string) {
	for _, ip := range ips {
		fmt.Printf("talk message sent to ip : %s \n", ip)
		talkIndividual(ip, message)
	}
}

/**
Method that send a message to the google home for the
message to be read out loud
 */
func talkIndividual(ip string, message string) {
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: ip,
		Lang:     "fr",
		Accent:   "FR",
	})
	if err != nil {
		fmt.Printf("unable to send message\n")
	}
	cli.SetLang("fr")
	cli.Notify(message)
	time.Sleep(3 * time.Second)
}