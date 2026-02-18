package test

import (
	"context"
	"log"
	"net"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Logstay/tech-challange-execution-service/config"
	app "github.com/Logstay/tech-challange-execution-service/internal"
	httpServer "github.com/Logstay/tech-challange-execution-service/internal/adapters/http"
	"github.com/Logstay/tech-challange-execution-service/pkg/cache"
	"github.com/Logstay/tech-challange-execution-service/pkg/constant"
	database "github.com/Logstay/tech-challange-execution-service/pkg/postgres"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"
)

var (
	dbFactory database.ConnectionFactory
	testCfg   *config.Config
	testURL   *url.URL
)

type ErrorResponse struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"Descrição do erro"`
}

const jwtExpirationTime = 60 * time.Minute

func TestMain(m *testing.M) {
	tearDown := setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func CleanAll() {
	_, _ = dbFactory.GetDB().Exec(`TRUNCATE TABLE consumable RESTART IDENTITY CASCADE`)
	_, _ = dbFactory.GetDB().Exec(`TRUNCATE TABLE inventory RESTART IDENTITY CASCADE`)
	_, _ = dbFactory.GetDB().Exec(`TRUNCATE TABLE service_definition RESTART IDENTITY CASCADE`)
}

func setup() func() {
	initialContext := context.Background()
	SetupEnv()

	testCfg, _ = config.LoadConfig()

	tracingProvider, err := tracing.NewTracingProvider(testCfg)
	if err != nil {
		log.Default().Println("buildTracing error: " + err.Error())
	}

	dbFactory, err = database.NewPostgresConnectionFactory(initialContext,
		testCfg.Database,
	)
	if err != nil {
		log.Default().Println("banco de dados de teste não iniciado... " + err.Error())
	}

	cacheInstance := cache.NewMemoryCache()
	defer cacheInstance.Close()

	serviceBuilder := app.NewServiceBuilder(dbFactory, tracingProvider, testCfg, cacheInstance)
	serviceFactory, err := serviceBuilder.Build()
	if err != nil {
		log.Default().Println("service factory build error: " + err.Error())
	}

	httpSrv := httpServer.NewServer(initialContext, serviceFactory, testCfg, tracingProvider)
	ts := httptest.NewServer(httpSrv.GetRouter())

	testURL, err = url.Parse(ts.URL)
	if err != nil {
		log.Default().Println("url do servidor de teste inválida... " + err.Error())
	}

	log.Default().Println("server iniciado...", testURL.Port())

	// Atualizar config de teste com URL do servidor
	testCfg.Server.Host = testURL.Hostname()
	testCfg.Server.Port = testURL.Port()

	return func() {
		ts.Close()
		_ = dbFactory.Close()
	}
}

func CleanViper() {
	SetupEnv()
}

func RestClient() *resty.Request {
	token, err := generateToken("11144477735")
	if err != nil {
		log.Fatal(err)
	}

	cli := resty.New()
	cli.SetBaseURL("http://" + net.JoinHostPort(
		testCfg.Server.Host,
		testCfg.Server.Port),
	).SetHeaders(map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": token,
	})
	return cli.R()
}

func SetupEnv() {
	if err := os.Setenv("DB_HOST", "localhost"); err != nil {
		return
	}
	if err := os.Setenv("DB_PORT", "5432"); err != nil {
		return
	}
	if err := os.Setenv("DB_NAME", "test_tech_challenge_db"); err != nil {
		return
	}
	if err := os.Setenv("DB_USER", "techmastertester"); err != nil {
		return
	}
	if err := os.Setenv("DB_PASSWORD", "techmastertester01"); err != nil {
		return
	}
	if err := os.Setenv("DB_SSLMODE", "disable"); err != nil {
		return
	}

	if err := os.Setenv("POOL_MIN", "1"); err != nil {
		return
	}
	if err := os.Setenv("POOL_MAX", "1"); err != nil {
		return
	}
	if err := os.Setenv("REST_SERVER_HOST", ""); err != nil {
		return
	}
	if err := os.Setenv("REST_SERVER_PORT", "8081"); err != nil {
		return
	}
	if err := os.Setenv("JWT_SECRET", "tech-challenge-os-service-secret"); err != nil {
		return
	}
	if err := os.Setenv("KAFKA_PRODUCER_HOSTS", "localhost:9092"); err != nil {
		return
	}
	if err := os.Setenv("KAFKA_CONSUMER_HOSTS", "localhost:9092"); err != nil {
		return
	}
	if err := os.Setenv("TRACING_SERVICE_NAME", "tech-challange-execution-service-test"); err != nil {
		return
	}
	if err := os.Setenv("TRACING_SERVICE_VERSION", "1.0.0-test"); err != nil {
		return
	}
}

func generateToken(cpf string) (string, error) {
	secret := testCfg.JWT.Secret
	expirationTime := time.Now().Add(jwtExpirationTime)

	claims := constant.JWTClaim{
		CPF: cpf,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   cpf,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
