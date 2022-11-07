package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "storage/grpc-storage"
	"sync"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStorageManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	str := "hi"

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := c.Set(ctx, &pb.SetRequest{Key: "admin", Value: []byte(str)})
			if err != nil {
				log.Fatalf("set data failed: %v", err)
				return
			}
			log.Printf("SET: %v", s)

		}()
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g, err := c.Get(ctx, &pb.GetRequest{Key: "admin"})
			if err != nil {
				log.Fatalf("didn't receive anything: %v", err)
				return
			}

			log.Printf("GOT: %v", g)
		}()
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d, err := c.Delete(ctx, &pb.DeleteRequest{Key: "admin"})
			if err != nil {
				log.Fatalf("didn't delete anything: %v", err)
				return
			}

			log.Printf("DElETED: %v", string(d.ResultDeleted))
		}()
	}
	wg.Wait()

	//g, err := c.Get(ctx, &pb.GetRequest{Key: "admin"})
	//if err != nil {
	//	log.Fatalf("didn't receive anything: %v", err)
	//	return
	//}

	//log.Printf("GOT: %v", g)
	//
	//d, err := c.Delete(ctx, &pb.DeleteRequest{Key: "admin"})
	//if err != nil {
	//	log.Fatalf("didn't delete anything: %v", err)
	//	return
	//}

	//log.Printf("DElETED: %v", string(d.ResultDeleted))

}
