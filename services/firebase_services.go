package services

import (
    "context"
    "firebase.google.com/go"
    "firebase.google.com/go/messaging"
    "google.golang.org/api/option"
    "log"
)

// Global variable accessible anywhere in this package
var MessagingClient *messaging.Client

// Initialize Firebase once
func InitFirebase() {
    opt := option.WithCredentialsFile("config/serviceAccountKey.json")
    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("Error initializing Firebase app: %v", err)
    }

    client, err := app.Messaging(context.Background())
    if err != nil {
        log.Fatalf("Error initializing Firebase Messaging client: %v", err)
    }

    MessagingClient = client
    log.Println("✅ Firebase initialized successfully")
}

// Send push notification
func SendNotification(token, title, body string) error {
    msg := &messaging.Message{
        Token: token,
        Notification: &messaging.Notification{
            Title: title,
            Body:  body,
        },
    }

    response, err := MessagingClient.Send(context.Background(), msg)
    if err != nil {
        log.Println("❌ Failed to send notification:", err)
        return err
    }

    log.Println("✅ Notification sent:", response)
    return nil
}
