package repository

import (
	"github.com/eantyshev/otus_go/calendar/pkg/appointment"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/eantyshev/otus_go/calendar/pkg/models"
)

var ap1 = models.Appointment{
	ID:              1,
	Summary:         "go seminar",
	Description:     "awesome Go seminar",
	StartsAt:        time.Date(2019, 10, 29, 17, 0, 0, 0, time.UTC),
	DurationMinutes: 100,
	IsRegular:       true,
	DaysOfWeek:      []time.Weekday{time.Tuesday, time.Thursday},
}

var ap2 = models.Appointment{
	ID:              2,
	Summary:         "sleep time",
	Description:     "wonderful dreams about programming Go",
	StartsAt:        time.Date(2019, 10, 29, 20, 0, 0, 0, time.UTC),
	DurationMinutes: 8 * 60,
	IsRegular:       true,
}

func TestStore(t *testing.T) {
	var r appointment.Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	if apStored, err := r.GetById(1); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, *apStored, ap1)
	}
}

func TestStoreConflict(t *testing.T) {
	var r appointment.Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	var apCopy = ap1
	if err := r.Store(&apCopy); err != models.ErrConflictId {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	var r appointment.Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	var apCopy = ap1
	apCopy.Summary = "modified summary"
	if err := r.Update(&apCopy); err != nil {
		t.Fatal(err)
	}
	if apStored, err := r.GetById(1); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, *apStored, apCopy)
	}
}

func TestNotFound(t *testing.T) {
	var r appointment.Repository = NewMapRepo()
	if _, err := r.GetById(1); err != models.ErrIdNotFound {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestFetch(t *testing.T) {
	var r appointment.Repository = NewMapRepo()
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
	var r appointment.Repository = NewMapRepo()
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
	var r appointment.Repository = NewMapRepo()
	if err := r.Store(&ap1); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(1); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(1); err != models.ErrIdNotFound {
		t.Fatalf("Unexpected err: %s", err)
	}
}
