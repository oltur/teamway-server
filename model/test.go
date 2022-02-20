package model

import (
	"errors"
	"github.com/oltur/teamway-server/types"
	"github.com/rs/xid"
)

type Test struct {
	ID             types.Id    `json:"id"`
	Title          string      `json:"title"`
	Questions      []*Question `json:"questions"`
	ScoreThreshold int         `json:"scoreThreshold"`
	PositiveResult string      `json:"positiveResult"`
	NegativeResult string      `json:"negativeResult"`
}

type Question struct {
	Title   string    `json:"title"`
	Answers []*Answer `json:"answers"`
}

type Answer struct {
	Title string `json:"title"`
	Score int    `json:"score"`
}

func TestsAll() (res []*Test, err error) {
	allTests := GetMapValuesForTests(testsByIds)
	res = allTests
	return
}

func TestOne(id types.Id) (res *Test, err error) {
	for k := range testsByIds {
		if id == k {
			res = testsByIds[k]
			return
		}
	}
	return nil, ErrNotFound
}

//// TestSave Internal use only
//func TestSave(req *Test) (err error) {
//	testsByIds[req.ID] = req
//	return
//}

func TestDelete(id types.Id) (err error) {
	if _, ok := testsByIds[id]; ok {
		err = ErrNotFound
		return
	}
	delete(testsByIds, id)
	return
}

var testsByIds = make(map[types.Id]*Test)

//func GetMapKeysForTests(m map[types.Id]*Test) (res []types.Id) {
//	res = make([]types.Id, len(m))
//	i := 0
//	for k := range m {
//		res[i] = k
//		i++
//	}
//	return
//}

func GetMapValuesForTests(m map[types.Id]*Test) (res []*Test) {
	res = make([]*Test, len(m))
	i := 0
	for _, v := range m {
		res[i] = v
		i++
	}
	return res
}

func TestInsert(req *AddTestRequest) (res *Test, err error) {
	test := &Test{
		ID:             types.Id(xid.New().String()),
		Title:          req.Title,
		Questions:      req.Questions,
		ScoreThreshold: req.ScoreThreshold,
		PositiveResult: req.PositiveResult,
		NegativeResult: req.NegativeResult,
	}

	if err = req.Validation(); err != nil {
		return
	}

	_, err = TestOne(test.ID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		} else {
			err = nil
		}
	} else {
		err = ErrIdExists
		return
	}

	testsByIds[test.ID] = test
	res = test
	return
}
