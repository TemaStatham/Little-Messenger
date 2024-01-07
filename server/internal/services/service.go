package services

import (
	"github.com/TemaStatham/Little-Messenger/internal/models"
	"github.com/TemaStatham/Little-Messenger/internal/repository"
)

// User - интерфейс, определяющий сервисы, связанные с пользователями.
type User interface {
	CreateUser(user *models.User) (uint, error)
	GetUserByEmail(email, password string) (models.User, error)
	GetUserByID(userID uint) (models.User, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(accessToken string) (uint, error)
	GetUsers() ([]models.Contact, error)
	CreateContact(userID1 string, userID2 string) error
}

// Chat - интерфейс, определяющий сервисы, связанные с чатами.
type Chat interface {
	CreateChatMember(userID, chatID uint) error
	CreatePrivateChat(user1ID, user2ID uint) error
	CreatePublicChat(creatorID uint, name string) error
}

// Service представляет основной уровень сервиса, объединяющий сервисы User, Chat и WebSocket.
type Service struct {
	User
	Chat
}

// NewService создает новый экземпляр Service с предоставленным репозиторием.
func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
		Chat: NewChatService(repos.Chat),
	}
}
