package main

import (
    "fmt"
    "log"

    "path/to/your/package/wabot" // Replace with the actual import path
)

func main() {
    clientID := "YOUR_CLIENT_ID"
    clientSecret := "YOUR_CLIENT_SECRET"

    wabotClient := wabot.NewWabotApiClient(clientID, clientSecret)

    // Authenticate
    err := wabotClient.Authenticate()
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }

    fmt.Println("Authenticated successfully.")

    // Get Templates
    templates, err := wabotClient.GetTemplates()
    if err != nil {
        log.Fatalf("Failed to get templates: %v", err)
    }

    for _, template := range templates {
        templateID := template["template_id"]
        templateName := template["name"]
        fmt.Printf("Template ID: %v, Name: %v\n", templateID, templateName)
    }

    // Send a message
    to := "+1234567890"           // Replace with the recipient's phone number
    templateID := "339"           // Replace with your template ID
    templateParams := []string{"John", "your email address"}

    err = wabotClient.SendMessage(to, templateID, templateParams)
    if err != nil {
        log.Fatalf("Failed to send message: %v", err)
    }

    fmt.Println("Message sent successfully.")

    // Logout
    err = wabotClient.Logout()
    if err != nil {
        log.Fatalf("Logout failed: %v", err)
    }

    fmt.Println("Logged out successfully.")
}
