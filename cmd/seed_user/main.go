package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"katseye/internal/infrastructure/config"
	"katseye/internal/infrastructure/persistence/mongodb"
)

func main() {
	email := flag.String("email", "", "Email do usuário")
	password := flag.String("password", "", "Senha em texto claro a ser criptografada")
	active := flag.Bool("active", true, "Define se o usuário ficará ativo")
	role := flag.String("role", "user", "Papel do usuário (admin, manager, user)")
	profileType := flag.String("profile_type", "service_account", "Tipo de perfil (service_account, partner_manager, consumer)")
	flag.Parse()

	if strings.TrimSpace(*email) == "" || strings.TrimSpace(*password) == "" {
		log.Fatal("email e password são obrigatórios")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("carregando configuração: %v", err)
	}

	emailNorm := strings.TrimSpace(strings.ToLower(*email))
	passHash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("gerando hash da senha: %v", err)
	}

	client, err := mongodb.NewMongoClient(cfg.Mongo.URI)
	if err != nil {
		log.Fatalf("conectando ao mongo: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("erro ao fechar conexão com mongo: %v", err)
		}
	}()

	database := client.Database(cfg.Mongo.Database)
	collection := database.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar se o usuário já existe
	var existingUser bson.M
	err = collection.FindOne(ctx, bson.M{"email": emailNorm}).Decode(&existingUser)
	
	var userID primitive.ObjectID
	if err == nil && existingUser != nil {
		// Usuário existe, usar o ID existente
		userID = existingUser["_id"].(primitive.ObjectID)
	} else {
		// Novo usuário, criar novo ID
		userID = primitive.NewObjectID()
	}

	// Definir permissões padrão baseadas no papel
	permissions := []string{}
	if *role == "admin" {
		permissions = []string{"users:manage"}
	}

	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"_id":           userID,
			"email":         emailNorm,
			"password_hash": string(passHash),
			"active":        *active,
			"role":          *role,
			"permissions":   permissions,
			"profile_type":  *profileType,
			"updated_at":    now,
		},
		"$setOnInsert": bson.M{
			"created_at": now,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{"email": emailNorm}, update, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatalf("inserindo usuário: %v", err)
	}

	switch {
	case result.MatchedCount > 0:
		log.Printf("usuário %s atualizado com sucesso", emailNorm)
	case result.UpsertedCount > 0:
		log.Printf("usuário %s criado com sucesso (id=%v)", emailNorm, userID)
	default:
		log.Printf("nenhuma alteração realizada para o usuário %s", emailNorm)
	}
}
