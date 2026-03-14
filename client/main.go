package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	serverAddr = ":8080"
	size       = 9
)

func main() {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer conn.Close()

	var puzzle [][]int
	err = json.NewDecoder(conn).Decode(&puzzle)
	if err != nil {
		log.Fatal("Failed to receive puzzle:", err)
	}

	fmt.Println("Sudoku puzzle received:")
	printBoard(puzzle)

	reader := bufio.NewReader(os.Stdin)

	for {
		board := puzzle
		fmt.Println("Enter your move in format: row col value (1-based indices)")
		fmt.Print("> ")
		var row, col, val int
		_, err := fmt.Fscanf(reader, "%d %d %d\n", &row, &col, &val)
		if err != nil {
			fmt.Println("Invalid input, try again")
			reader.ReadString('\n')
			continue
		}

		board[row-1][col-1] = val

		err = json.NewEncoder(conn).Encode(board)
		if err != nil {
			log.Println("Failed to send board:", err)
			break
		}

		var resp map[string]any
		err = json.NewDecoder(conn).Decode(&resp)
		if err != nil {
			log.Println("Failed to receive response:", err)
			break
		}

		fmt.Println("Server response:")
		if lives, ok := resp["lives"]; ok {
			fmt.Println("Lives remaining:", lives)
		}
		if errors, ok := resp["errors"]; ok {
			fmt.Println("Mistakes in this move:", errors)
		}
		if updatedBoard, ok := resp["board"]; ok {
			fmt.Println("Updated board:")
			boardInterface := updatedBoard.([]any)
			for i := range size {
				rowSlice := boardInterface[i].([]any)
				for j := range size {
					fmt.Printf("%d ", int(rowSlice[j].(float64)))
				}
				fmt.Println()
			}
		}
	}
}

func printBoard(board [][]int) {
	for i := range size {
		for j := range size {
			fmt.Printf("%d ", board[i][j])
		}
		fmt.Println()
	}
}
