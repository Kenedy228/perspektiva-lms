package attempt_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/attempt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should create attempt concrete", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(100).
			withMockQuestions(200).
			build(t, nil)

		//Assert
		assert.NotEqual(t, a.ID(), uuid.Nil)
		assert.NotEqual(t, a.EnrollmentID(), uuid.Nil)
		assert.NotEqual(t, a.QuizID(), uuid.Nil)
		assert.Equal(t, a.CountItems(), 200)
		assert.True(t, a.StartedAt().Before(time.Now()))
		assert.Equal(t, a.DeadlineAt().Sub(a.StartedAt()), time.Duration(time.Second*100))
		assert.True(t, a.FinishedAt().IsZero())
		assert.Equal(t, a.Status(), attempt.StatusInProgress)
	})

	t.Run("should create attempt, limit infinite", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(0).
			withMockQuestions(5).
			build(t, nil)

		//Assert
		assert.True(t, a.DeadlineAt().IsZero())
	})

	t.Run("should return err, if questions nil", func(t *testing.T) {
		//Arrange - Assert
		newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(100).
			build(t, attempt.ErrInvalidAttempt)
	})

	t.Run("should return err, if enrollmentID nil", func(t *testing.T) {
		//Arrange - Assert
		newAttemptBuilder().
			withQuiz(uuid.New()).
			withTimeLimit(100).
			withMockQuestions(1000).
			build(t, attempt.ErrInvalidAttempt)
	})

	t.Run("should return err, if quizID nil", func(t *testing.T) {
		//Arrange - Assert
		newAttemptBuilder().
			withEnrollment(uuid.New()).
			withTimeLimit(100).
			withMockQuestions(1000).
			build(t, attempt.ErrInvalidAttempt)
	})
}

func TestAddAnswer(t *testing.T) {
	t.Run("infinite deadline", func(t *testing.T) {
		t.Run("should add if item exists and has no answer yet", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(0).
				withMockQuestions(1).
				build(t, nil)

			//Act
			items := a.Items()
			err := a.AddAnswer(items[0].ID(), newAnswer())

			//Assert
			assert.NoError(t, err)
		})

		t.Run("should add if item exists and already has answer", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(0).
				withMockQuestions(1).
				build(t, nil)

			//Act
			items := a.Items()
			fErr := a.AddAnswer(items[0].ID(), newAnswer())
			sErr := a.AddAnswer(items[0].ID(), newAnswer())

			//Assert
			assert.NoError(t, fErr)
			assert.NoError(t, sErr)
		})

		t.Run("should return err if id is nil", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(0).
				withMockQuestions(1).
				build(t, nil)

			//Act
			err := a.AddAnswer(uuid.Nil, newAnswer())

			//Assert
			assert.ErrorIs(t, err, attempt.ErrUnexistingItem)
		})

		t.Run("should return err if id is not present", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(0).
				withMockQuestions(1).
				build(t, nil)

			//Act
			err := a.AddAnswer(uuid.New(), newAnswer())

			//Assert
			assert.ErrorIs(t, err, attempt.ErrUnexistingItem)
		})
	})

	t.Run("finite deadline", func(t *testing.T) {
		t.Run("can change before deadline", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(1).
				build(t, nil)

			//Act
			items := a.Items()
			err := a.AddAnswer(items[0].ID(), newAnswer())

			//Assert
			assert.NoError(t, err)
		})

		t.Run("can't change after deadline", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				//Arrange
				a := newAttemptBuilder().
					withEnrollment(uuid.New()).
					withQuiz(uuid.New()).
					withTimeLimit(100000).
					withMockQuestions(1).
					build(t, nil)

				//Act
				time.Sleep(time.Second * 200000)
				a.SetExpired()
				items := a.Items()
				err := a.AddAnswer(items[0].ID(), newAnswer())

				//Assert
				assert.ErrorIs(t, err, attempt.ErrNotModifiable)
				assert.Nil(t, items[0].Answer())
			})
		})
	})
}

func TestFinish(t *testing.T) {
	t.Run("can't finish if expired", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(500).
				withMockQuestions(1).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 1000)
			a.SetExpired()
			err := a.Finish()

			//Assert
			assert.ErrorIs(t, err, attempt.ErrInactive)
			assert.Equal(t, a.Status(), attempt.StatusExpired)
		})
	})

	t.Run("can finish if not finished yet and not expired", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(1000).
			withMockQuestions(1).
			build(t, nil)

		//Act
		err := a.Finish()

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, a.Status(), attempt.StatusFinished)
	})

	t.Run("can finish if not finished yet and has no time limit", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(0).
			withMockQuestions(1).
			build(t, nil)

		//Act
		err := a.Finish()

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, a.Status(), attempt.StatusFinished)
	})
}

func TestSetExpired(t *testing.T) {
	t.Run("can set expired if attempt in progress and has deadline, and time is expired", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(500).
				withMockQuestions(1).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 10000)
			err := a.SetExpired()

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, a.Status(), attempt.StatusExpired)
		})
	})

	t.Run("can't set expired if test has no limit", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(0).
			withMockQuestions(1).
			build(t, nil)

		//Act
		err := a.SetExpired()

		//Assert
		assert.ErrorIs(t, err, attempt.ErrInfiniteDeadline)
		assert.NotEqual(t, a.Status(), attempt.StatusExpired)
	})

	t.Run("can't set expired if attempt is not expired yet", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(1000).
				withMockQuestions(1).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 100)
			err := a.SetExpired()

			//Assert
			assert.ErrorIs(t, err, attempt.ErrNotExpiredYet)
			assert.NotEqual(t, a.Status(), attempt.StatusExpired)
		})
	})
}

func TestCanModify(t *testing.T) {
	t.Run("can't modify if attempt is finished", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(1000).
			withMockQuestions(1).
			build(t, nil)

		//Act
		a.Finish()
		can := a.CanModify()

		//Assert
		assert.False(t, can)
	})

	t.Run("can't modify if attempt is expired", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(1000).
				withTimeLimit(1000).
				withMockQuestions(1).
				build(t, nil)

			//Act
			time.Sleep(time.Second * 1000000)
			a.SetExpired()
			can := a.CanModify()

			//Assert
			assert.False(t, can)
		})
	})

	t.Run("can't modify if attempt is cancelled", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(1000).
			withMockQuestions(1).
			build(t, nil)

		//Act
		a.Cancel()
		can := a.CanModify()

		//Assert
		assert.False(t, can)
	})

	t.Run("can't modify if attempt is interrupted", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(1000).
			withMockQuestions(1).
			build(t, nil)

		//Act
		a.Interrupt()
		can := a.CanModify()

		//Assert
		assert.False(t, can)
	})

	t.Run("can modify if attempt not finished and not expired", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(1).
				build(t, nil)

			//Act
			can := a.CanModify()

			//Assert
			assert.True(t, can)
		})
	})
}

func TestScore(t *testing.T) {
	t.Run("not finished yet", func(t *testing.T) {
		t.Run("in progress", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(1).
				build(t, nil)

			//Act
			score, err := a.Score()

			//Assert
			assert.Equal(t, score, -1)
			assert.ErrorIs(t, err, attempt.ErrNotFinishedYet)
			assert.Equal(t, a.Status(), attempt.StatusInProgress)
		})

		t.Run("expired", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				//Arrange
				a := newAttemptBuilder().
					withEnrollment(uuid.New()).
					withQuiz(uuid.New()).
					withTimeLimit(100000).
					withMockQuestions(1).
					build(t, nil)

				//Act
				time.Sleep(time.Second * 2000000)
				a.SetExpired()
				score, err := a.Score()

				//Assert
				assert.Equal(t, score, -1)
				assert.ErrorIs(t, err, attempt.ErrNotFinishedYet)
				assert.Equal(t, a.Status(), attempt.StatusExpired)
			})
		})
	})

	t.Run("finished", func(t *testing.T) {
		t.Run("all correct answers return len(items) score", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				//Arrange
				questions, answers := makeQuestions(20, 0)

				a := newAttemptBuilder().
					withEnrollment(uuid.New()).
					withQuiz(uuid.New()).
					withTimeLimit(100000).
					withQuestions(questions).
					build(t, nil)

				//Act
				items := a.Items()
				for i := range items {
					a.AddAnswer(items[i].ID(), answers[i])
				}
				a.Finish()
				score, err := a.Score()

				//Assert
				assert.NoError(t, err)
				assert.Equal(t, 20, score)
			})
		})

		t.Run("all incorrect answers return 0 score", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				//Arrange
				questions, answers := makeQuestions(0, 20)

				a := newAttemptBuilder().
					withEnrollment(uuid.New()).
					withQuiz(uuid.New()).
					withTimeLimit(100000).
					withQuestions(questions).
					build(t, nil)

				//Act
				items := a.Items()
				for i := range items {
					a.AddAnswer(items[i].ID(), answers[i])
				}
				a.Finish()
				score, err := a.Score()

				//Assert
				assert.NoError(t, err)
				assert.Equal(t, 0, score)
			})
		})

		t.Run("for correct answer +1, for incorrect +0", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				//Arrange
				questions, answers := makeQuestions(12, 8)

				a := newAttemptBuilder().
					withEnrollment(uuid.New()).
					withQuiz(uuid.New()).
					withTimeLimit(100000).
					withQuestions(questions).
					build(t, nil)

				//Act
				items := a.Items()
				for i := range items {
					a.AddAnswer(items[i].ID(), answers[i])
				}
				a.Finish()
				score, err := a.Score()

				//Assert
				assert.NoError(t, err)
				assert.Equal(t, 12, score)
			})
		})
	})
}

func TestInterrupt(t *testing.T) {
	t.Run("can't interrupt attempt if it is'nt in progress", func(t *testing.T) {
		t.Run("can't interrupt expired attempt", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				//Arrange
				a := newAttemptBuilder().
					withEnrollment(uuid.New()).
					withQuiz(uuid.New()).
					withTimeLimit(100000).
					withMockQuestions(20).
					build(t, nil)

				//Act
				time.Sleep(time.Second * 200000)
				a.SetExpired()
				err := a.Interrupt()

				//Assert
				assert.ErrorIs(t, err, attempt.ErrInactive)
				assert.Equal(t, a.Status(), attempt.StatusExpired)
			})
		})

		t.Run("can't interrupt finished attempt", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(20).
				build(t, nil)

			//Act
			a.Finish()
			err := a.Interrupt()

			//Assert
			assert.ErrorIs(t, err, attempt.ErrInactive)
			assert.Equal(t, a.Status(), attempt.StatusFinished)
		})

		t.Run("can't interrupt cancelled attempt", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(20).
				build(t, nil)

			//Act
			a.Cancel()
			err := a.Interrupt()

			//Assert
			assert.ErrorIs(t, err, attempt.ErrInactive)
			assert.Equal(t, a.Status(), attempt.StatusCancelled)
		})
	})

	t.Run("can interrupt in progress attempt", func(t *testing.T) {
		//Arrange
		a := newAttemptBuilder().
			withEnrollment(uuid.New()).
			withQuiz(uuid.New()).
			withTimeLimit(100000).
			withMockQuestions(20).
			build(t, nil)

		//Act
		err := a.Interrupt()

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, a.Status(), attempt.StatusInterrupted)
	})
}

func TestCancel(t *testing.T) {
	t.Run("can cancel attempt in any status", func(t *testing.T) {
		t.Run("cancel finished attempt", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(10).
				build(t, nil)

			//Act
			a.Finish()
			a.Cancel()

			//Assert
			assert.Equal(t, a.Status(), attempt.StatusCancelled)
		})

		t.Run("can cancel expired attempt", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				//Arrange
				a := newAttemptBuilder().
					withEnrollment(uuid.New()).
					withQuiz(uuid.New()).
					withTimeLimit(100000).
					withMockQuestions(20).
					build(t, nil)

				//Act
				time.Sleep(time.Second * 200000)
				a.Cancel()

				//Assert
				assert.Equal(t, a.Status(), attempt.StatusCancelled)
			})
		})

		t.Run("can cancel interrupted attempt attempt", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(20).
				build(t, nil)

			//Act
			a.Cancel()

			//Assert
			assert.Equal(t, a.Status(), attempt.StatusCancelled)
		})

		t.Run("can cancel cancelled attempt", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(20).
				build(t, nil)

			//Act
			a.Cancel()
			a.Cancel()

			//Assert
			assert.Equal(t, a.Status(), attempt.StatusCancelled)
		})

		t.Run("can cancel in progress attempt", func(t *testing.T) {
			//Arrange
			a := newAttemptBuilder().
				withEnrollment(uuid.New()).
				withQuiz(uuid.New()).
				withTimeLimit(100000).
				withMockQuestions(20).
				build(t, nil)

			//Act
			a.Cancel()

			//Assert
			assert.Equal(t, a.Status(), attempt.StatusCancelled)
		})
	})
}
