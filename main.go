package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	fj "github.com/s-vvardenfell/CloudKeyValStorage/file_journal"
	"github.com/s-vvardenfell/CloudKeyValStorage/storage"
)

var logger fj.TransactionLogger

func initializeTransactionLog() error {
	var err error
	logger, err = fj.NewFileTransactionLogger("logs/transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := fj.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors: // Получает ошибки
		case e, ok = <-events:
			switch e.EventType {
			case fj.EventDelete:
				// Получено событие DELETE!
				err = storage.Delete(e.Key)
			case fj.EventPut:
				// Получено событие PUT!
				err = storage.Put(e.Key, e.Value)
			}
		}
	}
	logger.Run()
	return err
}

// keyValuePutHandler ожидает получить PUT-запрос с
// ресурсом "/v1/key/{key}".
func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Получить ключ из запроса
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	// Тело запроса хранит значение
	if err != nil {
		// Если возникла ошибка, сообщить о ней
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = storage.Put(key, string(value))
	// Сохранить значение как строку
	if err != nil {
		// Если возникла ошибка, сообщить о ней
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	logger.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated) // Все хорошо! Вернуть StatusCreated
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Извлечь ключ из запроса
	key := vars["key"]
	value, err := storage.Get(key)
	// Получить значение для данного ключа
	if errors.Is(err, storage.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(value))
	// Записать значение в ответ
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Извлечь ключ из запроса
	key := vars["key"]

	err := storage.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}

func main() {
	initializeTransactionLog()

	r := mux.NewRouter()
	// Зарегистрировать keyValuePutHandler как обработчик HTTP-запросов PUT,
	// в которых указан путь "/v1/{key}"
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
