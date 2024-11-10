package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
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
		Server    string `yaml:"server"`
		Port      int    `yaml:"port"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		Templates struct {
			ReservationSubject string `yaml:"reservation_subject"`
			CertificateSubject string `yaml:"certificate_subject"`
		} `yaml:"subject_templates"`
	} `yaml:"mail"`
	ValidateElements struct {
		Regex         string `yaml:"regex"`
		ValidElements map[string]struct {
			From int `yaml:"from"`
			To   int `yaml:"to"`
		} `yaml:"valid_elements"`
	} `yaml:"validate_elements"`
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
	MidRegex      *regexp.Regexp
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
			log.Fatalf(`Error parsing "client_session.expire": %v`, err)
		} else if cacheExpire, err := time.ParseDuration(config.Cache.Expiration); err != nil {
			log.Fatalf(`Error parsing "cache.expiration": %v`, err)
		} else if cachePurge, err := time.ParseDuration(config.Cache.Purge); err != nil {
			log.Fatalf(`Error parsing "cache.purge": %v`, err)
		} else if reservationExpire, err := time.ParseDuration(config.Reservation.Expiration); err != nil {
			log.Fatalf(`Error parsing "reservation.expiration": %v`, err)

			// parse the templates
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
				MidRegex: regexp.MustCompile(config.ValidateElements.Regex),
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
