package main

import (
	"encoding/json"
	"log"
	"os"
	"sudoku-server/sudoku"

	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

const (
	addr         = ":8080"
	RoomString   = "room"
	PlayerString = "player"
)

func main() {

	f, err := os.OpenFile("logs", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer f.Close()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	m := melody.New()
	playerCh := make(chan *Player)

	go func() {
		var room *Room
		for {

			// Создаем новую комнату
			room = NewRoom()

			// Заполняем её двумя игроками. Назначаем комнату игроку
			for i := range room.players {
				room.players[i] = <-playerCh
				room.players[i].Session.Set(RoomString, room)
			}

			for i := range room.players {
				room.players[i].Puzzle = sudoku.CopyGrid(room.Puzzle)
			}

			fm := &FirstMessage{
				Puzzle: room.Puzzle,
				Time:   gameTime,
			}

			jsonData, _ := json.Marshal(fm)
			for i := range room.players {
				room.players[i].Session.Write(jsonData)
			}

		}
	}()

	m.HandleConnect(func(s *melody.Session) {

		// Создаем нового игрока и передаем в канал для послания первого сообщения
		player := NewPlayer(s)

		logger.Println(LogNewPlayer)

		s.Set(PlayerString, player)
		go func() { playerCh <- player }()
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		// Нахохим игрока по его адресу
		player := s.MustGet(PlayerString).(*Player)

		// Нахохим комнату по roomId, что есть у игрока
		room := s.MustGet(RoomString).(*Room)

		// Ответ от клиента
		msgDTO := &MessageDTO{}
		json.Unmarshal(msg, msgDTO)

		// Завершаем соединение если это конец игры
		if msgDTO.IsEnd {
			room.Mu.Lock()
			if room.Closed {
				room.Mu.Unlock()

				log.Print(LogCloseRoom)
				logger.Writer().Write([]byte(LogCloseRoom))

				return
			}
			room.Closed = true

			players := room.players
			room.Mu.Unlock()

			for _, p := range players {
				p.Session.Close()
			}
			return
		}

		// Назнаначем игроку обновленный пазл
		player.Mu.Lock()
		player.Puzzle = msgDTO.Puzzle
		puzzleCopy := sudoku.CopyGrid(player.Puzzle)
		player.Mu.Unlock()

		// Проверяем решено ли судоку у игрока
		solved := sudoku.IsSolved(puzzleCopy, room.Solution)

		// Индексы ошибки. Будут равны -1 -1 если ошибок нету
		valid := sudoku.ValidAnswer(puzzleCopy, room.Solution)

		// Второй игрок нужен чтобы клиент мог высчитать сколько его противник заполнил клеток
		var secondPlayer *Player

		// Находим второго игрока
		for _, p := range room.players {
			if p != player {
				secondPlayer = p
				break
			}
		}

		// Если второго игрока нету то завершаем
		if secondPlayer == nil {
			return
		}

		// Получаем пазл второго игроков
		secondPlayer.Mu.Lock()
		secondPuzzle := sudoku.CopyGrid(secondPlayer.Puzzle)
		secondPlayer.Mu.Unlock()

		// Заполняем структуру сообщения
		sendmsgDTO := &SendMessageDTO{
			IsValid:  valid,
			IsSolved: solved,
			Puzzle:   secondPuzzle,
		}

		// Отправляем данные клиенту
		jsonData, _ := json.Marshal(sendmsgDTO)
		player.Session.Write(jsonData)

	})

	m.HandleDisconnect(func(s *melody.Session) {
		// Пытаемся получить игрока
		p, ok := s.Get(PlayerString)
		if !ok {
			return
		}

		// Делаем привидение к игроку
		player := p.(*Player)

		// Пытаемся получить комнату
		r, ok := s.Get(RoomString)
		if !ok {
			return
		}

		// Делаем привидение к комнате
		room := r.(*Room)

		/* Закрываем комнату если открыта, иначе пропускаем так как соединения
		уже закрыты*/
		room.Mu.Lock()
		if room.Closed {
			room.Mu.Unlock()
			return
		}
		room.Closed = true

		players := room.players
		room.Mu.Unlock()

		// Закрыаем все соединения
		for _, p := range players {
			if p != player {
				p.Session.Close()
			}
		}
	})

	r.GET("/ws", func(ctx *gin.Context) {
		m.HandleRequest(ctx.Writer, ctx.Request)
	})

	logger.Println(LogStartServer)
	r.Run(addr)
}
