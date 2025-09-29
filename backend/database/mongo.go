package database

//ESTOS METODOS FUNCIONAN PERO NO ESCALAN
//LO SACAMOS A DB.go

// var Client *mongo.Client

// func Connect() error {
// 	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

// 	client, err := mongo.Connect(context.Background(), clientOptions)

// 	if err != nil {
// 		return err
// 	}

// 	err = client.Ping(context.Background(), nil)
// 	if err != nil {
// 		return err
// 	}

// 	Client = client
// 	fmt.Printf("Client Connect: %v\n", &Client)

// 	return nil
// }

// func Disconnect() error {
// 	fmt.Printf("Client Disconnect: %v\n", &Client)

// 	return Client.Disconnect(context.Background())
// }
