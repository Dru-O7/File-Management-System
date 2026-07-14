package email

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"

	"office-file-sharing/backend/internal/shared/config"
	"office-file-sharing/backend/internal/shared/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SendNotificationEmail fetches the recipient, parses the notification payload,
// renders the appropriate email template, and dispatches the email using SMTP.
func SendNotificationEmail(db *gorm.DB, notifID uuid.UUID) {
	// We read the latest notification record from the DB
	var notif models.Notification
	if err := db.Preload("Recipient").First(&notif, "id = ?", notifID).Error; err != nil {
		log.Printf("[Email Service] Failed to fetch notification %s: %v", notifID, err)
		return
	}

	cfg := config.Load()
	if cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		log.Println("[Email Service] SMTP host or user is not configured. Skipping email sending.")
		return
	}

	if notif.Recipient.Email == "" {
		log.Printf("[Email Service] Recipient user %s (%s) has no email address. Skipping.", notif.Recipient.Name, notif.Recipient.ID)
		return
	}

	// Parse payload JSON
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(notif.Payload), &payload); err != nil {
		log.Printf("[Email Service] Failed to parse payload JSON for notification %s: %v", notif.ID, err)
		return
	}

	docTitle, _ := payload["document_title"].(string)
	actorName, _ := payload["actor_name"].(string)
	uploaderName, _ := payload["uploader_name"].(string)
	message, _ := payload["message"].(string)

	var subject, body string

	switch notif.Template {
	case "action_required":
		subject = fmt.Sprintf("Action Required: %s", docTitle)
		body = fmt.Sprintf(`
			<h3>Action Required</h3>
			<p>Hello %s,</p>
			<p>A document titled <strong>"%s"</strong> has been submitted and is pending your review.</p>
			<p>Submitted by: %s</p>
			<p>Please log in to the portal to take action.</p>
		`, notif.Recipient.Name, docTitle, uploaderName)
	case "approved":
		subject = fmt.Sprintf("Document Approved: %s", docTitle)
		body = fmt.Sprintf(`
			<h3>Document Approved</h3>
			<p>Hello %s,</p>
			<p>Your document titled <strong>"%s"</strong> has been approved by <strong>%s</strong>.</p>
		`, notif.Recipient.Name, docTitle, actorName)
	case "rejected":
		subject = fmt.Sprintf("Document Rejected: %s", docTitle)
		body = fmt.Sprintf(`
			<h3>Document Rejected</h3>
			<p>Hello %s,</p>
			<p>Your document titled <strong>"%s"</strong> has been rejected by <strong>%s</strong>.</p>
		`, notif.Recipient.Name, docTitle, actorName)
	case "sent_back":
		subject = fmt.Sprintf("Document Sent Back: %s", docTitle)
		body = fmt.Sprintf(`
			<h3>Document Sent Back for Revision</h3>
			<p>Hello %s,</p>
			<p>Your document titled <strong>"%s"</strong> has been sent back for revision by <strong>%s</strong>.</p>
			<p>Please review any remarks and update the document.</p>
		`, notif.Recipient.Name, docTitle, actorName)
	case "closed":
		subject = fmt.Sprintf("Document Signed and Closed: %s", docTitle)
		body = fmt.Sprintf(`
			<h3>Document Closed</h3>
			<p>Hello %s,</p>
			<p>Your document titled <strong>"%s"</strong> has been successfully signed and closed by <strong>%s</strong>.</p>
		`, notif.Recipient.Name, docTitle, actorName)
	case "sla_warning":
		subject = fmt.Sprintf("SLA Breach Alert: %s", docTitle)
		body = fmt.Sprintf(`
			<h3>SLA Breach Warning</h3>
			<p>Hello %s,</p>
			<p>%s</p>
		`, notif.Recipient.Name, message)
	default:
		subject = fmt.Sprintf("Document Management Update")
		body = fmt.Sprintf(`
			<h3>Update Notification</h3>
			<p>Hello %s,</p>
			<p>There is a new update regarding a document action on the portal.</p>
		`, notif.Recipient.Name)
	}

	err := SendMail(cfg, []string{notif.Recipient.Email}, subject, body)
	if err != nil {
		log.Printf("[Email Service] Failed to send email to %s: %v", notif.Recipient.Email, err)
		db.Model(&notif).Update("status", "failed")
		return
	}

	log.Printf("[Email Service] Email successfully sent to %s for template '%s'", notif.Recipient.Email, notif.Template)
	now := time.Now()
	db.Model(&notif).Updates(map[string]interface{}{"status": "sent", "sent_at": &now})
}

func SendMail(cfg *config.Config, to []string, subject, body string) error {
	header := make(map[string]string)
	header["From"] = cfg.SMTPFrom
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPassword, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	return smtp.SendMail(addr, auth, cfg.SMTPUser, to, []byte(message))
}
