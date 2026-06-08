package whatsapp

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"

	mediapkg "github.com/dimaskiddo/go-whatsapp-multidevice-rest/pkg/media"
)

// WebhookPayload is the unified envelope sent to the webhook URL for all events.
type WebhookPayload struct {
	Phone     string      `json:"phone"`
	Event     string      `json:"event"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// --- Message Event ---

// WebhookQuotedMessage holds a summary of the message being replied to.
type WebhookQuotedMessage struct {
	ID     string `json:"id"`
	Sender string `json:"sender"`
	Type   string `json:"type"`
	Text   string `json:"text,omitempty"`
}

type WebhookMessageData struct {
	ID           string    `json:"id"`
	From         string    `json:"from"`
	FromAlt      string    `json:"from_alt,omitempty"`
	FromName     string    `json:"from_name"`
	Chat         string    `json:"chat"`
	RecipientAlt string    `json:"recipient_alt,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	IsGroup      bool      `json:"is_group"`
	IsFromMe     bool      `json:"is_from_me"`
	IsEdit       bool      `json:"is_edit,omitempty"`
	IsViewOnce   bool      `json:"is_view_once,omitempty"`
	Type         string    `json:"type"`
	Text         string    `json:"text,omitempty"`
	MimeType     string    `json:"mime_type,omitempty"`
	// Media contains base64-encoded bytes only when MEDIA_STORAGE=none (default).
	Media string `json:"media,omitempty"`
	// MediaURL is populated when MEDIA_STORAGE=local or MEDIA_STORAGE=s3.
	MediaURL      string                `json:"media_url,omitempty"`
	QuotedMessage *WebhookQuotedMessage `json:"quoted_message,omitempty"`
}

// storeMedia downloads the given message attachment and either saves it via the
// configured storage backend (returning a URL) or returns raw base64 bytes.
func storeMedia(ctx context.Context, client *whatsmeow.Client, msg whatsmeow.DownloadableMessage, msgID, mimeType string) (media, mediaURL string) {
	if client == nil {
		return
	}
	data, err := client.Download(ctx, msg)
	if err != nil {
		return
	}

	if mediapkg.MediaStorage != nil {
		// Derive a deterministic filename from the message ID and MIME type.
		ext := mimeToExt(mimeType)
		filename := fmt.Sprintf("%s%s", msgID, ext)
		url, err := mediapkg.MediaStorage.Save(ctx, filename, mimeType, data)
		if err == nil {
			mediaURL = url
			return
		}
		// Fall through to base64 on error.
	}

	// Fallback / MEDIA_STORAGE=none: embed as base64.
	media = base64.StdEncoding.EncodeToString(data)
	return
}

// mimeToExt maps common MIME types to file extensions.
func mimeToExt(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	case "video/mp4":
		return ".mp4"
	case "video/3gpp":
		return ".3gp"
	case "audio/ogg":
		return ".ogg"
	case "audio/mpeg":
		return ".mp3"
	case "audio/aac":
		return ".aac"
	case "application/pdf":
		return ".pdf"
	}
	return ".bin"
}

// buildQuotedMessage extracts reply context from a contextInfo proto.
func buildQuotedMessage(ctx *waE2E.ContextInfo) *WebhookQuotedMessage {
	if ctx == nil {
		return nil
	}
	qID := ctx.GetStanzaID()
	if qID == "" {
		return nil
	}
	qSender := ctx.GetParticipant()
	qMsg := ctx.GetQuotedMessage()
	qType := "unknown"
	qText := ""
	if qMsg != nil {
		switch {
		case qMsg.GetConversation() != "":
			qType = "text"
			qText = qMsg.GetConversation()
		case qMsg.GetExtendedTextMessage() != nil:
			qType = "text"
			qText = qMsg.GetExtendedTextMessage().GetText()
		case qMsg.GetImageMessage() != nil:
			qType = "image"
			qText = qMsg.GetImageMessage().GetCaption()
		case qMsg.GetVideoMessage() != nil:
			qType = "video"
			qText = qMsg.GetVideoMessage().GetCaption()
		case qMsg.GetAudioMessage() != nil:
			qType = "audio"
		case qMsg.GetDocumentMessage() != nil:
			qType = "document"
			qText = qMsg.GetDocumentMessage().GetCaption()
		case qMsg.GetStickerMessage() != nil:
			qType = "sticker"
		case qMsg.GetLocationMessage() != nil:
			qType = "location"
		case qMsg.GetContactMessage() != nil:
			qType = "contact"
		case qMsg.GetReactionMessage() != nil:
			qType = "reaction"
			qText = qMsg.GetReactionMessage().GetText()
		}
	}
	return &WebhookQuotedMessage{
		ID:     qID,
		Sender: qSender,
		Type:   qType,
		Text:   qText,
	}
}

func buildMessagePayload(client *whatsmeow.Client, v *events.Message) WebhookPayload {
	msgType := "unknown"
	msgText := ""
	msgMime := ""
	msgMedia := ""
	msgMediaURL := ""
	var msgQuoted *WebhookQuotedMessage

	msg := v.Message
	if msg != nil {
		switch {
		case msg.GetConversation() != "":
			msgType = "text"
			msgText = msg.GetConversation()
		case msg.GetExtendedTextMessage() != nil:
			msgType = "text"
			msgText = msg.GetExtendedTextMessage().GetText()
			msgQuoted = buildQuotedMessage(msg.GetExtendedTextMessage().GetContextInfo())
		case msg.GetImageMessage() != nil:
			msgType = "image"
			msgText = msg.GetImageMessage().GetCaption()
			msgMime = msg.GetImageMessage().GetMimetype()
			msgMedia, msgMediaURL = storeMedia(context.Background(), client, msg.GetImageMessage(), v.Info.ID, msgMime)
			msgQuoted = buildQuotedMessage(msg.GetImageMessage().GetContextInfo())
		case msg.GetVideoMessage() != nil:
			msgType = "video"
			msgText = msg.GetVideoMessage().GetCaption()
			msgMime = msg.GetVideoMessage().GetMimetype()
			msgMedia, msgMediaURL = storeMedia(context.Background(), client, msg.GetVideoMessage(), v.Info.ID, msgMime)
			msgQuoted = buildQuotedMessage(msg.GetVideoMessage().GetContextInfo())
		case msg.GetAudioMessage() != nil:
			msgType = "audio"
			msgMime = msg.GetAudioMessage().GetMimetype()
			msgMedia, msgMediaURL = storeMedia(context.Background(), client, msg.GetAudioMessage(), v.Info.ID, msgMime)
			msgQuoted = buildQuotedMessage(msg.GetAudioMessage().GetContextInfo())
		case msg.GetDocumentMessage() != nil:
			msgType = "document"
			msgText = msg.GetDocumentMessage().GetCaption()
			msgMime = msg.GetDocumentMessage().GetMimetype()
			msgMedia, msgMediaURL = storeMedia(context.Background(), client, msg.GetDocumentMessage(), v.Info.ID, msgMime)
			msgQuoted = buildQuotedMessage(msg.GetDocumentMessage().GetContextInfo())
		case msg.GetStickerMessage() != nil:
			msgType = "sticker"
			msgMime = msg.GetStickerMessage().GetMimetype()
			msgMedia, msgMediaURL = storeMedia(context.Background(), client, msg.GetStickerMessage(), v.Info.ID, msgMime)
			msgQuoted = buildQuotedMessage(msg.GetStickerMessage().GetContextInfo())
		case msg.GetLocationMessage() != nil:
			msgType = "location"
			msgQuoted = buildQuotedMessage(msg.GetLocationMessage().GetContextInfo())
		case msg.GetContactMessage() != nil:
			msgType = "contact"
			msgQuoted = buildQuotedMessage(msg.GetContactMessage().GetContextInfo())
		case msg.GetReactionMessage() != nil:
			msgType = "reaction"
			msgText = msg.GetReactionMessage().GetText()
		case msg.GetPollCreationMessage() != nil:
			msgType = "poll"
			msgText = msg.GetPollCreationMessage().GetName()
			msgQuoted = buildQuotedMessage(msg.GetPollCreationMessage().GetContextInfo())
		}
	}

	senderAlt := ""
	if !v.Info.SenderAlt.IsEmpty() {
		senderAlt = v.Info.SenderAlt.String()
	}
	recipientAlt := ""
	if !v.Info.RecipientAlt.IsEmpty() {
		recipientAlt = v.Info.RecipientAlt.String()
	}

	return WebhookPayload{
		Event:     "message",
		Timestamp: time.Now(),
		Data: WebhookMessageData{
			ID:            v.Info.ID,
			From:          v.Info.Sender.String(),
			FromAlt:       senderAlt,
			FromName:      v.Info.PushName,
			Chat:          v.Info.Chat.String(),
			RecipientAlt:  recipientAlt,
			Timestamp:     v.Info.Timestamp,
			IsGroup:       v.Info.IsGroup,
			IsFromMe:      v.Info.IsFromMe,
			IsEdit:        v.IsEdit,
			IsViewOnce:    v.IsViewOnce,
			Type:          msgType,
			Text:          msgText,
			MimeType:      msgMime,
			Media:         msgMedia,
			MediaURL:      msgMediaURL,
			QuotedMessage: msgQuoted,
		},
	}
}

// --- Receipt Event ---

type WebhookReceiptData struct {
	MessageIDs    []string          `json:"message_ids"`
	From          string            `json:"from"`
	Chat          string            `json:"chat"`
	IsGroup       bool              `json:"is_group"`
	Timestamp     time.Time         `json:"timestamp"`
	Type          types.ReceiptType `json:"type"`
	MessageSender string            `json:"message_sender,omitempty"`
}

func buildReceiptPayload(v *events.Receipt) WebhookPayload {
	msgSender := ""
	if !v.MessageSender.IsEmpty() {
		msgSender = v.MessageSender.String()
	}
	ids := make([]string, len(v.MessageIDs))
	copy(ids, v.MessageIDs)
	return WebhookPayload{
		Event:     "receipt",
		Timestamp: time.Now(),
		Data: WebhookReceiptData{
			MessageIDs:    ids,
			From:          v.Sender.String(),
			Chat:          v.Chat.String(),
			IsGroup:       v.IsGroup,
			Timestamp:     v.Timestamp,
			Type:          v.Type,
			MessageSender: msgSender,
		},
	}
}

// --- Chat Presence Event ---

type WebhookChatPresenceData struct {
	From  string `json:"from"`
	Chat  string `json:"chat"`
	State string `json:"state"`
	Media string `json:"media"`
}

func buildChatPresencePayload(v *events.ChatPresence) WebhookPayload {
	return WebhookPayload{
		Event:     "chat_presence",
		Timestamp: time.Now(),
		Data: WebhookChatPresenceData{
			From:  v.Sender.String(),
			Chat:  v.Chat.String(),
			State: string(v.State),
			Media: string(v.Media),
		},
	}
}

// --- Presence Event ---

type WebhookPresenceData struct {
	From        string    `json:"from"`
	Unavailable bool      `json:"unavailable"`
	LastSeen    time.Time `json:"last_seen,omitempty"`
}

func buildPresencePayload(v *events.Presence) WebhookPayload {
	return WebhookPayload{
		Event:     "presence",
		Timestamp: time.Now(),
		Data: WebhookPresenceData{
			From:        v.From.String(),
			Unavailable: v.Unavailable,
			LastSeen:    v.LastSeen,
		},
	}
}

// --- Group Info Event ---

type WebhookGroupInfoData struct {
	JID       string    `json:"jid"`
	Sender    string    `json:"sender,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Name      string    `json:"name,omitempty"`
	Topic     string    `json:"topic,omitempty"`
	Join      []string  `json:"join,omitempty"`
	Leave     []string  `json:"leave,omitempty"`
	Promote   []string  `json:"promote,omitempty"`
	Demote    []string  `json:"demote,omitempty"`
}

func jidSliceToStrings(jids []types.JID) []string {
	if len(jids) == 0 {
		return nil
	}
	result := make([]string, len(jids))
	for i, j := range jids {
		result[i] = j.String()
	}
	return result
}

func buildGroupInfoPayload(v *events.GroupInfo) WebhookPayload {
	sender := ""
	if v.Sender != nil {
		sender = v.Sender.String()
	}
	name := ""
	if v.Name != nil {
		name = v.Name.Name
	}
	topic := ""
	if v.Topic != nil {
		topic = v.Topic.Topic
	}
	return WebhookPayload{
		Event:     "group_info",
		Timestamp: time.Now(),
		Data: WebhookGroupInfoData{
			JID:       v.JID.String(),
			Sender:    sender,
			Timestamp: v.Timestamp,
			Name:      name,
			Topic:     topic,
			Join:      jidSliceToStrings(v.Join),
			Leave:     jidSliceToStrings(v.Leave),
			Promote:   jidSliceToStrings(v.Promote),
			Demote:    jidSliceToStrings(v.Demote),
		},
	}
}

// --- Joined Group Event ---

type WebhookJoinedGroupData struct {
	JID    string `json:"jid"`
	Name   string `json:"name"`
	Reason string `json:"reason,omitempty"`
	Type   string `json:"type,omitempty"`
}

func buildJoinedGroupPayload(v *events.JoinedGroup) WebhookPayload {
	return WebhookPayload{
		Event:     "joined_group",
		Timestamp: time.Now(),
		Data: WebhookJoinedGroupData{
			JID:    v.JID.String(),
			Name:   v.GroupInfo.Name,
			Reason: v.Reason,
			Type:   v.Type,
		},
	}
}

// --- Picture Event ---

type WebhookPictureData struct {
	JID       string    `json:"jid"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Removed   bool      `json:"removed"`
	PictureID string    `json:"picture_id,omitempty"`
}

func buildPicturePayload(v *events.Picture) WebhookPayload {
	return WebhookPayload{
		Event:     "picture",
		Timestamp: time.Now(),
		Data: WebhookPictureData{
			JID:       v.JID.String(),
			Author:    v.Author.String(),
			Timestamp: v.Timestamp,
			Removed:   v.Remove,
			PictureID: v.PictureID,
		},
	}
}

// --- Connection Events ---

type WebhookConnectionData struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

func buildConnectedPayload() WebhookPayload {
	return WebhookPayload{
		Event:     "connection",
		Timestamp: time.Now(),
		Data:      WebhookConnectionData{Status: "connected"},
	}
}

func buildDisconnectedPayload() WebhookPayload {
	return WebhookPayload{
		Event:     "connection",
		Timestamp: time.Now(),
		Data:      WebhookConnectionData{Status: "disconnected"},
	}
}

func buildLoggedOutPayload(v *events.LoggedOut) WebhookPayload {
	return WebhookPayload{
		Event:     "connection",
		Timestamp: time.Now(),
		Data:      WebhookConnectionData{Status: "logged_out", Reason: v.Reason.String()},
	}
}

// --- Call Events ---

type WebhookCallData struct {
	CallID    string    `json:"call_id"`
	From      string    `json:"from"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason,omitempty"`
}

func buildCallOfferPayload(v *events.CallOffer) WebhookPayload {
	return WebhookPayload{
		Event:     "call",
		Timestamp: time.Now(),
		Data: WebhookCallData{
			CallID:    v.CallID,
			From:      v.From.String(),
			Timestamp: v.Timestamp,
			Status:    "offer",
		},
	}
}

func buildCallTerminatePayload(v *events.CallTerminate) WebhookPayload {
	return WebhookPayload{
		Event:     "call",
		Timestamp: time.Now(),
		Data: WebhookCallData{
			CallID:    v.CallID,
			From:      v.From.String(),
			Timestamp: v.Timestamp,
			Status:    "terminate",
			Reason:    v.Reason,
		},
	}
}

func buildCallRejectPayload(v *events.CallReject) WebhookPayload {
	return WebhookPayload{
		Event:     "call",
		Timestamp: time.Now(),
		Data: WebhookCallData{
			CallID:    v.CallID,
			From:      v.From.String(),
			Timestamp: v.Timestamp,
			Status:    "reject",
		},
	}
}
