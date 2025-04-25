package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultPort string = ":6379"
const saveFileName string = "data.json"

type Config struct {
	Port string
}
type Server struct {
	Config   Config
	Listener net.Listener
	Database map[string]string
}

func NewServer(config Config) *Server {
	if len(config.Port) == 0 {
		config.Port = defaultPort
	}
	return &Server{Config: config, Database: LoadDatabase()}
}

func (server *Server) Listen() error {
	slog.Info("Starting GoRedis Server at port " + server.Config.Port)
	listener, err := net.Listen("tcp", server.Config.Port)
	if err != nil {
		return err
	}
	server.Listener = listener

	for {
		connection, err := server.Listener.Accept()
		if err != nil {
			return err
		}
		go server.HandleConnection(connection)
	}
}

func (server *Server) HandleConnection(connection net.Conn) {
	slog.Info("Connection established.")
	defer connection.Close()
	buffer := bufio.NewReader(connection)

	for {
		data, err := buffer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				slog.Info("Client disconnected.")
			} else {
				slog.Error("Error in reading data.")
			}
			return
		} else if len(strings.TrimSpace(data)) == 0 {
			continue
		}
		go server.ExecuteCommand(connection, data)
	}
}

func LoadDatabase() map[string]string {
	var database map[string]string
	file, err := os.Open(saveFileName)
	if err != nil {
		return map[string]string{}
	}

	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Database unreadable.")
	}

	json.Unmarshal(data, &database)
	return database
}

func SaveDatabase(database map[string]string) error {
	data, err := json.Marshal(database)
	if err != nil {
		slog.Error("Error serializing data.")
		return err
	}

	err = os.WriteFile(saveFileName, data, 0644)
	if err != nil {
		slog.Error("Error saving database.")
		return err
	}

	return nil
}

func (server *Server) ExecuteCommand(connection net.Conn, data string) {
	var response string
	data = strings.TrimSpace(data)
	dataSplit := strings.Split(data, " ")
	command := dataSplit[0]

	if strings.EqualFold(data, "PING") {
		response = SerializeIntoRESP("PONG")
	} else if strings.EqualFold(data, "SAVE") {
		err := SaveDatabase(server.Database)
		if err != nil {
			response = SerializeIntoRESP(err.Error())
		}
		response = SerializeIntoRESP("OK")
	} else if strings.EqualFold(command, "GET") {
		key := dataSplit[1]
		response = SerializeIntoRESP(server.Database[key])
		if len(response) == 0 {
			response = "KEY NOT FOUND"
		}
	} else if strings.EqualFold(command, "SET") {
		key := dataSplit[1]
		value := dataSplit[2]
		if len(dataSplit) == 5 {
			expiryFlag := dataSplit[3]
			expiryTime := dataSplit[4]
			go server.ExpireKeys(key, expiryFlag, expiryTime)
		}
		server.Database[key] = value
		response = SerializeIntoRESP("OK")
	} else if strings.EqualFold(command, "DELETE") {
		key := dataSplit[1]
		delete(server.Database, key)
		response = SerializeIntoRESP("OK")
	} else if strings.EqualFold(command, "EXISTS") {
		key := dataSplit[1]
		if len(server.Database[key]) == 0 {
			response = SerializeIntoRESP("NO")
		} else {
			response = SerializeIntoRESP("YES")
		}
	} else if strings.EqualFold(command, "INCR") {
		key := dataSplit[1]
		if len(server.Database[key]) == 0 {
			server.Database[key] = "1"
			response = SerializeIntoRESP("OK")
		} else {
			number, err := strconv.Atoi(server.Database[key])
			if err != nil {
				response = SerializeIntoRESP("NON-INTEGER VALUE")
			} else {
				server.Database[key] = strconv.Itoa(number + 1)
				response = SerializeIntoRESP("OK")
			}
		}
	} else if strings.EqualFold(command, "DECR") {
		key := dataSplit[1]
		if len(server.Database[key]) == 0 {
			server.Database[key] = "-1"
			response = SerializeIntoRESP("OK")
		} else {
			number, err := strconv.Atoi(server.Database[key])
			if err != nil {
				response = SerializeIntoRESP("NON-INTEGER VALUE")
			} else {
				server.Database[key] = strconv.Itoa(number - 1)
				response = SerializeIntoRESP("OK")
			}
		}
	}

	_, err := connection.Write([]byte(response + "\n"))
	if err != nil {
		slog.Error("Error in responding after evaluation of command.")
	}
	slog.Info("Responded with " + response)
}

func (server *Server) ExpireKeys(key string, expiryFlag string, expiryTime string) {
	expiration, err := strconv.Atoi(expiryTime)
	if err != nil {
		slog.Error("Non-integer value passed in for expiry time.")
		return
	}
	slog.Info("Data expiry task queued for key: " + key)

	if strings.EqualFold(expiryFlag, "EX") {
		time.Sleep(time.Second * time.Duration(expiration))
	} else if strings.EqualFold(expiryFlag, "PX") {
		time.Sleep(time.Millisecond * time.Duration(expiration))
	}

	delete(server.Database, key)
	slog.Info("The deleted key is: " + key)
}

func SerializeIntoRESP(data string) string {
	return data
}

func main() {
	serverConfig := Config{}
	server := NewServer(serverConfig)
	server.Listen()
}
