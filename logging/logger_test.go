package logging

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("123 %d", 111)
}

func TestLogger(t *testing.T) {
	logger := New()
	logger.Debugf("%s%d\n", "gda", 123)
}

func BenchmarkBytes(b *testing.B) {
	var buffer bytes.Buffer
	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()

	for i := 0; i < b.N; i++ {
		var buf [64]byte

		buf[0] = byte('0' + (year/1000)%10)
		buf[1] = byte('0' + (year/100)%10)
		buf[2] = byte('0' + (year/10)%10)
		buf[3] = byte('0' + year%10)
		buf[4] = byte('0' + (month/10)%10)
		buf[5] = byte('0' + month%10)
		buf[6] = byte('0' + (day/10)%10)
		buf[7] = byte('0' + day%10)
		buf[8] = ' '
		buf[9] = byte('0' + (hour/10)%10)
		buf[10] = byte('0' + hour%10)
		buf[11] = ':'
		buf[12] = byte('0' + (minute/10)%10)
		buf[13] = byte('0' + minute%10)
		buf[14] = ':'
		buf[15] = byte('0' + (second/10)%10)
		buf[16] = byte('0' + second%10)
		/*
			buf[0] = digits[(year/1000)%10]
			buf[1] = digits[(year/100)%10]
			buf[2] = digits[(year/10)%10]
			buf[3] = digits[year%10]
			buf[4] = digits[(month/10)%10]
			buf[5] = digits[month%10]
			buf[6] = digits[(day/10)%10]
			buf[7] = digits[day%10]
			buf[8] = ' '
			buf[9] = digits[(hour/10)%10]
			buf[10] = digits[hour%10]
			buf[11] = ':'
			buf[12] = digits[(minute/10)%10]
			buf[13] = digits[minute%10]
			buf[14] = ':'
			buf[15] = digits[(second/10)%10]
			buf[16] = digits[second%10]
		*/
		// buf[17] = ' '
		buffer.Reset()
		buffer.Write(buf[:17])
	}
}

func BenchmarkTimeFormat(b *testing.B) {
	var buffer bytes.Buffer
	now := time.Now()
	for i := 0; i < b.N; i++ {
		buf := now.Format("20060102 15:04:05")
		buffer.Reset()
		buffer.WriteString(buf)
	}
}

func BenchmarkBytes2(b *testing.B) {
	var buffer bytes.Buffer
	var buf [64]byte
	for i := 0; i < b.N; i++ {
		line := 9584111
		j := 63
		for j >= 0 {
			buf[j] = digits[line%10]
			j--
			line /= 10
			if line == 0 {
				break
			}
		}
		buffer.Reset()
		buffer.Write(buf[j:])
	}
}

func BenchmarkTimeFormat2(b *testing.B) {
	var buffer bytes.Buffer
	for i := 0; i < b.N; i++ {
		line := 9584111
		buf := strconv.Itoa(line)
		buffer.Reset()
		buffer.WriteString(buf)
	}
}

func BenchmarkFPrint(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		buffer := &bytes.Buffer{}

		for pb.Next() {
			buffer.Reset()
			fmt.Fprintf(buffer, "19181716151ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff413 %s %d %x", "jjjifdao", 957236, 1237)
			buffer.Bytes()
		}
	})
}

func BenchmarkSPrint(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		buffer := &bytes.Buffer{}
		for pb.Next() {
			buffer.Reset()
			buffer.WriteString(fmt.Sprintf("19181716151ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff413 %s %d %x", "jjjifdao", 957236, 1237))
			buffer.Bytes()
		}
	})
}

func www(b *bytes.Buffer, s string) {
	b.WriteString(s)
}
