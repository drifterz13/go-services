package main

func main() {
	conn, err := NewGrpcConn()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	srv := NewServer(conn.GetTaskClient(), conn.GetUserClient())
	srv.Serve()
}
