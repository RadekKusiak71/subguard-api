package tasks

import (
	"bytes"
	"log"
	"net/smtp"
	"text/template"
	"time"

	"github.com/RadekKusiak71/subguard-api/internal/config"
	"github.com/RadekKusiak71/subguard-api/internal/subscriptions"
	"github.com/RadekKusiak71/subguard-api/internal/users"
)

type SubscriptionCron struct {
	subStore  subscriptions.SubscriptionStore
	userStore users.UserStore
}

type EmailData struct {
	Subscriptions []subscriptions.Subscription
	Year          int
}

func NewSubscriptionCron(subStore subscriptions.SubscriptionStore, userStore users.UserStore) *SubscriptionCron {
	return &SubscriptionCron{subStore: subStore, userStore: userStore}
}

func (sc *SubscriptionCron) StartReminderCron() func() {
	return func() {
		now := time.Now().Format(time.RFC3339)
		log.Printf("[%s] Running subscription reminder cron...", now)

		subs, err := sc.subStore.GetExpiringSoon()
		if err != nil {
			log.Printf("[%s] Error fetching expiring subscriptions: %v", now, err)
			return
		}

		if len(subs) == 0 {
			log.Printf("[%s] No subscriptions expiring tomorrow.", now)
			return
		}

		usersSubs := make(map[int][]subscriptions.Subscription)
		for _, sub := range subs {
			usersSubs[sub.UserID] = append(usersSubs[sub.UserID], sub)
		}

		for userID, userSubs := range usersSubs {
			if err := sc.sendMail(userID, userSubs); err != nil {
				log.Printf("[ERROR] Failed to send email to user %d: %v", userID, err)
			}

			if err := sc.subStore.UpdateNextPaymentBatch(userSubs); err != nil {
				log.Printf("[ERROR] Failed to update next_payment_at for user %d: %v", userID, err)
			}
		}
	}
}

func (sc *SubscriptionCron) sendMail(userID int, expiringSubs []subscriptions.Subscription) error {
	user, err := sc.userStore.Get(userID)
	if err != nil {
		log.Printf("[ERROR] Couldn't find user with ID: %d", userID)
		return err
	}

	tmpl, err := template.ParseFiles("internal/templates/email_template.html")
	if err != nil {
		log.Printf("[ERROR] Failed to parse email template: %v", err)
		return err
	}

	data := EmailData{
		Subscriptions: expiringSubs,
		Year:          time.Now().Year(),
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("[ERROR] Failed to execute email template: %v", err)
		return err
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	sender := config.Config.EmailSender
	password := config.Config.EmailPassword

	headers := map[string]string{
		"From":         sender,
		"To":           user.Email,
		"Subject":      "Reminder: Subscriptions Expiring Tomorrow",
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(k + ": " + v + "\r\n")
	}
	msg.WriteString("\r\n" + body.String())

	auth := smtp.PlainAuth("", sender, password, smtpHost)
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{user.Email}, msg.Bytes()); err != nil {
		log.Printf("[ERROR] Failed to send email to %s: %v", user.Email, err)
		return err
	}

	log.Printf("[INFO] Email successfully sent to %s", user.Email)
	return nil
}
