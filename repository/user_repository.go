package repository

import (
	"context"
	"log"

	"github.com/hongdangcseiu/go-back-end/domain"
	"github.com/hongdangcseiu/go-back-end/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(db mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}

func (ur *userRepository) Create(c context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)
	_, err := collection.InsertOne(c, user)
	return err
}

func (ur *userRepository) Fetch(c context.Context) ([]domain.User, error) {
	collection := ur.database.Collection(ur.collection)
	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})
	cursor, err := collection.Find(c, bson.D{}, opts)

	if err != nil {
		return nil, err
	}

	var users []domain.User
	err = cursor.All(c, &users)
	if users == nil {
		return []domain.User{}, err
	}

	return users, err
}

func (ur *userRepository) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	log.Print("user_repository.GetUserByEmail handler...")
	collection := ur.database.Collection(ur.collection)
	var user domain.User
	err := collection.FindOne(c, bson.M{"email": email}).Decode(&user)
	return user, err
}

func (ur *userRepository) GetUserByID(c context.Context, id string) (domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	var user domain.User

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)
	return user, err
}
func (ur *userRepository) GetUserByUserName(c context.Context, username string) (domain.User, error) {
	log.Print("user_repository.GetUserByUserName handler...")
	collection := ur.database.Collection(ur.collection)
	var user domain.User
	err := collection.FindOne(c, bson.M{"username": username}).Decode(&user)
	return user, err
}

func (ur *userRepository) UpdateUser(c context.Context, user domain.User) error {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"_id": user.ID}
	log.Print("postrepository.edit: ", filter)
	update := bson.M{
		"$set": bson.M{
			"name":         user.Name,
			"bio":          user.Bio,
			"profile_pic":  user.ProfilePic,
			"social_media": user.SocialMedia,
		},
	}
	log.Print("userrepository.edit: update:", update)

	_, err := collection.UpdateOne(c, filter, update)
	if err != nil {
		return err
	}

	return nil
}
