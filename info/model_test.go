package info

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInfo_copyRightYears(t *testing.T) {
	Convey("TestInfo_copyRightYears", t, func() {
		type tmpStruct struct {
			start  uint16
			update uint16
			expect string
		}

		l := []tmpStruct{
			{
				start:  0,
				update: 0,
				expect: "",
			},
			{
				start:  0,
				update: 1,
				expect: "1",
			},
			{
				start:  1,
				update: 0,
				expect: "1",
			},
			{
				start:  1,
				update: 1,
				expect: "1",
			},
			{
				start:  1,
				update: 2,
				expect: "1-2",
			},
			{
				start:  2018,
				update: 2021,
				expect: "2018-2021",
			},
		}

		for k, v := range l {
			Convey(fmt.Sprintf("Index[%d] %d-%d => %s", k, v.start, v.update, v.expect), func() {
				i := Info{}
				i.CopyrightStart = v.start
				i.CopyrightUpdate = v.update

				So(i.copyRightYears(), ShouldEqual, v.expect)
			})
		}

		// Convey("1-2", func() {
		// i := Info{}
		// 	i.CopyrightStart = 1
		// 	i.CopyrightUpdate = 2

		// 	So(i.copyRightYears(), ShouldEqual, "1-2")
		// })
	})
}
