package cars

import (
	"context"
	"log"

	"github.com/nadirbasalamah/go-vrent/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Car represents car model in database
type Car struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Make      string             `bson:"make"`
	Name      string             `bson:"name"`
	Available bool               `bson:"available"`
}

// Add func to add new car data
func (car Car) Add() Car {
	collection := database.Database.Collection("car")

	data := Car{
		Make:      car.Make,
		Name:      car.Name,
		Available: true,
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Fatalf("Error occured when inserting: %v", err)
		return Car{}
	}

	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatalf("Error occured when inserting: %v", err)
		return Car{}
	}
	return data
}

// Edit func to edit car data
func (car Car) Edit() Car {
	collection := database.Database.Collection("car")

	data := &Car{}
	filter := bson.M{"_id": car.ID}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		log.Fatalf("Cannot find car with specified ID: %v", err)
		return Car{}
	}

	data.Make = car.Make
	data.Name = car.Name
	data.Available = car.Available

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		log.Fatalf("Cannot update object in mongoDB: %v", updateErr)
		return Car{}
	}
	return car
}

// Delete func to delete a car data
func (car Car) Delete() primitive.ObjectID {
	collection := database.Database.Collection("car")

	filter := bson.M{"_id": car.ID}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatalf("Cannot update object in mongoDB: %v", err)
		return primitive.ObjectID{}
	}
	return car.ID
}

// Get func to get all cars data
func Get() []Car {
	collection := database.Database.Collection("car")
	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		log.Fatalf("Cannot retrieve all cars data: %v", err)
		return []Car{}
	}

	defer cur.Close(context.Background())
	var cars []Car

	for cur.Next(context.Background()) {
		var car Car
		err := cur.Decode(&car)
		if err != nil {
			log.Fatalf("Cannot decode data: %v", err)
			return []Car{}
		}
		cars = append(cars, car)
	}
	return cars
}

// FindCar func to find car data by id
func FindCar(carID primitive.ObjectID) *Car {
	collection := database.Database.Collection("car")

	data := &Car{}
	filter := bson.M{"_id": carID}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		log.Fatalf("Cannot find car with specified ID: %v", err)
		return &Car{}
	}

	return data
}
