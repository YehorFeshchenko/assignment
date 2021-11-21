package middleware

import (
	"assignment/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type Currency int

const (
	FACTOR       = 1000
	EURString    = "EUR"
	USDString    = "USD"
	DepString    = "deposit"
	WithdrString = "withdrawal"
)

//const URI = "mongodb+srv://admin:admin@cluster0.qmwwt.mongodb.net/assignmentDB?retryWrites=true&w=majority"
var URI string

func StoreTransaction(t models.Transaction) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatalf("Unable to connect to mongoDB: %+v", err)
		return false
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	transactionsColl := client.Database("assignmentDB").Collection("transactions")

	transactionDoc := bson.D{{"user_id", t.UserId}, {"currency", t.CurrencyStr},
		{"amount", t.Amount * FACTOR}, {"time_placed", t.TimePlaced}, {"type", t.OperationTypeStr}}
	_, err = transactionsColl.InsertOne(context.TODO(), transactionDoc)
	if err != nil {
		log.Fatalf("Failed to insert doc: %+v", err)
		return false
	}

	updateAccount(t, client, ctx)

	return true
}

func updateAccount(t models.Transaction, cl *mongo.Client, ct context.Context) {
	ctx, cancel := context.WithTimeout(ct, 20*time.Second)
	defer cancel()
	accountColl := cl.Database("assignmentDB").Collection("accounts")

	filterCursor, err := accountColl.Find(ctx, bson.M{"user_id": t.UserId})
	if err != nil {
		log.Fatal(err)
	}
	var accountsFilter []bson.M
	if err = filterCursor.All(ctx, &accountsFilter); err != nil {
		log.Fatal(err)
	}

	//if result has no docs, simply creating new one
	if accountsFilter == nil {
		accountDoc := createAccountDoc(t)
		_, err := accountColl.InsertOne(ctx, accountDoc)
		if err != nil {
			log.Fatalf("Failed to insert doc: %+v", err)
		}
		return
	} else if len(accountsFilter) > 1 {
		log.Fatalf("Unexpected multiple document")
		return
	}

	//initial document with data for user_id provided
	accountDoc := accountsFilter[0]

	//performing update calculations
	switch t.CurrencyStr {
	case EURString:
		currentBalance := accountDoc["EURBalance"]
		currentAmount := t.Amount
		switch t.OperationTypeStr {
		case DepString:
			accountDoc["EURBalance"] = currentBalance.(float64) + currentAmount*FACTOR
		case WithdrString:
			accountDoc["EURBalance"] = currentBalance.(float64) - currentAmount*FACTOR
		default:
			log.Fatalf("Invalid operation: %+v", t.OperationTypeStr)
			return
		}
	case USDString:
		currentBalance := accountDoc["USDBalance"]
		currentAmount := t.Amount
		switch t.OperationTypeStr {
		case DepString:
			accountDoc["USDBalance"] = currentBalance.(float64) + currentAmount*FACTOR
		case WithdrString:
			accountDoc["USDBalance"] = currentBalance.(float64) - currentAmount*FACTOR
		default:
			log.Fatalf("Invalid operation: %+v", t.OperationTypeStr)
			return
		}
	default:
		log.Fatalf("Invalid currency: %+v", t.CurrencyStr)
		return
	}

	_, err = accountColl.ReplaceOne(ctx, bson.M{"user_id": t.UserId}, accountDoc)
	if err != nil {
		log.Fatal(err)
	}

}

func createAccountDoc(t models.Transaction) bson.D {
	if t.OperationTypeStr != DepString {
		log.Fatalf("Invalid operation: %+v", t.OperationTypeStr)
		return nil
	}
	var accDoc bson.D
	switch t.CurrencyStr {
	case EURString:
		accDoc = bson.D{{"user_id", t.UserId}, {"USDBalance", 0.0}, {"EURBalance", t.Amount * FACTOR}}
	case USDString:
		accDoc = bson.D{{"user_id", t.UserId}, {"USDBalance", t.Amount * FACTOR}, {"EURBalance", 0.0}}
	default:
		log.Fatalf("Invalid currency: %+v", t.CurrencyStr)
		return nil
	}

	return accDoc
}

func GetUserBalance(userId string) interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatalf("Unable to connect to mongoDB: %+v", err)
		return false
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	accountColl := client.Database("assignmentDB").Collection("accounts")

	var result models.Account
	err = accountColl.FindOne(ctx, bson.M{"user_id": userId}).Decode(&result)
	if err != nil {
		log.Fatalf("Failed to find account doc: %+v", err)
	}
	if err == mongo.ErrNoDocuments {
		return err.Error()
	}

	result.USDBalance = result.USDBalance / FACTOR
	result.EURBalance = result.EURBalance / FACTOR

	return result
}
