package repository

import (
	"fmt"

	"github.com/TemaStatham/Little-Messenger/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	publicTypeName  = "public"
	privateTypeName = "private"
)

// ChatPostgres :
type ChatPostgres struct {
	db *sqlx.DB
}

// NewChatPostgres :
func NewChatPostgres(db *sqlx.DB) *ChatPostgres {
	return &ChatPostgres{db: db}
}

// CreatePublicChat создает общую беседу
func (c *ChatPostgres) CreatePublicChat(creatorID uint, name string) error {
	chatID, err := c.CreateChat(publicTypeName)
	if err != nil {
		return fmt.Errorf("error creating public chat: %v", err)
	}

	query := `
		INSERT INTO public_chats (name, creation_date, creator_user_id, chat_id)
		VALUES ($1, CURRENT_DATE, $2, $3)
	`
	_, err = c.db.Exec(query, name, creatorID, chatID)
	if err != nil {
		return fmt.Errorf("error creating public chat: %v", err)
	}

	err = c.CreateChatMember(creatorID, chatID)
	if err != nil {
		return fmt.Errorf("error creating public chat: %v", err)
	}

	return nil
}

// CreatePrivateChat создает личные сообщения
func (c *ChatPostgres) CreatePrivateChat(user1ID, user2ID uint) error {
	chatID, err := c.CreateChat(privateTypeName)
	if err != nil {
		return fmt.Errorf("error creating private chat: %v", err)
	}

	query := `
		INSERT INTO private_chats (user1_id, user2_id, chat_id)
		VALUES ($1, $2, $3)
	`
	_, err = c.db.Exec(query, user1ID, user2ID, chatID)
	if err != nil {
		return fmt.Errorf("error creating private chat: %v", err)
	}

	return nil
}

// CreateChat создает общий чат в таблице chats
func (c *ChatPostgres) CreateChat(chatType string) (uint, error) {
	query := `
		INSERT INTO chats (type_name)
		VALUES ($1)
		RETURNING id
	`

	var chatID uint
	err := c.db.Get(&chatID, query, chatType)
	if err != nil {
		return 0, fmt.Errorf("error creating chat: %v", err)
	}

	return chatID, nil
}

// CreateChatMember добавляет пользователя к участникам публичной беседы
func (c *ChatPostgres) CreateChatMember(userID, chatID uint) error {
	query := `
		INSERT INTO chat_members (user_id, chat_id)
		VALUES ($1, $2)
	`
	_, err := c.db.Exec(query, userID, chatID)
	if err != nil {
		return fmt.Errorf("error create chat member: %v", err)
	}

	return nil
}

// GetUserPublicChats получает публичные чаты, в которых участвует пользователь
func (c *ChatPostgres) GetUserPublicChats(userID uint) ([]models.Conversation, error) {
	query := `
        SELECT pc.chat_id, pc.name, pc.creation_date, u.id as creator_id, u.username as creator_username
        FROM public_chats pc
        JOIN chat_members cm ON pc.chat_id = cm.chat_id
        JOIN users u ON pc.creator_user_id = u.id
        WHERE cm.user_id = $1
    `
	var userPublicChats []models.Conversation
	if err := c.db.Select(&userPublicChats, query, userID); err != nil {
		return nil, fmt.Errorf("error getting user chats: %v", err)
	}

	// Заполнение Messages нужно добавить в зависимости от вашей логики,
	// например, может потребоваться получить сообщения для каждого чата
	// for i := range userPublicChats {
	// 	messages, err := c.GetChatMessages(userPublicChats[i].ID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error getting chat messages: %v", err)
	// 	}
	// 	userPublicChats[i].Messages = messages
	// }

	return userPublicChats, nil
}

// GetUserPrivateChats получает личные чаты, в которых участвует пользователь
func (c *ChatPostgres) GetUserPrivateChats(userID uint) ([]models.Chat, error) {
	query := `
        SELECT pc.chat_id, u1.id as user1_id, u1.username as user1_username,
               u2.id as user2_id, u2.username as user2_username
        FROM private_chats pc
        JOIN users u1 ON pc.user1_id = u1.id
        JOIN users u2 ON pc.user2_id = u2.id
        WHERE u1.id = $1 OR u2.id = $1
    `
	var userPrivateChats []models.Chat
	if err := c.db.Select(&userPrivateChats, query, userID); err != nil {
		return nil, fmt.Errorf("error getting user private chats: %v", err)
	}

	// Заполнение Messages нужно добавить в зависимости от вашей логики,
	// например, может потребоваться получить сообщения для каждого чата
	// for i := range userPrivateChats {
	// 	messages, err := c.GetChatMessages(userPrivateChats[i].ID)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("error getting chat messages: %v", err)
	// 	}
	// 	userPrivateChats[i].Messages = messages
	// }

	return userPrivateChats, nil
}

// GetChatMessages получает сообщения для указанного чата
func (c *ChatPostgres) GetChatMessages(chatID uint) ([]models.Message, error) {
	query := `
        SELECT id, user_id, content, send_time
        FROM messages
        WHERE chat_id = $1
    `
	var messages []models.Message
	if err := c.db.Select(&messages, query, chatID); err != nil {
		return nil, fmt.Errorf("error getting chat messages: %v", err)
	}
	return messages, nil
}
