package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"

	// init driver.
	_ "github.com/lib/pq"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/api"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(os.Getenv("ENV_FILE"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	time.Local = loc
}

type CalendarSuite struct {
	suite.Suite
	ctx         context.Context
	ctxCancel   context.CancelFunc
	client      api.EventServiceClient
	db          *sqlx.DB
	mu          sync.Mutex
	rabbitCh    *amqp.Channel
	rabbitQueue string
	senderLog   string
}

func (s *CalendarSuite) SetupSuite() {
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())
	s.ctx = metadata.NewOutgoingContext(s.ctx, metadata.Pairs("x-uid", "test-uid"))

	grpcConn, err := grpc.Dial(viper.GetString("APP_HOST"), grpc.WithInsecure())
	if err != nil {
		s.FailNow(err.Error())
	}
	s.client = api.NewEventServiceClient(grpcConn)

	db, err := sqlx.Connect("postgres", viper.GetString("DB_DSN"))
	if err != nil {
		s.FailNow(err.Error())
	}
	s.db = db

	amqpConn, err := amqp.Dial(viper.GetString("RABBIT_DSN"))
	if err != nil {
		s.FailNow(err.Error())
	}

	ch, err := amqpConn.Channel()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.rabbitCh = ch
	s.rabbitQueue = viper.GetString("RABBIT_QUEUE")
	s.senderLog = viper.GetString("SENDER_LOG_FILE")
}

func (s *CalendarSuite) TearDownSuite() {
	defer s.db.Close()
	defer s.ctxCancel()
}

func (s *CalendarSuite) SetupTest() {
	s.mu.Lock()
	s.cleanupDb()
	s.cleanupRabbit()
}

func (s *CalendarSuite) TearDownTest() {
	s.mu.Unlock()
}

func (s *CalendarSuite) TestAddEvent() {
	now := time.Now()
	beginDate := s.toProtoTimestamp(now.Add(time.Hour))
	endDate := s.toProtoTimestamp(now.Add(2 * time.Hour))

	event, err := s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "title-1",
		Description: "description-1",
		BeginDate:   beginDate,
		EndDate:     endDate,
	})
	if err != nil {
		s.FailNow(err.Error())
		return
	}

	s.Equal("title-1", event.Title)
	s.Equal("description-1", event.Description)
	s.Equal(beginDate.AsTime(), event.BeginDate.AsTime())
	s.Equal(endDate.AsTime(), event.EndDate.AsTime())
}

func (s *CalendarSuite) TestDayList() {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	e, err := s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "t1",
		Description: "d1",
		BeginDate:   s.toProtoTimestamp(today.Add(23 * time.Hour)),
		EndDate:     s.toProtoTimestamp(today.Add(24 * time.Hour)),
	})
	if err != nil {
		s.FailNow(err.Error())
	}
	fmt.Println(e)

	e, err = s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "t2",
		Description: "d2",
		BeginDate:   s.toProtoTimestamp(today.Add(25*time.Hour + time.Second)),
		EndDate:     s.toProtoTimestamp(today.Add(26 * time.Hour)),
	})
	if err != nil {
		s.FailNow(err.Error())
	}
	fmt.Println(e)

	events, err := s.client.DayList(s.ctx, &api.DayListRequest{
		Day: s.toProtoTimestamp(now),
	})
	s.Require().NoError(err)
	s.Require().Len(events.Items, 1)
	s.Require().Equal("t1", events.Items[0].Title)
}

func (s *CalendarSuite) TestWeekList() {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	e, err := s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "t1",
		Description: "d1",
		BeginDate:   s.toProtoTimestamp(today.AddDate(0, 0, 6).Add(23 * time.Hour)),
		EndDate:     s.toProtoTimestamp(today.AddDate(0, 0, 7)),
	})
	if err != nil {
		s.FailNow(err.Error())
	}
	fmt.Println(e)

	e, err = s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "t2",
		Description: "d2",
		BeginDate:   s.toProtoTimestamp(today.AddDate(0, 0, 7)),
		EndDate:     s.toProtoTimestamp(today.AddDate(0, 0, 7).Add(time.Hour)),
	})
	if err != nil {
		s.FailNow(err.Error())
	}
	fmt.Println(e)

	events, err := s.client.WeekList(s.ctx, &api.WeekListRequest{
		BeginDate: s.toProtoTimestamp(now),
	})
	s.Require().NoError(err)
	s.Require().Len(events.Items, 1)
	s.Require().Equal("t1", events.Items[0].Title)
}

func (s *CalendarSuite) TestMonthList() {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	e, err := s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "t1",
		Description: "d1",
		BeginDate:   s.toProtoTimestamp(today.AddDate(0, 1, 0).Add(-time.Hour)),
		EndDate:     s.toProtoTimestamp(today.AddDate(0, 1, 0)),
	})
	if err != nil {
		s.FailNow(err.Error())
	}
	fmt.Println(e)

	e, err = s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "t2",
		Description: "d2",
		BeginDate:   s.toProtoTimestamp(today.AddDate(0, 1, 0)),
		EndDate:     s.toProtoTimestamp(today.AddDate(0, 1, 0).Add(time.Hour)),
	})
	if err != nil {
		s.FailNow(err.Error())
	}
	fmt.Println(e)

	events, err := s.client.MonthList(s.ctx, &api.MonthListRequest{
		BeginDate: s.toProtoTimestamp(now),
	})
	s.Require().NoError(err)
	s.Require().Len(events.Items, 1)
	s.Require().Equal("t1", events.Items[0].Title)
}

func (s *CalendarSuite) TestSender() {
	now := time.Now()
	_, err := s.client.Create(s.ctx, &api.CreateEventRequest{
		Title:       "t1",
		Description: "d1",
		BeginDate:   s.toProtoTimestamp(now.Add(5 * time.Minute)),
		EndDate:     s.toProtoTimestamp(now.Add(time.Hour)),
	})
	if err != nil {
		s.FailNow(err.Error())
	}

	time.Sleep(7 * time.Second)

	b, err := ioutil.ReadFile(s.senderLog)
	if err != nil {
		s.FailNow(err.Error())
	}
	str := string(b)

	s.Require().True(strings.Contains(str, `\"Title\":\"t1\"`))
}

func (s *CalendarSuite) toProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	protoTs, err := ptypes.TimestampProto(t)
	if err != nil {
		s.FailNow(err.Error())
	}
	return protoTs
}

func (s *CalendarSuite) cleanupDb() {
	_, err := s.db.ExecContext(s.ctx, "DROP SCHEMA public CASCADE; "+
		"CREATE SCHEMA public; "+
		"GRANT ALL ON SCHEMA public TO postgres; "+
		"GRANT ALL ON SCHEMA public TO public;",
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	err = goose.Up(s.db.DB, viper.GetString("MIGRATIONS_DIR"))
	if err != nil {
		s.FailNow(err.Error())
	}
}

func (s *CalendarSuite) cleanupRabbit() {
	_, err := s.rabbitCh.QueuePurge(s.rabbitQueue, false)
	if err != nil {
		s.FailNow(err.Error())
	}

	//file, err := os.OpenFile(s.senderLog, os.O_RDWR|os.O_CREATE, 0777)
	file, err := os.OpenFile(s.senderLog, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		s.FailNow(err.Error())
		return
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		s.FailNow(err.Error())
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		s.FailNow(err.Error())
	}
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
