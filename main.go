package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	qrterminal "github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/multierr"
	_ "modernc.org/sqlite"
)

func main() {
	recipients := flag.String("recipients", "", "list of comma separated phone numbers")
	message := flag.String("message", "", "message to send")
	flag.Parse()

	if recipients == nil || *recipients == "" {
		log.Fatal("recipients must be specified")
	}

	users := strings.Split(strings.TrimSpace(*recipients), ",")
	if len(users) == 0 {
		log.Fatal("invalid recipients")
	}

	if message == nil || *message == "" {
		log.Fatal("message must be specified")
	}

	ctx := context.Background()

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New(ctx, "sqlite", "file:wazzaaa.db?_pragma=foreign_keys(1)", dbLog)
	if err != nil {
		log.Fatal(err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, err := client.GetQRChannel(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		if err = client.Connect(); err != nil {
			log.Fatal(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				log.Println("QR code:", evt.Code)
			} else {
				log.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		if err = client.Connect(); err != nil {
			log.Fatal(err)
		}
	}
	defer client.Disconnect()

	var errs error
	for _, user := range users {
		targetJID := types.NewJID(user, types.DefaultUserServer)

		_, err = client.SendMessage(ctx, targetJID, &waE2E.Message{
			Conversation: message,
		})
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	if errs != nil {
		log.Fatalf("Send whatsapp messages: %v", errs)
	}
}
