package repository

import (
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	uuid1 = uuid.MustParse("11112222333344445555666677778888")
	uuid2 = uuid.MustParse("11112222333344445555666677779999")
)

var ap1 = appointment.Appointment{
	Uuid:        uuid1,
	Summary:     "go seminar",
	Description: "awesome Go seminar",
	StartsAt:    time.Date(2019, 10, 29, 17, 0, 0, 0, time.UTC),
	Duration:    100 * time.Minute,
	Owner:       "user1",
}

var ap2 = appointment.Appointment{
	Uuid:        uuid2,
	Summary:     "sleep time",
	Description: "wonderful dreams about programming Go",
	StartsAt:    time.Date(2019, 10, 29, 20, 0, 0, 0, time.UTC),
	Duration:    8 * time.Hour,
	Owner:       "user1",
}

func TestStore(t *testing.T) {
	var r Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	if apStored, err := r.GetById(uuid1); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, *apStored, ap1)
	}
}

func TestStoreConflict(t *testing.T) {
	var r Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	var apCopy = ap1
	if err := r.Store(&apCopy); err != appointment.ErrConflictId {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	var r Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	var apCopy = ap1
	apCopy.Summary = "modified summary"
	if err := r.Update(&apCopy); err != nil {
		t.Fatal(err)
	}
	if apStored, err := r.GetById(uuid1); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, *apStored, apCopy)
	}
}

func TestNotFound(t *testing.T) {
	var r Repository = NewMapRepo()
	if _, err := r.GetById(uuid1); err != appointment.ErrIdNotFound {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestFetch(t *testing.T) {
	var r Repository = NewMapRepo()
	r.Store(&ap1)
	r.Store(&ap2)
	aps, timeEnd, err := r.Fetch(time.Date(2019, 10, 29, 0, 0, 0, 0, time.UTC), 2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, aps, &ap1)
	assert.Contains(t, aps, &ap2)
	assert.Equal(t, len(aps), 2)
	assert.Equal(t, timeEnd, time.Date(2019, 10, 30, 4, 0, 0, 0, time.UTC))
}

func TestFetchRepeat(t *testing.T) {
	var r Repository = NewMapRepo()
	r.Store(&ap1)
	r.Store(&ap2)
	aps, timeEnd, err := r.Fetch(time.Date(2019, 10, 29, 0, 0, 0, 0, time.UTC), 1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(aps), 1)
	apsOther, timeEnd2, err2 := r.Fetch(timeEnd, 1)
	if err2 != nil {
		t.Fatal(err2)
	}
	aps = append(aps, apsOther...)
	assert.Contains(t, aps, &ap1)
	assert.Contains(t, aps, &ap2)
	assert.Equal(t, len(aps), 2)
	assert.Equal(t, timeEnd2, time.Date(2019, 10, 30, 4, 0, 0, 0, time.UTC))
}

func TestDelete(t *testing.T) {
	var r Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(uuid1); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(uuid1); err != appointment.ErrIdNotFound {
		t.Fatalf("Unexpected err: %s", err)
	}
}
