package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nadirbasalamah/go-vrent/graph/auth"
	"github.com/nadirbasalamah/go-vrent/graph/generated"
	"github.com/nadirbasalamah/go-vrent/graph/internal/cars"
	"github.com/nadirbasalamah/go-vrent/graph/internal/users"
	"github.com/nadirbasalamah/go-vrent/graph/model"
	"github.com/nadirbasalamah/go-vrent/graph/pkg/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) AddCar(ctx context.Context, input model.NewCar) (*model.Car, error) {
	car := cars.Car{}
	car.Make = input.Make
	car.Name = input.Name
	car.Price = input.Price

	result := car.Add()
	return &model.Car{ID: result.ID.Hex(), Make: result.Make, Name: result.Name, Price: result.Price, Available: result.Available}, nil
}

func (r *mutationResolver) EditCar(ctx context.Context, input model.EditCar) (*model.Car, error) {
	car := cars.Car{}
	car.ID, _ = primitive.ObjectIDFromHex(input.CarID)
	car.Make = input.Make
	car.Name = input.Name
	car.Price = input.Price

	result := car.Edit()
	return &model.Car{ID: result.ID.Hex(), Make: result.Make, Name: result.Name, Price: result.Price, Available: result.Available}, nil
}

func (r *mutationResolver) DeleteCar(ctx context.Context, input *model.DeleteCar) (string, error) {
	car := cars.Car{}
	car.ID, _ = primitive.ObjectIDFromHex(input.CarID)

	result := car.Delete()
	return result.Hex(), nil
}

func (r *mutationResolver) AddUser(ctx context.Context, input model.NewUser) (string, error) {
	user := users.User{}
	user.Name = input.Username
	user.Password = input.Password
	user.Register()
	token, err := jwt.GenerateToken(user.Name)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	user := users.User{}
	user.Name = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Name)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) ReserveCar(ctx context.Context, input model.ReserveCar) (*model.User, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.User{}, fmt.Errorf("access denied")
	}

	carID, _ := primitive.ObjectIDFromHex(input.CarID)

	updatedUser := user.ReserveCar(carID)
	return &model.User{ID: updatedUser.ID.Hex(), Name: updatedUser.Name, Reservation: &model.Car{ID: updatedUser.Reservation.ID.Hex(), Make: updatedUser.Reservation.Make, Name: updatedUser.Reservation.Name, Price: updatedUser.Reservation.Price}}, nil
}

func (r *mutationResolver) ReturnCar(ctx context.Context, input model.ReturnCar) (*model.User, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.User{}, fmt.Errorf("access denied")
	}

	carID, _ := primitive.ObjectIDFromHex(input.CarID)

	updatedUser := user.ReturnCar(carID)
	return &model.User{ID: updatedUser.ID.Hex(), Name: updatedUser.Name, Reservation: nil}, nil
}

func (r *queryResolver) Cars(ctx context.Context) ([]*model.Car, error) {
	results := []*model.Car{}
	carsData := []cars.Car{}

	carsData = cars.Get()
	for _, car := range carsData {
		results = append(results, &model.Car{ID: car.ID.Hex(), Make: car.Make, Name: car.Name, Price: car.Price, Available: car.Available})
	}
	return results, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	results := []*model.User{}
	usersData := []users.User{}
	reservation := model.Car{}

	usersData = users.Get()
	for _, user := range usersData {
		if user.Reservation != nil {
			reservation = model.Car{ID: user.Reservation.ID.Hex(), Make: user.Reservation.Make, Name: user.Reservation.Name}
		}
		results = append(results, &model.User{ID: user.ID.Hex(), Name: user.Name, Reservation: &reservation})
	}
	return results, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
