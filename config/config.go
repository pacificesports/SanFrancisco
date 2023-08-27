package config

import (
	"os"
	"san_francisco/model"
	"strings"
)

var Service = model.Service{
	Name:        os.Getenv("SERVICE_NAME"),
	StatusEmail: os.Getenv("STATUS_EMAIL"),
	URL:         "http://" + strings.ToLower(os.Getenv("SERVICE_NAME")) + ":" + Port,
	Version:     Version,
}

var RinconService = model.Service{}

var Version = "2.3.0"
var Env = os.Getenv("ENV")
var Port = os.Getenv("PORT")
var RinconPort = os.Getenv("RINCON_PORT")
var JaegerPort = os.Getenv("JAEGER_PORT")

var PostgresHost = os.Getenv("POSTGRES_HOST")
var PostgresDatabase = os.Getenv("POSTGRES_DATABASE")
var PostgresUser = os.Getenv("POSTGRES_USER")
var PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
var PostgresPort = os.Getenv("POSTGRES_PORT")

var DiscordToken = os.Getenv("DISCORD_TOKEN")
var DiscordGuild = os.Getenv("DISCORD_GUILD")
var DiscordChannel = os.Getenv("DISCORD_CHANNEL")

var FirebaseServiceAccountEncoded = os.Getenv("FIREBASE_SERVICE_ACCOUNT")

var APIKeys []model.APIKey
