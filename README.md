# imersao_fullcycle

### estrutura de pastas:

go - projeto golang

- cmd
  - trade
    - main.go
- internal
  - infra
    - kafka
      - consumer.go
      - producer.go
  - market
    - dto
      - order-output.dto.go
      - trade-input.dto.go
      - transaction-output.dto.go
    - entity
      - asset.go
      - book_test.go
      - book.go
      - investor.go
      - order_queue.go
      - order.go
      - transaction.go
    - tranformer
      - transform-input.go
      - transform-output.go
