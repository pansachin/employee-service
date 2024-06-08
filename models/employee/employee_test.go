package employee_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"

	"github.com/pansachin/employee-service/models/employee"
	"github.com/pansachin/employee-service/pkg/database"
	"github.com/pansachin/employee-service/pkg/database/dbtest"
	"github.com/pansachin/employee-service/pkg/validate"
)

type TestSuite struct {
	db       *sqlx.DB
	log      *slog.Logger
	teardown func()

	tableName string
	ctx       context.Context
	rwmux     *sync.RWMutex
}

var ts TestSuite
var TableName = "employee"

func TestMain(m *testing.M) {
	success := m.Run()
	ts.teardown()
	os.Exit(success)
}

func setupTestDB() error {
	// Make sure we have sufficient permission for the db user
	q := fmt.Sprintf("create table if not exists test_db.%s like %s.%s", ts.tableName, dbtest.UnitDbConfig.Name, ts.tableName)
	if _, err := ts.db.ExecContext(ts.ctx, q); err != nil {
		return fmt.Errorf("creating test_db.%s test table: %v", ts.tableName, err)
	}

	q = fmt.Sprintf("truncate table test_db.%s", ts.tableName)
	if _, err := ts.db.ExecContext(ts.ctx, q); err != nil {
		return fmt.Errorf("truncating test_db.%s test table: %v", ts.tableName, err)
	}

	q = fmt.Sprintf(`
		insert into test_db.%s
		select * from %s.%s
		order by created_on desc limit 100`,
		ts.tableName,
		dbtest.UnitDbConfig.Name,
		ts.tableName,
	)
	if _, err := ts.db.ExecContext(ts.ctx, q); err != nil {
		return fmt.Errorf("error copying data to test_db.%s test table: %v", ts.tableName, err)
	}

	return nil
}

func registerTestSuite(t *testing.T) {
	if ts.db == nil {
		log, db, teardown := dbtest.NewUnit(t)
		ctx := context.Background()
		rwmux := &sync.RWMutex{}
		ts = TestSuite{
			db:        db,
			log:       log,
			teardown:  teardown,
			ctx:       ctx,
			tableName: TableName,
			rwmux:     rwmux,
		}

		t.Logf("Create test database test_db.%s", ts.tableName)
		if err := setupTestDB(); err != nil {
			t.Fatalf("Failed to create test_db.%s: %v", ts.tableName, err)
		}
	}
}

func Test_EmployeeCRUD(t *testing.T) {
	registerTestSuite(t)

	// Use throughout the test
	rsc := employee.NewCore(ts.log, ts.db, ts.rwmux)

	t.Log("Given the need to work with Employees")
	{
		testID := 1
		t.Logf("\tTest %d:\tWhen handling a single Employee.", testID)
		{
			// hard coded for easy testing
			now := time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC)

			var rt employee.NewEmployee
			data := rt.GenerateFakeData(1)

			// CREATE
			newRecord, err := rsc.Create(ts.ctx, data[0], now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to create Employee : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to create Employee", dbtest.Success, testID)
			testID++

			// QUERY BY ID
			fetchedRecord, err := rsc.QueryByID(ts.ctx, newRecord.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve Employee by ID: %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve Employee by ID.", dbtest.Success, testID)
			testID++

			// DIFF New and Fetch
			if diff := cmp.Diff(newRecord, fetchedRecord); diff != "" {
				t.Fatalf("\t%s\tTest %d:\tShould get back the same Employee. Diff:\n%s", dbtest.Failed, testID, diff)
			}
			t.Logf("\t%s\tTest %d:\tShould get back the same Employee", dbtest.Success, testID)
			testID++

			// UPDATE - NON-EXISTING RECORD
			us := employee.UpdateEmployee{
				Position: dbtest.StringPointer("Seniro Software Engineer"),
			}
			err = rsc.Update(ts.ctx, "923498273", us, now)
			if !errors.Is(err, employee.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to update non-existing source : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:Should NOT be able to update non-existing source", dbtest.Success, testID)
			testID++

			// UPDATE
			if err := rsc.Update(ts.ctx, newRecord.ID, us, now); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to update the Employee: %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to update the Employee", dbtest.Success, testID)
			testID++

			// SOFT DELETE
			if err := rsc.Delete(ts.ctx, newRecord.ID, time.Now()); err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to soft delete Employee : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to soft delete Employee", dbtest.Success, testID)
			testID++

			// QUERY BY ID FOR DELETED
			_, err = rsc.QueryByID(ts.ctx, newRecord.ID)
			if !errors.Is(err, employee.ErrNotFound) {
				t.Fatalf("\t%s\tTest %d:\tShould NOT be able to retrieve Employee by id : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:Should NOT be able to retrieve Employee by id", dbtest.Success, testID)
			testID++

		}
	}
}

func Test_EmployeePagination(t *testing.T) {
	registerTestSuite(t)

	// Use throughout the test
	rsc := employee.NewCore(ts.log, ts.db, ts.rwmux)

	// Generate some data
	var rt employee.NewEmployee
	data := rt.GenerateFakeData(3)
	if err := rsc.Seed(ts.ctx, data); err != nil {
		t.Fatalf("Failed seeding for pagination: %v", err)
	}

	t.Log("Given the need to page through Employee records.")
	{
		testID := 1
		t.Logf("\tTest %d:\tWhen paging through 2 statues.", testID)
		{
			pagi := database.NewPagination()
			pagi.PerPage = 1

			// GET FIRST RECORD
			s1, err := rsc.Query(ts.ctx, pagi)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve Employee for page 1 : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve Employee for page 1", dbtest.Success, testID)
			testID++

			if len(s1) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single Employee", dbtest.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single Employee", dbtest.Success, testID)
			testID++

			// GET SECOND RECORD
			pagi.Page = 1
			s2, err := rsc.Query(ts.ctx, pagi)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve Employee for page 2 : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve Employee for page 2", dbtest.Success, testID)
			testID++

			if len(s2) != 1 {
				t.Fatalf("\t%s\tTest %d:\tShould have a single Employee : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have a single Employee", dbtest.Success, testID)
			testID++

			// COMPARE
			if s1[0].ID == s2[0].ID {
				t.Logf("\t\tTest %d:\tEmployee 1: %v", testID, s1[0].ID)
				t.Logf("\t\tTest %d:\tEmployee 2: %v", testID, s2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different Employees : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould have different Employees", dbtest.Success, testID)
			testID++

			// GET 3 RECORDS AND MAKE SURE THE ABOVE 2 MATCH
			pagi.Page = 0
			pagi.PerPage = 3
			three, err := rsc.Query(ts.ctx, pagi)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tShould be able to retrieve Employee for 3 records : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to retrieve Employee for 3 records", dbtest.Success, testID)
			testID++

			if len(three) != 3 {
				t.Fatalf("\t%s\tTest %d:\tShould 3 Employees : %s", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould 3 Employees", dbtest.Success, testID)
			testID++

			// COMPARE
			if s1[0].ID != three[0].ID || s2[0].ID != three[1].ID {
				t.Logf("\t\tTest %d:\tEmployee 1 - Expected: %v, Got: %v", testID, three[0].ID, s1[0].ID)
				t.Logf("\t\tTest %d:\tEmployee 2 - Expected: %v, Got: %v", testID, three[1].ID, s2[0].ID)
				t.Fatalf("\t%s\tTest %d:\tShould have different Employees", dbtest.Failed, testID)
			}
			t.Logf("\t%s\tTest %d:\tShould have different Employees", dbtest.Success, testID)
		}
	}
}
func Test_EmployeeCRUDValidation(t *testing.T) {
	registerTestSuite(t)

	// Use throughout the test
	rtc := employee.NewCore(ts.log, ts.db, ts.rwmux)

	t.Log("Given the need to get specific CRUD model validation error messages")
	{
		testID := 1
		t.Logf("\tTest %d:\tWhen dealing with status CRUD validation errors", testID)
		{
			// hard coded for easy testing
			now := time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC)

			// CREATE Validation Check
			_, err := rtc.Create(ts.ctx, employee.NewEmployee{}, now)
			if !validate.IsFieldErrors(err) {
				t.Fatalf("\t%s\tTest %d [Create]:\tExpecting validation field errors : %v", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d [Create]:\tExpecting validation field errors", dbtest.Success, testID)
			testID++

			errs := validate.GetFieldErrors(err).Fields()
			expecting := 1
			if len(errs) != expecting {
				t.Errorf("\t%s\tTest %d [Create]:\tValidation error count", dbtest.Failed, testID)
				t.Logf("\t\t\tTest %d:\tGot: %v", testID, len(errs))
				t.Logf("\t\t\tTest %d:\tExp: %v", testID, expecting)
				t.Logf("\t\t\tTest %d:\tError Message: %v", testID, err)
			} else {
				t.Logf("\t%s\tTest %d [Create]:\tValidation error count.", dbtest.Success, testID)
			}

		}
	}
}
