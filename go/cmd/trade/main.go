package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/GabrielMessiasdaRosa/imersao_fullcycle/go/internal/infra/kafka"
	"github.com/GabrielMessiasdaRosa/imersao_fullcycle/go/internal/market/dto"
	"github.com/GabrielMessiasdaRosa/imersao_fullcycle/go/internal/market/entity"
	"github.com/GabrielMessiasdaRosa/imersao_fullcycle/go/internal/market/transformer"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// precisa de um canal de entra e um de saida para nossas orders
	// ordersIn é o canal que vai receber as orders
	ordersIn := make(chan *entity.Order)

	// ordersOut é o canal que vai enviar as orders
	ordersOut := make(chan *entity.Order)

	// wg é um waitgroup que vai esperar todas as goroutines terminarem
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	// KafkaMsgChan é o canal que vai receber as mensagens do kafka
	KafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}
	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	// cria uma goroutine para consumir as mensagens do kafka
	go kafka.Consume(KafkaMsgChan) // cria uma nova goroutine para consumir as mensagens do kafka que é uma nova thread

	// rcebe do canal kafkaMsgChan e envia para o input do nosso sistema, processa e envia para o output do nosso sistema
	book := entity.NewBook(ordersIn, ordersOut, wg)
	go book.Trade() // go routine é muito bom!

	go func() {
		for msg := range KafkaMsgChan {
			wg.Add(1)
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("orders"), "output")
	}
}
