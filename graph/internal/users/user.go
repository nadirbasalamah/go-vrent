package users

import (
	"context"
	"log"

	"github.com/nadirbasalamah/go-vrent/database"
	"github.com/nadirbasalamah/go-vrent/graph/internal/cars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User represents user model for database
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Password    string             `bson:"password"`
	Reservation *cars.Car          `bson:"reservation"`
}

// Register register new user
func (user *User) Register() {
	collection := database.Database.Collection("user")

	bs, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		log.Fatalf("Error when generating password: %v", err)
	}

	password := string(bs)

	data := User{
		Name:        user.Name,
		Password:    password,
		Reservation: nil,
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatalf("Error occured when inserting: %v", err)
	}

	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatalf("Error occured when inserting: %v", err)
	}
}

// Authenticate check if user's name and password is match
func (user *User) Authenticate() bool {
	collection := database.Database.Collection("user")

	data := &User{}
	filter := bson.M{"name": user.Name}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		log.Fatalf("Cannot find user with specified name: %v", err)
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(user.Password))
	return err == nil
}

// GetUserIDByUsername check if a user exists in database by given username
func GetUserIDByUsername(username string) (primitive.ObjectID, error) {
	collection := database.Database.Collection("user")

	data := &User{}
	filter := bson.M{"name": username}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		log.Fatalf("Cannot find user with specified ID: %v", err)
		return primitive.ObjectID{}, err
	}

	return data.ID, nil
}

// Get func to get all users data
func Get() []User {
	collection := database.Database.Collection("user")
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		log.Fatalf("Cannot retrieve all users data: %v", err)
		return []User{}
	}

	defer cur.Close(context.Background())
	var users []User

	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			log.Fatalf("Cannot decode data: %v", err)
			return []User{}
		}
		users = append(users, user)
	}
	return users
}

// ReserveCar to reserve a car for authenticated users
func (user *User) ReserveCar(carID primitive.ObjectID) User {
	collection := database.Database.Collection("user")

	data := &User{}
	filter := bson.M{"_id": user.ID}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		log.Fatalf("Cannot find user with specified ID: %v", err)
		return User{}
	}

	reservedCar := cars.FindCar(carID)
	data.Reservation = reservedCar

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		log.Fatalf("Cannot update object in mongoDB: %v", updateErr)
		return User{}
	}
	return *data
}
