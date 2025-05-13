LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/note_v1/note.proto

gen-cert:
	mkdir -p certificates
	# Генерация закрытого ключа для центра сертификации (CA)
	openssl genrsa -out ca.key 4096
	# Создание самоподписанного сертификата CA на основе ключа (срок действия 365 дней)
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
	# Генерация закрытого ключа для gRPC-сервиса
	openssl genrsa -out service.key 4096
	# Создание запроса на сертификат (CSR) для сервиса с использованием конфигурационного файла
	openssl req -new -key service.key -out service.csr -config certificate.conf
	# Подписание CSR с помощью CA, чтобы выпустить действительный сертификат для сервиса
	# Используется конфигурационный файл для указания расширений сертификата
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
    		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext