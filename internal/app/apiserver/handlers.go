package apiserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"DZ2/internal/app/middleware"
	"DZ2/internal/app/models"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gorilla/mux"
)

type Error struct {
	Error string `json:"error"`
}

type Message struct {
	Message string `json:"message"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func (api *APIServer) GetStock(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	articles, err := api.store.Article().SelectAll()
	if err != nil {
		api.logger.Info(err)
		msg := Error{
			Error: "We have some troubles to accessing articles in database. Try later",
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if len(articles) == 0 {
		writer.WriteHeader(400)
		msg := Error{
			Error: "No autos found in DataBase",
		}
		json.NewEncoder(writer).Encode(msg)

	} else {
		api.logger.Info("Get All Cars GET /stock")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(articles)
	}
}

func (api *APIServer) PostCar(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post Car POST /auto/{mark}")
	mark := mux.Vars(req)["mark"]
	var car models.Car
	err := json.NewDecoder(req.Body).Decode(&car)
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Error{
			Error: "Provided json is invalid",
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	car.Mark = mark
	_, created, err := api.store.Article().Create(&car)
	if err != nil {
		api.logger.Info("Troubles while creating new car:", err)
		msg := Error{
			Error: "We have some troubles to accessing database. Try again",
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	} else if !created {
		writer.WriteHeader(400)
		msg := Error{
			Error: "Auto with that mark already exists",
		}
		json.NewEncoder(writer).Encode(msg)
	} else {
		writer.WriteHeader(201)
		msg := Message{
			Message: "Auto created",
		}
		json.NewEncoder(writer).Encode(msg)
	}
}

func (api *APIServer) PutCar(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post Car PUT /auto/{mark}")
	mark := mux.Vars(req)["mark"]
	var car models.Car
	err := json.NewDecoder(req.Body).Decode(&car)
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Error{
			Error: "Provided json is invalid",
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	car.Mark = mark
	_, updated, err := api.store.Article().Update(&car)
	if err != nil {
		api.logger.Info("Troubles while updating a car:", err)
		msg := Error{
			Error: "We have some troubles to accessing database. Try again",
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return

	} else if !updated {
		writer.WriteHeader(404)
		msg := Error{
			Error: "Auto with that mark not found",
		}
		json.NewEncoder(writer).Encode(msg)
	} else {
		writer.WriteHeader(202)
		msg := Message{
			Message: "Auto updated",
		}
		json.NewEncoder(writer).Encode(msg)
	}
}

func (api *APIServer) GetArticleById(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get Car by ID /articles/{id}")
	id := mux.Vars(req)["mark"]
	article, ok, err := api.store.Article().FindCarById(id)

	log.Println(article, ok, err)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (articles) with id. err:", err)
		msg := Error{
			Error: "We have some troubles to accessing database. Try again",
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	if !ok {
		api.logger.Info("Can not find article with that ID in database")
		msg := Error{
			Error: "Car with that ID does not exists in database.",
		}

		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(202)
	json.NewEncoder(writer).Encode(article)

}

func (api *APIServer) DeleteArticleById(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Delete Car by Id DELETE /articles/{mark}")
	mark := mux.Vars(req)["mark"]

	_, ok, err := api.store.Article().FindCarById(mark)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (articles) with mark. err:", err)
		msg := Error{
			Error: "We have some troubles to accessing database. Try again",
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	if !ok {
		api.logger.Info("Can not find article with that ID in database")
		msg := Error{
			Error: "Car with that ID does not exists in database.",
		}

		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	_, deleted, err := api.store.Article().DeleteById(mark)
	if err != nil {
		api.logger.Info("Troubles while deleting database element from table (articles) with mark. err:", err)
		msg := Error{
			Error: "We have some troubles to accessing database. Try again",
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	} else if !deleted {
		msg := Error{
			Error: "Auto with that mark not found",
		}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
	} else {
		writer.WriteHeader(202)
		msg := Error{
			Error: fmt.Sprintf("Auto deleted"),
		}
		json.NewEncoder(writer).Encode(msg)
	}
}

func (api *APIServer) PostUserRegister(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post User Register POST /register")
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Error{
			Error: "Provided json is invalid",
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Пытаемся найти пользователя с таким логином в бд
	_, ok, err := api.store.User().FindByLogin(user.Login)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
		msg := Error{
			Error: "We have some troubles to accessing database. Try again",
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Смотрим, если такой пользователь уже есть - то никакой регистрации мы не делаем!
	if ok {
		api.logger.Info("User with that ID already exists")
		msg := Error{
			Error: "User with that login already exists in database",
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Теперь пытаемся добавить в бд
	_, err = api.store.User().Create(&user)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
		msg := Error{
			Error: "We have some troubles to accessing database. Try again",
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	msg := Message{
		Message: fmt.Sprint("User created. Try to auth"),
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)

}

func (api *APIServer) PostToAuth(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post to Auth POST /auth")
	var userFromJson models.User
	err := json.NewDecoder(req.Body).Decode(&userFromJson)
	//Обрабатываем случай, если json - вовсе не json или в нем какие-либо пробелмы
	if err != nil {
		api.logger.Info("Invalid json recieved from client")
		msg := Error{
			Error: "Provided json is invalid",
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Необходимо попытаться обнаружить пользователя с таким login в бд
	userInDB, ok, err := api.store.User().FindByLogin(userFromJson.Login)
	// Проблема доступа к бд
	if err != nil {
		api.logger.Info("Can not make user search in database:", err)
		msg := Error{
			Error: "We have some troubles while accessing database",
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Если подключение удалось , но пользователя с таким логином нет
	if !ok {
		api.logger.Info("User with that login does not exists")
		msg := Error{
			Error: "User with that login does not exists in database. Try register first",
		}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//Если пользователь с таким логином ест ьв бд - проверим, что у него пароль совпадает с фактическим
	if userInDB.Password != userFromJson.Password {
		api.logger.Info("Invalid credetials to auth")
		msg := Error{
			Error: "Your password is invalid",
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	//Теперь выбиваем токен как знак успешной аутентифкации
	token := jwt.New(jwt.SigningMethodHS256)             // Тот же метод подписания токена, что и в JwtMiddleware.go
	claims := token.Claims.(jwt.MapClaims)               // Дополнительные действия (в формате мапы) для шифрования
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() //Время жизни токена
	claims["admin"] = true
	claims["name"] = userInDB.Login
	tokenString, err := token.SignedString(middleware.SecretKey)
	//В случае, если токен выбить не удалось!
	if err != nil {
		api.logger.Info("Can not claim jwt-token")
		msg := Error{
			Error: "We have some troubles. Try again",
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	//В случае, если токен успешно выбит - отдаем его клиенту
	msg := Message{
		Message: tokenString,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)

}
