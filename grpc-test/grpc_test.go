package grpc_test

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	ser "storage/grpc-server"
	pb "storage/grpc-storage"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
)

// Testing Memcache
func TestSet(t *testing.T) {
	var m = &ser.StorageManagementServer{
		Storage: make(map[string][]byte),
	}

	var ctx, _ = context.WithTimeout(context.Background(), time.Second*1000)

	setreq := &pb.SetRequest{
		Key:   "admin",
		Value: []byte{7},
	}
	Convey("Testing All Memcache methods", t, func() {
		valSet, _ := m.Set(ctx, setreq)

		getreq := &pb.GetRequest{
			Key: valSet.Key,
		}
		g, _ := m.Get(ctx, getreq)

		Convey("Keys from set() and get() functions are equal", func() {
			So(g.Key, ShouldEqual, valSet.Key)
		})

		deletereq := &pb.DeleteRequest{
			Key: valSet.Key,
		}

		_, err := m.Delete(ctx, deletereq)
		Convey("Delete function has no errors ", func() {
			So(err, ShouldEqual, nil)
		})

		log.Printf("Key from get(): %v", getreq.Key)
		log.Printf("Key from set(): %v", setreq.Key)
		log.Printf("Key from delete(): %v", deletereq.Key)

	})
}

// / Testing GRPC server
func TestGRPC(t *testing.T) {
	Convey("Testing connection", t, func() {
		lis, err := net.Listen("tcp", ":9997")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		So(err, ShouldEqual, nil)

		wg := &sync.WaitGroup{}
		wg.Add(2)
		s := grpc.NewServer()
		var sr = &ser.StorageManagementServer{
			Grpc: s,
		}

		Convey("Testing launch of the server", func(c C) {
			pb.RegisterStorageManagementServer(s, sr)
			log.Printf("server listening at %v", lis.Addr())
			go func(wg *sync.WaitGroup) {
				log.Printf("before turning on to the server")
				err := s.Serve(lis)
				if err != nil {
					log.Printf("connection error: %v", err)
				}

				log.Printf("here we have already stopped server")
				wg.Done()

				c.So(err, ShouldEqual, nil)
			}(wg)
			Convey("Testing stop of the server", func(c C) {
				go func(sr *ser.StorageManagementServer, wg *sync.WaitGroup) {
					log.Printf("awaiting before stop")
					time.Sleep(time.Second * 3)
					log.Printf("stopping server")
					sr.StopGrpcServer()
					if err != nil {
						log.Printf("connection error: %v", err)
					}
					log.Printf("after stop of the server")
					wg.Done()
					c.So(err, ShouldEqual, nil)
				}(sr, wg)
				wg.Wait()
			})
		})
	})
}
