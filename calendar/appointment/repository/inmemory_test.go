package repository

import (
    "testing"
    "time"
    "github.com/stretchr/testify/assert"

    "github.com/eantyshev/otus_go/calendar/models"
    "github.com/eantyshev/otus_go/calendar/appointment"
)

var ap1 = models.Appointment{
    ID:          1,
    Summary:     "go seminar",
    Description: "awesome Go seminar",
    StartsAt:   time.Date(2019, 10, 29, 17, 0, 0, 0, time.UTC),
    IsRegular:   true,
    DaysOfWeek:        []models.WeekDay{models.Tuesday, models.Thursday},
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

func TestNotFound(t *testing.T) {
    var r appointment.Repository = NewMapRepo()
    if _, err := r.GetById(1); err != models.ErrIdNotFound {
        t.Fatalf("Unexpected error: %s", err)
    }
}
