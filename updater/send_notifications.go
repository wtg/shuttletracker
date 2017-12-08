package updater

import (		
	"net/smtp"

	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"

)

func GetEmail(phone_number string, input_carrier string) (string) {
	// Map to hold cell carrier names to access the appropriate email address
	var carrier_emails = map[string]string {
		"AT&T": "txt.att.net",
		"Sprint": "messaging.sprintpcs.com",
		"T-Mobile": "tmomail.net",
		"Verison": "vtext.com",

		// TODO: add more potential carriers
	}

	// If the carrier is in the map, return the recipient email address
	carrier, in_map := carrier_emails[input_carrier]

	if in_map{
		return phone_number + "@" + carrier
	} else {
		log.Debugf("Could not find carrier: %s", input_carrier)
		return ""
	}

}

func CreateMessage(current_stop string, next_stop string) ([]byte){
	
	var message_body string = "The shuttle is at " + current_stop + ".\nThe next stop is " + next_stop + "."
	
	msg := []byte("RPI Shuttle Tracker Notification\r\n" +
		message_body + "\r\n")

	return msg
}

func Send(notifications []model.Notification, current_stop string, next_stop string) (int){

	// Get recipient(s) email address, and create message
	var to_emails []string 
	for i := range notifications {
		to_emails = append(to_emails, GetEmail(notifications[i].PhoneNumber, notifications[i].Carrier))
	}

	message := CreateMessage(current_stop, next_stop)

	// Authenticate sender email
	auth := smtp.PlainAuth("", "shuttletrackertest@gmail.com", "shuttletracker2017", "smtp.gmail.com")
	
	// Connect to the server, authenticate, set the sender and recipient, and send
	var sent int = 0
	for i := range to_emails {
		var to = []string{to_emails[i]}
		err := smtp.SendMail("smtp.gmail.com:587", auth, "shuttletrackertest@gmail.com", to, message)
		
		if err != nil {
			log.Debugf("Message send error: %v", err)
		} else {
			log.Debugf("Message sent")
			sent++
		}
	}
	return sent
	
}