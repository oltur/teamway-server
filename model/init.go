package model

import "fmt"

func init() {
	// init users
	addUserReq := &AddUserRequest{
		UserName: "User1",
		Password: "1",
	}
	user1, err := UserInsert(addUserReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User1 id: %s\n", user1.ID)
	addUserReq = &AddUserRequest{
		UserName: "User2",
		Password: "2",
	}
	user2, err := UserInsert(addUserReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User2 id: %s\n", user2.ID)

	// init tests
	addTestReq := &AddTestRequest{
		Title: "Are you an introvert or an extrovert?",
		Questions: []*Question{
			{
				Title: "You’ve been sitting in the doctor’s waiting room for more than 25 minutes. You:",
				Answers: []*Answer{
					{
						Title: "Look at your watch every two minutes",
						Score: 1,
					},
					{
						Title: "Bubble with inner anger, but keep quiet",
						Score: 2,
					},
					{
						Title: "Explain to other equally impatient people in the room that the doctor is always running late",
						Score: 3,
					},
					{
						Title: "Complain in a loud voice, while tapping your foot impatiently",
						Score: 4,
					},
				},
			},
			{
				Title: "You’re having an animated discussion with a colleague regarding a project that you’re in charge of. You:",
				Answers: []*Answer{
					{
						Title: "Don’t dare contradict them",
						Score: 1,
					},
					{
						Title: "Think that they are obviously right",
						Score: 2,
					},
					{
						Title: "Defend your own point of view, tooth and nail",
						Score: 3,
					},
					{
						Title: "Continuously interrupt your colleague",
						Score: 4,
					},
				},
			},
			{
				Title: "You are taking part in a guided tour of a museum. You:",
				Answers: []*Answer{
					{
						Title: "Are a bit too far towards the back so don’t really hear what the guide is saying",
						Score: 1,
					},
					{
						Title: "Follow the group without question",
						Score: 2,
					},
					{
						Title: "Make sure that everyone is able to hear properly",
						Score: 3,
					},
					{
						Title: "Are right up the front, adding your own comments in a loud voice",
						Score: 4,
					},
				},
			},
			{
				Title: "During dinner parties at your home, you have a hard time with people who:",
				Answers: []*Answer{
					{
						Title: "Ask you to tell a story in front of everyone else",
						Score: 1,
					},
					{
						Title: "Talk privately between themselves",
						Score: 2,
					},
					{
						Title: "Hang around you all evening",
						Score: 3,
					},
					{
						Title: "Always drag the conversation back to themselves",
						Score: 4,
					},
				},
			},
			{
				Title: "You crack a joke at work, but nobody seems to have noticed. You:",
				Answers: []*Answer{
					{
						Title: "Think it’s for the best — it was a lame joke anyway",
						Score: 1,
					},
					{
						Title: "Wait to share it with your friends after work",
						Score: 2,
					},
					{
						Title: "Try again a bit later with one of your colleagues",
						Score: 3,
					},
					{
						Title: "Keep telling it until they pay attention",
						Score: 4,
					},
				},
			},
		},
		ScoreThreshold: 15,
		PositiveResult: "You are extrovert",
		NegativeResult: "You are introvert",
	}
	test1, err := TestInsert(addTestReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Test1 id: %s\n", test1.ID)

	addTestReq = &AddTestRequest{
		Title: "Sample test",
		Questions: []*Question{
			{
				Title: "Question #1",
				Answers: []*Answer{
					{
						Title: "1",
						Score: 1,
					},
					{
						Title: "2",
						Score: 2,
					},
					{
						Title: "3",
						Score: 3,
					},
					{
						Title: "4",
						Score: 4,
					},
				},
			},
			{
				Title: "Question #2",
				Answers: []*Answer{
					{
						Title: "1",
						Score: 1,
					},
					{
						Title: "2",
						Score: 2,
					},
					{
						Title: "3",
						Score: 3,
					},
					{
						Title: "4",
						Score: 4,
					},
				},
			},
			{
				Title: "Question #3",
				Answers: []*Answer{
					{
						Title: "1",
						Score: 1,
					},
					{
						Title: "2",
						Score: 2,
					},
					{
						Title: "3",
						Score: 3,
					},
					{
						Title: "4",
						Score: 4,
					},
				},
			},
			{
				Title: "Question #4",
				Answers: []*Answer{
					{
						Title: "1",
						Score: 1,
					},
					{
						Title: "2",
						Score: 2,
					},
					{
						Title: "3",
						Score: 3,
					},
					{
						Title: "4",
						Score: 4,
					},
				},
			},
			{
				Title: "Question #5",
				Answers: []*Answer{
					{
						Title: "1",
						Score: 1,
					},
					{
						Title: "2",
						Score: 2,
					},
					{
						Title: "3",
						Score: 3,
					},
					{
						Title: "4",
						Score: 4,
					},
				},
			},
		},
		ScoreThreshold: 15,
		PositiveResult: "yes",
		NegativeResult: "no",
	}
	test2, err := TestInsert(addTestReq)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Test2 id: %s\n", test2.ID)

	// init testsTaken
}
