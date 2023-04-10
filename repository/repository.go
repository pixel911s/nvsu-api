package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"nvsu-api/model"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrProductNotFound = errors.New("product not found")
	ErrOrderNotFound   = errors.New("order not found")
)

type repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) Repository {
	return &repository{db: db}
}

func (r repository) GetUser(ctx context.Context, email string) (model.User, error) {
	var out user
	err := r.db.
		Collection("users").
		FindOne(ctx, bson.M{"email": email}).
		Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}
	return toModel(out), nil

}

func (r repository) GetProduct(ctx context.Context, id string) (model.Product, error) {
	var out product

	oid, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": oid}

	err := r.db.
		Collection("products").
		FindOne(ctx, query).
		Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Product{}, ErrProductNotFound
		}
		return model.Product{}, err
	}
	return toProductModel(out), nil
}

func (r repository) GetOrder(ctx context.Context, id string) (model.Order, error) {
	var out order

	oid, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": oid}

	err := r.db.
		Collection("orders").
		FindOne(ctx, query).
		Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Order{}, ErrOrderNotFound
		}
		return model.Order{}, err
	}
	return toOrderModel(out), nil
}

func (r repository) GetProducts(ctx context.Context) ([]*model.Product, error) {
	opt := options.FindOptions{}
	query := bson.M{}

	cursor, err := r.db.
		Collection("products").
		Find(ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var products []*model.Product
	for cursor.Next(ctx) {
		product := &model.Product{}
		// product = toProductModel(product)

		err := cursor.Decode(product)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	return products, nil
}

func (r repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	out, err := r.db.
		Collection("users").
		InsertOne(ctx, fromModel(user))
	if err != nil {
		return model.User{}, err
	}
	user.ID = out.InsertedID.(primitive.ObjectID).String()
	return user, nil
}

func (r repository) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	out, err := r.db.
		Collection("products").
		InsertOne(ctx, fromProductModel(product))
	if err != nil {
		return model.Product{}, err
	}
	product.ID = out.InsertedID.(primitive.ObjectID).String()
	return product, nil
}

func (r repository) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	out, err := r.db.
		Collection("orders").
		InsertOne(ctx, fromOrderModel(order))
	if err != nil {
		return model.Order{}, err
	}
	order.ID = out.InsertedID.(primitive.ObjectID).String()
	return order, nil
}

func (r repository) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	in := bson.M{}
	if user.Name != "" {
		in["name"] = user.Name
	}
	if user.Password != "" {
		in["password"] = user.Password
	}
	out, err := r.db.
		Collection("users").
		UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": in})
	if err != nil {
		return model.User{}, err
	}
	if out.MatchedCount == 0 {
		return model.User{}, ErrUserNotFound
	}
	return user, nil
}

func (r repository) UpdateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	in := bson.M{}
	oid, _ := primitive.ObjectIDFromHex(order.ID)

	if order.ID != "" {
		in["id"] = oid
	}

	if order.Status != "" {
		in["status"] = order.Status
	}

	out, err := r.db.
		Collection("orders").
		UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": in})
	if err != nil {
		return model.Order{}, err
	}
	if out.MatchedCount == 0 {
		return model.Order{}, ErrOrderNotFound
	}
	return order, nil
}

func (r repository) DeleteUser(ctx context.Context, email string) error {
	out, err := r.db.
		Collection("users").
		DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}
	if out.DeletedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

type user struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}

type product struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Price int                `bson:"price,omitempty"`
}

type order struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	OrderNumber string             `bson:"order_number,omitempty"`
	Price       int                `bson:"price,omitempty"`
	Qty         int                `bson:"qty,omitempty"`
	Total       int                `bson:"total,omitempty"`
	CustomerID  string             `bson:"customer_id,omitempty"`
	Status      string             `bson:"status,omitempty"`
	Remark      string             `bson:"remark,omitempty"`
}

func fromModel(in model.User) user {
	return user{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}
}

func fromProductModel(in model.Product) product {
	return product{
		Name:  in.Name,
		Price: in.Price,
	}
}

func fromOrderModel(in model.Order) order {
	return order{
		OrderNumber: in.OrderNumber,
		Price:       in.Price,
		Qty:         in.Qty,
		Total:       in.Total,
		CustomerID:  in.CustomerID,
		Status:      in.Status,
		Remark:      in.Remark,
	}
}

func toModel(in user) model.User {
	return model.User{
		ID:       in.ID.String(),
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	}
}

func toProductModel(in product) model.Product {
	return model.Product{
		ID:    in.ID.String(),
		Name:  in.Name,
		Price: in.Price,
	}
}

func toOrderModel(in order) model.Order {
	return model.Order{
		ID:          in.ID.String(),
		OrderNumber: in.OrderNumber,
		Price:       in.Price,
		Qty:         in.Qty,
		Total:       in.Total,
		CustomerID:  in.CustomerID,
		Status:      in.Status,
		Remark:      in.Remark,
	}
}
