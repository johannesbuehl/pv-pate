package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

type ConfigYaml struct {
	LogLevel string `yaml:"log_level"`
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"database"`
	Cache struct {
		Expiration string `yaml:"expiration"`
		Purge      string `yaml:"purge"`
	} `yaml:"cache"`
	ClientSession struct {
		JwtSignature string `yaml:"jwt_signature"`
		Expire       string `yaml:"expire"`
	} `yaml:"client_session"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Reservation struct {
		Expiration string `yaml:"expiration"`
	} `yaml:"reservation"`
	Mail struct {
		Server   string `yaml:"server"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Template struct {
			Subject   string `json:"body"`
			Body      string `yaml:"body"`
			BodyPlain string `yaml:"body_plain"`
		} `json:"template"`
	} `yaml:"mail"`
}

type CacheConfig struct {
	Expiration time.Duration
	Purge      time.Duration
}

type ReservationConfig struct {
	Expiration time.Duration
}

type ConfigStruct struct {
	ConfigYaml
	LogLevel      zerolog.Level
	SessionExpire time.Duration
	Cache         CacheConfig
	Reservation   ReservationConfig
	Templates     ConfigTemplates
}

type ConfigTemplates struct {
	Subject   *template.Template
	Body      *template.Template
	BodyPlain *template.Template
}

var config ConfigStruct

var logger zerolog.Logger

type specificLevelWriter struct {
	io.Writer
	Level zerolog.Level
}

func (w specificLevelWriter) WriteLevel(l zerolog.Level, p []byte) (int, error) {
	if l >= w.Level {
		return w.Write(p)
	} else {
		return len(p), nil
	}
}

type Payload struct {
	jwt.RegisteredClaims
	CustomClaims map[string]any
}

func (config ConfigStruct) signJWT(val any) (string, error) {
	valMap, err := strucToMap(val)

	if err != nil {
		return "", err
	}

	payload := Payload{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.SessionExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		CustomClaims: valMap,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return t.SignedString([]byte(config.ClientSession.JwtSignature))
}

func loadConfig() ConfigStruct {
	config := ConfigYaml{}

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Panic().Msgf("Error opening config-file: %q", err)
	}

	reader := bytes.NewReader(yamlFile)

	dec := yaml.NewDecoder(reader)
	dec.KnownFields(true)
	err = dec.Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config-file: %v", err)
		os.Exit(1)
	}

	if logLevel, err := zerolog.ParseLevel(config.LogLevel); err != nil {
		panic(fmt.Errorf("can't parse log-level: %v", err))
	} else {
		var configStruct ConfigStruct

		// parse the durations
		if session_expire, err := time.ParseDuration(config.ClientSession.Expire); err != nil {
			fmt.Fprintf(os.Stderr, `Error parsing "client_session.expire": %v`, err)
			os.Exit(1)
		} else if cacheExpire, err := time.ParseDuration(config.Cache.Expiration); err != nil {
			fmt.Fprintf(os.Stderr, `Error parsing "cache.expiration": %v`, err)
			os.Exit(1)
		} else if cachePurge, err := time.ParseDuration(config.Cache.Purge); err != nil {
			fmt.Fprintf(os.Stderr, `Error parsing "cache.purge": %v`, err)
			os.Exit(1)

		} else if reservationExpire, err := time.ParseDuration(config.Reservation.Expiration); err != nil {
			fmt.Fprintf(os.Stderr, `Error parsing "reservation.expiration": %v`, err)
			os.Exit(1)

			// parse the templates
		} else if mailSubjectTemplate, err := template.New("mailSubject").Parse(config.Mail.Template.Subject); err != nil {
			fmt.Fprintf(os.Stderr, `Error parsing "mail.template.subject": %v`, err)

			os.Exit(1)
		} else if mailBodyTemplate, err := template.New("mailBody").Parse(config.Mail.Template.Body); err != nil {
			fmt.Fprintf(os.Stderr, `Error parsing "mail.template.body": %v`, err)

			os.Exit(1)
		} else if mailBodyPlainTemplate, err := template.New("mailBodyPlain").Parse(config.Mail.Template.BodyPlain); err != nil {
			fmt.Fprintf(os.Stderr, `Error parsing "mail.template.body_plain": %v`, err)

			os.Exit(1)
		} else {
			configStruct = ConfigStruct{
				ConfigYaml:    config,
				LogLevel:      logLevel,
				SessionExpire: session_expire,
				Cache: CacheConfig{
					Expiration: cacheExpire,
					Purge:      cachePurge,
				},
				Reservation: ReservationConfig{
					Expiration: reservationExpire,
				},
				Templates: ConfigTemplates{
					Subject:   mailSubjectTemplate,
					Body:      mailBodyTemplate,
					BodyPlain: mailBodyPlainTemplate,
				},
			}
		}

		return configStruct
	}
}

func init() {
	config = loadConfig()

	// try to set the log-level
	zerolog.SetGlobalLevel(config.LogLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// create the console output
	outputConsole := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		NoColor: true,
	}

	// create the logfile output
	outputLog := &lumberjack.Logger{
		Filename:  "logs/backend.log",
		MaxAge:    7,
		LocalTime: true,
	}

	// create a multi-output-writer
	multi := zerolog.MultiLevelWriter(
		specificLevelWriter{
			Writer: outputConsole,
			Level:  config.LogLevel,
		},
		specificLevelWriter{
			Writer: outputLog,
			Level:  config.LogLevel,
		},
	)

	// create a logger-instance
	logger = zerolog.New(multi).With().Timestamp().Logger()
}
