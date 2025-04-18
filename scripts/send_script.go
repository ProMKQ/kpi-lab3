package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// 5 кроків для руху фігури по діагоналі
	for i := 0; i < 5; i++ {
		// координати зміщуються по діагоналі
		x := 0.5 + float64(i)*0.03
		y := 0.5 + float64(i)*0.03

		// команда, яка буде відправлена серверу
		script := fmt.Sprintf(`reset
white
figure %.2f %.2f
update`, x, y)

		// HTTP POST-запит на сервер
		resp, err := http.Post("http://localhost:17000/", "text/plain", bytes.NewBufferString(script))
		if err != nil {
			fmt.Println("Помилка запиту:", err)
			continue
		}
		resp.Body.Close()

		// затримка 1 секунда між кадрами
		time.Sleep(1 * time.Second)
	}
}
