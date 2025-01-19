package config

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection
var CustomerCollection *mongo.Collection
var EmployeeCollection *mongo.Collection
var ItemCollection *mongo.Collection
var SupplierCollection *mongo.Collection
var TransactionCollection *mongo.Collection
var ReportCollection *mongo.Collection

// InitMongoDB untuk menginisialisasi koneksi ke MongoDB
func InitMongoDB() error {
    uri := "mongodb+srv://karamissuu:karamissu1@cluster0.lyovb.mongodb.net/?retryWrites=true&w=majority"
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB: ", err)
        return err
    }

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal("MongoDB Ping failed: ", err)
        return err
    }

    log.Println("MongoDB connected successfully")
    Client = client
    UserCollection = Client.Database("apkclaundry").Collection("user")
    CustomerCollection = Client.Database("apkclaundry").Collection("pelanggan")
    EmployeeCollection = Client.Database("apkclaundry").Collection("karyawan")
	ItemCollection = client.Database("apkclaundry").Collection("barang")
	SupplierCollection = client.Database("apkclaundry").Collection("supplier")
	TransactionCollection = client.Database("apkclaundry").Collection("transaksi")
	ReportCollection = client.Database("apkclaundry").Collection("laporan")

    return nil
}
