package phpgo

import (
	"fmt"
	"testing"
	"time"
)

func TestDateStrToTime(t *testing.T) {
	date := "2017-07-04 12:23:12"
	t1 := DateStrToTime(date)
	t.Log(t1.Unix())
}

func TestDateFormat(t *testing.T) {
	format := "Y y m n M F d j D l g G h H a A i s"
	s := DateFormat(format)
	t.Log("DateFormat", s)
}

func TestDateDate4(t *testing.T) {
	s := DateDate("U", 1499075877)
	t.Log("DateDate", s)
}

func TestDateDate3(t *testing.T) {
	t.Log(fmt.Sprintf("%d", 1499075877))
	t.Logf("%d", 1499075877)
	s := DateDate("g G u e O P c r U", 1499075877)
	t.Log(time.Now().Month().String())
	t.Log(time.Now().Unix())
	t.Log("DateDate", s)
}

func TestDateDate2(t *testing.T) {
	s := DateDate("M n y a A g", 1482278400)
	t.Log(time.Now().Month().String())
	t.Log("DateDate", s)
}

func TestDateDate1(t *testing.T) {
	t.Log(time.Now().Weekday().String())
	t.Log(time.Now().YearDay())
	s := DateDate("Y-m-d H:i:s D j l", 0)
	t.Log("DateDate", s)
}

func TestDateDate(t *testing.T) {
	s := DateDate("Y-m-d H:i:s", 0)
	t.Log("DateDate", s)
}

func TestDatemdHi(t *testing.T) {
	s := DatemdHi()
	t.Log("DatemdHi", s)
}

func TestDateYmdHi(t *testing.T) {
	s := DateYmdHi()
	t.Log("DateYmdHi", s)
}

func TestDateYmd(t *testing.T) {
	s := DateYmd()
	t.Log("DateYmd", s)
}

func TestDateHis(t *testing.T) {
	s := DateHis()
	t.Log("DateHis", s)
}

func TestDateYmdHis(t *testing.T) {
	s := DateYmdHis()
	t.Log("DateYmdHis", s)
}
