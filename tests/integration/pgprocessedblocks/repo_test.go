package pgprocessedblocks

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/allora-network/allora-producer/app/domain"
	"github.com/allora-network/allora-producer/infra"
)

type ProcessedBlockRepositoryTestSuite struct {
	suite.Suite
	ctx context.Context

	db                 *pgxpool.Pool
	pgContainer        *postgres.PostgresContainer
	pgConnectionString string
}

var fixedTime = time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

func (s *ProcessedBlockRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	pgContainer, err := postgres.Run(s.ctx, "postgres:16",
		postgres.WithDatabase("producerdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	s.Require().NoError(err)

	connStr, err := pgContainer.ConnectionString(s.ctx, "sslmode=disable")
	s.Require().NoError(err)

	db, err := pgxpool.Connect(s.ctx, connStr)
	s.Require().NoError(err)

	s.pgContainer = pgContainer
	s.db = db
	s.pgConnectionString = connStr

	err = s.db.Ping(s.ctx)
	s.NoError(err)
}

func (s *ProcessedBlockRepositoryTestSuite) TearDownSuite() {
	s.db.Close()
	err := s.pgContainer.Terminate(s.ctx)
	s.NoError(err)
}

func (s *ProcessedBlockRepositoryTestSuite) SetupTest() {
	err := infra.CreateTables(s.ctx, s.db)
	s.Require().NoError(err)
}

func (s *ProcessedBlockRepositoryTestSuite) TearDownTest() {
	err := infra.DropTables(s.ctx, s.db)
	s.Require().NoError(err)
}

func (s *ProcessedBlockRepositoryTestSuite) TestInitialDatabaseState() {
	var count int
	err := s.db.QueryRow(s.ctx, "SELECT COUNT(*) FROM processed_blocks").Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(0, count)

	err = s.db.QueryRow(s.ctx, "SELECT COUNT(*) FROM processed_block_events").Scan(&count)
	s.Require().NoError(err)
	s.Require().Equal(0, count)
}

func (s *ProcessedBlockRepositoryTestSuite) TestSaveProcessedBlock() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	block := domain.ProcessedBlock{
		Height:      1,
		ProcessedAt: fixedTime,
		Status:      domain.StatusCompleted,
	}
	err = repo.SaveProcessedBlock(s.ctx, block)
	s.Require().NoError(err)

	lastBlock, err := repo.GetLastProcessedBlock(s.ctx)
	s.Require().NoError(err)
	// ID is autogenerated by the database
	s.Require().Equal(block.ProcessedAt.Unix(), lastBlock.ProcessedAt.Unix())
	s.Require().Equal(block.Status, lastBlock.Status)
	s.Require().Equal(block.Height, lastBlock.Height)
}

func (s *ProcessedBlockRepositoryTestSuite) TestGetLastProcessedBlock() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	block := domain.ProcessedBlock{
		Height:      1,
		ProcessedAt: fixedTime,
		Status:      domain.StatusCompleted,
	}
	err = repo.SaveProcessedBlock(s.ctx, block)
	s.Require().NoError(err)

	block2 := domain.ProcessedBlock{
		Height:      2,
		ProcessedAt: fixedTime,
		Status:      domain.StatusCompleted,
	}
	err = repo.SaveProcessedBlock(s.ctx, block2)
	s.Require().NoError(err)

	block3 := domain.ProcessedBlock{
		Height:      3,
		ProcessedAt: fixedTime,
		Status:      domain.StatusFailed,
	}
	err = repo.SaveProcessedBlock(s.ctx, block3)
	s.Require().NoError(err)

	lastBlock, err := repo.GetLastProcessedBlock(s.ctx)
	s.Require().NoError(err)
	// ID is autogenerated by the database
	s.Require().Equal(block2.ProcessedAt.Unix(), lastBlock.ProcessedAt.Unix())
	s.Require().Equal(block2.Status, lastBlock.Status)
	s.Require().Equal(block2.Height, lastBlock.Height)
}

func (s *ProcessedBlockRepositoryTestSuite) TestGetLastProcessedBlock_NoBlocks() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	lastBlock, err := repo.GetLastProcessedBlock(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(domain.ProcessedBlock{}, lastBlock)
}

func (s *ProcessedBlockRepositoryTestSuite) TestGetLastProcessedBlock_NoCompletedBlocks() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	block := domain.ProcessedBlock{
		Height:      1,
		ProcessedAt: fixedTime,
		Status:      domain.StatusFailed,
	}
	err = repo.SaveProcessedBlock(s.ctx, block)
	s.Require().NoError(err)

	lastBlock, err := repo.GetLastProcessedBlock(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(domain.ProcessedBlock{}, lastBlock)
}

func (s *ProcessedBlockRepositoryTestSuite) TestSaveProcessedBlockEvent() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	event := domain.ProcessedBlockEvent{
		Height:      1,
		ProcessedAt: fixedTime,
		Status:      domain.StatusCompleted,
	}
	err = repo.SaveProcessedBlockEvent(s.ctx, event)
	s.Require().NoError(err)

	lastEvent, err := repo.GetLastProcessedBlockEvent(s.ctx)
	s.Require().NoError(err)
	// ID is autogenerated by the database
	s.Require().Equal(event.ProcessedAt.Unix(), lastEvent.ProcessedAt.Unix())
	s.Require().Equal(event.Status, lastEvent.Status)
	s.Require().Equal(event.Height, lastEvent.Height)
}

func (s *ProcessedBlockRepositoryTestSuite) TestGetLastProcessedBlockEvent() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	event := domain.ProcessedBlockEvent{
		Height:      1,
		ProcessedAt: fixedTime,
		Status:      domain.StatusCompleted,
	}
	err = repo.SaveProcessedBlockEvent(s.ctx, event)
	s.Require().NoError(err)

	event2 := domain.ProcessedBlockEvent{
		Height:      2,
		ProcessedAt: fixedTime,
		Status:      domain.StatusCompleted,
	}
	err = repo.SaveProcessedBlockEvent(s.ctx, event2)
	s.Require().NoError(err)

	event3 := domain.ProcessedBlockEvent{
		Height:      3,
		ProcessedAt: fixedTime,
		Status:      domain.StatusFailed,
	}
	err = repo.SaveProcessedBlockEvent(s.ctx, event3)
	s.Require().NoError(err)

	lastEvent, err := repo.GetLastProcessedBlockEvent(s.ctx)
	s.Require().NoError(err)
	// ID is autogenerated by the database
	s.Require().Equal(event2.ProcessedAt.Unix(), lastEvent.ProcessedAt.Unix())
	s.Require().Equal(event2.Status, lastEvent.Status)
	s.Require().Equal(event2.Height, lastEvent.Height)
}

func (s *ProcessedBlockRepositoryTestSuite) TestGetLastProcessedBlockEvent_NoEvents() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	lastEvent, err := repo.GetLastProcessedBlockEvent(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(domain.ProcessedBlockEvent{}, lastEvent)
}

func (s *ProcessedBlockRepositoryTestSuite) TestGetLastProcessedBlockEvent_NoCompletedEvents() {
	repo, err := infra.NewPgProcessedBlock(s.db)
	s.Require().NoError(err)

	event := domain.ProcessedBlockEvent{
		Height:      1,
		ProcessedAt: fixedTime,
		Status:      domain.StatusFailed,
	}
	err = repo.SaveProcessedBlockEvent(s.ctx, event)
	s.Require().NoError(err)

	lastEvent, err := repo.GetLastProcessedBlockEvent(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(domain.ProcessedBlockEvent{}, lastEvent)
}

func TestProcessedBlockRepository(t *testing.T) {
	suite.Run(t, new(ProcessedBlockRepositoryTestSuite))
}
