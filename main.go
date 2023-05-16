package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"go.uber.org/multierr"
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

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:wazzaaa.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatal(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatal(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	ctx := context.Background()
	signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, err := client.GetQRChannel(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if err = client.Connect(); err != nil {
			log.Fatal(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				terminal := qrcodeTerminal.New()
				terminal.Get(evt.Code).Print()
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

	var errs error
	for _, user := range users {
		targetJID := types.NewJID(user, types.DefaultUserServer)

		_, err = client.SendMessage(ctx, targetJID, &waProto.Message{
			Conversation: message,
		})
		if err != nil {
			errs = multierr.Append(errs, err)
		}
	}
	if errs != nil {
		log.Fatalf("Send whatsapp messages: %v", errs)
	}

	client.Disconnect()
}
