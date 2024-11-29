package inferer

import (
	"bufio"
	"daily-dashboard-backend/src/data"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

func (e Endpoint) SendMessage(convo *data.Conversation, w *http.ResponseWriter) string {
	fmt.Println("Sending our conversation over: ", *convo)

	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", e.Host, e.Port))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer conn.Close()

	// Convert Conversation to Byte Array
	convoBytes, err := json.Marshal(convo)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}
	convoBytesStr := string(convoBytes) + "<sobadd>"

	// Write to the LLM - Currently I only send over the latest message, no full context
	_, err = conn.Write([]byte(convoBytesStr))
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	// Initialize flusher to flush response back to frontend subsequently
	flusher, ok := (*w).(http.Flusher)
	if !ok {
		fmt.Println("Flusher error: ", err)
		return ""
	}

	// Prepare reader to read from our TCP Connection
	reader := bufio.NewReader(conn)
	buffer := make([]byte, 1024)

	// Read from Stream
	result := ""
	for {
		nBytes, err := reader.Read(buffer)

		// Completed / Error handling
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Reading error: ", err)
			break
		}

		// Add Response Chunk to Result - So that we can append our response to the DB later
		result += string(buffer[:nBytes])

		// Write Response Chunk back to Frontend
		_, err = fmt.Fprintf(*w, "%s", buffer[:nBytes])
		flusher.Flush()

		if err != nil {
			fmt.Println("Response writing error: ", err)
			break
		}
	}

	return result
}
