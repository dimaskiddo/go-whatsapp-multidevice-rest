package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/env"
	"github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/log"
)

// Webhook configuration loaded from .env
var (
	webhookURL   string
	webhookTypes []string // parsed from WEBHOOK_TYPE (comma-separated: user,group,status)
)

func init() {
	// Load webhook URL (optional, can be empty)
	url, err := env.GetEnvString("WEBHOOK_URL")
	if err == nil && len(url) > 0 {
		webhookURL = url
		log.Print(nil).Infof("Webhook URL configured: %s", webhookURL)
	} else {
		log.Print(nil).Info("No webhook URL configured (WEBHOOK_URL is empty)")
	}

	// Load webhook types (default: all three)
	rawTypes, err := env.GetEnvString("WEBHOOK_TYPE")
	if err == nil && len(rawTypes) > 0 {
		parts := strings.Split(rawTypes, ",")
		for _, p := range parts {
			t := strings.TrimSpace(strings.ToLower(p))
			if t == "user" || t == "group" || t == "status" {
				webhookTypes = append(webhookTypes, t)
			}
		}
		log.Print(nil).Infof("Webhook types configured: %s", strings.Join(webhookTypes, ", "))
	} else {
		// Default: enable all types
		webhookTypes = []string{"user", "group", "status"}
		log.Print(nil).Info("No WEBHOOK_TYPE set, defaulting to user,group,status")
	}
}

// WebhookPayload is the JSON payload sent to the webhook URL.
type WebhookPayload struct {
	JID       string `json:"jid"`
	From      string `json:"from"`
	Chat      string `json:"chat"`
	Message   string `json:"message"`
	MessageID string `json:"message_id"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source"`
}

// getMessageType determines the message category (user/group/status) from the chat JID.
func getMessageType(chat types.JID) string {
	switch chat.Server {
	case types.GroupServer:
		return "group"
	case types.BroadcastServer:
		return "status"
	default:
		return "user"
	}
}

// shouldSendWebhook checks if the given message type should trigger a webhook.
func shouldSendWebhook(msgType string) bool {
	if len(webhookURL) == 0 {
		return false
	}
	for _, t := range webhookTypes {
		if t == msgType {
			return true
		}
	}
	return false
}

// sendWebhook sends a POST request to the configured webhook URL with message details.
func sendWebhook(jid string, v *events.Message, msgType string, messageText string, senderDisplay string, chatDisplay string) {
	if len(webhookURL) == 0 {
		return
	}

	payload := WebhookPayload{
		JID:       jid,
		From:      senderDisplay,
		Chat:      chatDisplay,
		Message:   messageText,
		MessageID: v.Info.ID,
		Type:      msgType,
		Timestamp: v.Info.Timestamp.Unix(),
		Source:    v.Info.SourceString(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Print(nil).Errorf("[%s] Failed to marshal webhook payload: %v", maskJID(jid), err)
		return
	}

	// Send POST request with JSON body
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewReader(jsonData))
		if err != nil {
			log.Print(nil).Errorf("[%s] Failed to create webhook request: %v", maskJID(jid), err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Print(nil).Errorf("[%s] Webhook request failed: %v", maskJID(jid), err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			log.Print(nil).Warnf("[%s] Webhook returned non-2xx status: %d", maskJID(jid), resp.StatusCode)
		} else {
			log.Print(nil).Debugf("[%s] Webhook sent successfully (%s, type=%s)", maskJID(jid), v.Info.ID, msgType)
		}
	}()
}

// WhatsAppRegisterEventHandler registers event handlers on a whatsmeow client
// to handle incoming messages, connection state changes, history sync, and more.
// This is critical for WhatsApp to consider this client as a proper web client.
func WhatsAppRegisterEventHandler(jid string) {
	if WhatsAppClient[jid] == nil {
		return
	}

	WhatsAppClient[jid].AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			handleMessage(jid, v)
		case *events.Receipt:
			handleReceipt(jid, v)
		case *events.Connected:
			handleConnected(jid, v)
		case *events.Disconnected:
			handleDisconnected(jid, v)
		case *events.LoggedOut:
			handleLoggedOut(jid, v)
		case *events.StreamReplaced:
			handleStreamReplaced(jid, v)
		case *events.PairSuccess:
			handlePairSuccess(jid, v)
		case *events.HistorySync:
			handleHistorySync(jid, v)
		case *events.IdentityChange:
			handleIdentityChange(jid, v)
		case *events.Presence:
			// Optional: track contact presence
		case *events.ChatPresence:
			// Optional: track typing indicators
		case *events.KeepAliveTimeout:
			log.Print(nil).Warnf("[%s] Keepalive timeout (%d errors, last success: %v)", maskJID(jid), v.ErrorCount, v.LastSuccess)
		case *events.KeepAliveRestored:
			log.Print(nil).Infof("[%s] Keepalive restored", maskJID(jid))
		case *events.StreamError:
			log.Print(nil).Errorf("[%s] Stream error: code=%s", maskJID(jid), v.Code)
		case *events.ConnectFailure:
			log.Print(nil).Errorf("[%s] Connect failure: reason=%s", maskJID(jid), v.Reason)
		case *events.TemporaryBan:
			log.Print(nil).Warnf("[%s] Temporary ban: %s (expires in %v)", maskJID(jid), v.Code, v.Expire)
		case *events.OfflineSyncPreview:
			log.Print(nil).Infof("[%s] Offline sync preview: %d messages, %d receipts, %d notifications",
				maskJID(jid), v.Messages, v.Receipts, v.Notifications)
		case *events.OfflineSyncCompleted:
			log.Print(nil).Infof("[%s] Offline sync completed: %d messages", maskJID(jid), v.Count)
		}
	})
}

// resolvePhoneNumber resolves a JID to a human-readable phone number.
// In newer WhatsApp, users are identified by LIDs (@lid) instead of phone numbers.
// The actual phone number is available via SenderAlt/RecipientAlt in MessageSource,
// or via the store's LID-to-PN mapping.
func resolvePhoneNumber(jid types.JID, altJID types.JID) string {
	// If we have an alternative JID (phone number), use it
	if !altJID.IsEmpty() && altJID.Server == types.DefaultUserServer {
		return altJID.User
	}

	// If the JID itself is already a phone number, use it directly
	if jid.Server == types.DefaultUserServer {
		return jid.User
	}

	// For LIDs or other server types, return the full JID as fallback
	return jid.String()
}

// resolveChatNumber resolves the chat JID to a human-readable identifier.
// For DMs, this shows the phone number; for groups, it shows the group JID.
func resolveChatNumber(chat types.JID, recipientAlt types.JID) string {
	if chat.Server == types.GroupServer {
		return chat.String()
	}

	// For DMs and status, try to resolve to phone number
	if !recipientAlt.IsEmpty() && recipientAlt.Server == types.DefaultUserServer {
		return recipientAlt.User
	}
	if chat.Server == types.DefaultUserServer {
		return chat.User
	}
	return chat.String()
}

// extractMessageText extracts the human-readable text content from a message.
func extractMessageText(v *events.Message) string {
	if v.Message == nil {
		return ""
	}

	switch {
	case v.Message.GetConversation() != "":
		return v.Message.GetConversation()
	case v.Message.GetExtendedTextMessage() != nil:
		return v.Message.GetExtendedTextMessage().GetText()
	case v.Message.GetImageMessage() != nil:
		cap := v.Message.GetImageMessage().GetCaption()
		if cap != "" {
			return cap
		}
		return "[Image]"
	case v.Message.GetVideoMessage() != nil:
		cap := v.Message.GetVideoMessage().GetCaption()
		if cap != "" {
			return cap
		}
		return "[Video]"
	case v.Message.GetDocumentMessage() != nil:
		cap := v.Message.GetDocumentMessage().GetCaption()
		if cap != "" {
			return cap
		}
		return "[Document]"
	case v.Message.GetAudioMessage() != nil:
		return "[Audio]"
	case v.Message.GetStickerMessage() != nil:
		return "[Sticker]"
	case v.Message.GetLocationMessage() != nil:
		return fmt.Sprintf("[Location] %.6f,%.6f",
			v.Message.GetLocationMessage().GetDegreesLatitude(),
			v.Message.GetLocationMessage().GetDegreesLongitude())
	case v.Message.GetContactMessage() != nil:
		return "[Contact] " + v.Message.GetContactMessage().GetDisplayName()
	case v.Message.GetPollCreationMessage() != nil:
		return "[Poll] " + v.Message.GetPollCreationMessage().GetName()
	case v.Message.GetReactionMessage() != nil:
		return "[Reaction] " + v.Message.GetReactionMessage().GetText()
	case v.Message.GetProtocolMessage() != nil:
		return "[Protocol]"
	default:
		return "[Unknown type]"
	}
}

// handleMessage processes incoming messages:
// - Logs the message with resolved phone numbers
// - Sends read receipt to acknowledge receipt
// - Sends webhook if configured
func handleMessage(jid string, v *events.Message) {
	info := v.Info
	chat := info.Chat
	//messageID := info.ID

	// Resolve sender and chat to phone numbers (LID → PN)
	senderDisplay := resolvePhoneNumber(info.Sender, info.SenderAlt)
	chatDisplay := resolveChatNumber(chat, info.RecipientAlt)

	// Extract message text/content
	messageText := extractMessageText(v)

	log.Print(nil).Infof("[%s] Message received from %s in %s: %s",
		maskJID(jid), senderDisplay, chatDisplay, messageText)

	// Determine message type for webhook filtering
	msgType := getMessageType(chat)

	// Send webhook if this message type should trigger one
	if shouldSendWebhook(msgType) {
		sendWebhook(jid, v, msgType, messageText, senderDisplay, chatDisplay)
	}

	// Send read receipt to acknowledge the message
	// This is critical for WhatsApp to consider this client as a proper web client
	// err := WhatsAppMarkRead(jid, chat, []string{messageID})
	// if err != nil {
	// 	log.Print(nil).Errorf("[%s] Failed to mark message %s as read: %v", maskJID(jid), messageID, err)
	// }
}

// handleReceipt processes delivery/read receipts from other users
func handleReceipt(jid string, v *events.Receipt) {
	log.Print(nil).Debugf("[%s] Receipt from %s: type=%s, messageIDs=%v",
		maskJID(jid), v.MessageSource.Chat, v.Type, v.MessageIDs)
}

// handleConnected handles successful connection
func handleConnected(jid string, v *events.Connected) {
	log.Print(nil).Infof("[%s] WhatsApp client connected successfully", maskJID(jid))

	// Set presence to available so WhatsApp knows we're active
	WhatsAppPresence(jid, true)

	// Log connection duration stats
	if client := WhatsAppClient[jid]; client != nil {
		lastConn := client.LastSuccessfulConnect
		if !lastConn.IsZero() {
			log.Print(nil).Infof("[%s] Connection established, last successful was at %s",
				maskJID(jid), lastConn.Format(time.RFC3339))
		}
	}
}

// handleDisconnected handles disconnection events
func handleDisconnected(jid string, v *events.Disconnected) {
	log.Print(nil).Warnf("[%s] WhatsApp client disconnected", maskJID(jid))

	// The client has EnableAutoReconnect set to true, so it will try to reconnect automatically.
	// We just need to log the event and update state.
}

// handleLoggedOut handles when the device is logged out remotely
func handleLoggedOut(jid string, v *events.LoggedOut) {
	log.Print(nil).Warnf("[%s] WhatsApp client logged out (onConnect=%v, reason=%s)",
		maskJID(jid), v.OnConnect, v.Reason)

	// Clean up the client from our map
	if client := WhatsAppClient[jid]; client != nil {
		client.Disconnect()
		WhatsAppClient[jid] = nil
		delete(WhatsAppClient, jid)
	}
}

// handleStreamReplaced handles when another client connects with the same credentials
func handleStreamReplaced(jid string, v *events.StreamReplaced) {
	log.Print(nil).Warnf("[%s] WhatsApp stream replaced by another client", maskJID(jid))

	// Clean up the client from our map
	if client := WhatsAppClient[jid]; client != nil {
		WhatsAppClient[jid] = nil
		delete(WhatsAppClient, jid)
	}
}

// handlePairSuccess handles successful device pairing
func handlePairSuccess(jid string, v *events.PairSuccess) {
	log.Print(nil).Infof("[%s] WhatsApp pairing successful (platform=%s, business=%s)",
		maskJID(jid), v.Platform, v.BusinessName)
}

// handleHistorySync processes history sync blobs from the primary device.
// This is critical for WhatsApp to know the client has received historical messages.
func handleHistorySync(jid string, v *events.HistorySync) {
	if v.Data == nil {
		return
	}

	syncType := v.Data.GetSyncType()
	progress := v.Data.GetProgress()
	chunkOrder := v.Data.GetChunkOrder()

	log.Print(nil).Infof("[%s] History sync: type=%d, progress=%d, chunkOrder=%d",
		maskJID(jid), syncType, progress, chunkOrder)

	// Log conversation count if available
	conversations := v.Data.GetConversations()
	if conversations != nil {
		log.Print(nil).Infof("[%s] History sync contains %d conversations", maskJID(jid), len(conversations))
	}

	// Mark history sync as processed by the store
	// This acknowledges to WhatsApp that we've received the history data
	if client := WhatsAppClient[jid]; client != nil && client.Store != nil {
		// The history sync is already handled by whatsmeow's store internally
		log.Print(nil).Infof("[%s] History sync processed successfully", maskJID(jid))
	}
}

// handleIdentityChange handles when a contact changes their identity
func handleIdentityChange(jid string, v *events.IdentityChange) {
	log.Print(nil).Infof("[%s] Identity changed for %s (implicit=%v)",
		maskJID(jid), v.JID, v.Implicit)
}

// WhatsAppMarkRead sends read receipts for the specified message IDs in a chat.
// This is essential for WhatsApp to know the client has read messages.
// MarkRead takes: context, messageIDs, timestamp, chat JID, sender JID (empty for own messages).
// func WhatsAppMarkRead(jid string, chatJID types.JID, messageIDs []string) error {
// 	client, ok := WhatsAppClient[jid]
// 	if !ok || client == nil {
// 		return nil
// 	}

// 	if !client.IsConnected() || !client.IsLoggedIn() {
// 		log.Print(nil).Debugf("[%s] Cannot mark read: client not connected/logged in", maskJID(jid))
// 		return nil
// 	}

// 	err := client.MarkRead(context.Background(), messageIDs, time.Now(), chatJID, types.EmptyJID)
// 	if err != nil {
// 		log.Print(nil).Errorf("[%s] Failed to mark read: %v", maskJID(jid), err)
// 		return err
// 	}

// 	log.Print(nil).Debugf("[%s] Marked %d message(s) as read in %s", maskJID(jid), len(messageIDs), chatJID)
// 	return nil
// }

// maskJID masks the JID for logging (shows last 4 digits)
func maskJID(jid string) string {
	if len(jid) <= 4 {
		return "****"
	}
	return jid[:len(jid)-4] + "xxxx"
}
