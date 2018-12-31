package google_talk

import (
	"fmt"
	"github.com/evalphobia/google-home-client-go/googlehome"
	"time"
)

func Talk(ips []string, message string) {
	for _, ip := range ips {
		fmt.Printf("talk message sent to ip : %s", ip)
		talkIndividual(ip, message)
	}
}

func talkIndividual(ip string, message string) {
	cli, err := googlehome.NewClientWithConfig(googlehome.Config{
		Hostname: ip,
		Lang:     "fr",
		Accent:   "FR",
	})
	if err != nil {
		fmt.Printf("unable to send message")
	}
	cli.SetLang("fr")
	cli.Notify(message)
	time.Sleep(3 * time.Second)
}