package test

import (
	"context"
	"errors"
	"fmt"
	"iam-permission/internal/app/api/component"
	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"

	"iam-permission/internal/app/api/route"
	"iam-permission/internal/app/migration"
	"iam-permission/internal/pkg/config"
	"iam-permission/internal/pkg/db"
	"iam-permission/internal/pkg/redis"
)

var (
	defaultMigrationDir = "file://../../../../migration/postgres"

	testPublicRouter *gin.Engine
	testAdminRouter  *gin.Engine
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupSuite() {
	fmt.Println("============ Start running tests for package `api`... ============")

	// Load config
	err := config.Load("../../../pkg/config/file/test.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Init components
	err = component.InitComponents()
	if err != nil {
		log.Fatal(err)
	}

	// Migrate DB
	err = migrateDB()
	if err != nil {
		log.Fatal(err)
	}

	// HTTP mocks
	//httpmock.ActivateNonDefault(...)

	// setup router
	if testPublicRouter == nil {
		gin.SetMode(config.Instance().Http.Mode)
		testPublicRouter = gin.New()
		route.Register(testPublicRouter, false)
	}

	if testAdminRouter == nil {
		gin.SetMode(config.Instance().Http.Mode)
		testAdminRouter = gin.New()
		route.Register(testAdminRouter, true)
	}
}

func migrateDB() error {
	m, err := migration.NewMigration(db.Instance(), defaultMigrationDir, "")
	if err != nil {
		return err
	}

	err = m.MigrateUp(0)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (suite *TestSuite) SetupTest() {
	// http
	httpmock.Reset()

	// postgres
	if err := db.Instance().Exec("TRUNCATE TABLE services CASCADE").Error; err != nil {
		panic(err)
	}

	// redis
	redis.ClientInstance().FlushAll(context.Background())
}

func (suite *TestSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
