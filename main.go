package main

func main() {
	server := NewAPIServer(":8008")
	server.Run()
}
