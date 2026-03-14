package postgres

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/akshay/productiv-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a new GORM database connection and runs auto-migrations.
func NewDB(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	log.Println("database migrations completed")
	return db, nil
}

// runMigrations auto-migrates all domain models.
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Pillar{},
		&domain.FastingSession{},
		&domain.GymSession{},
		&domain.MeditationSession{},
		&domain.RetentionStreak{},
	)
}

// Seed inserts the default user, pillars, and 3 months of dummy data if they don't exist.
func Seed(db *gorm.DB) error {
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	if userCount > 0 {
		return nil
	}

	user := domain.User{
		Name:     "Akshay",
		Timezone: "Asia/Kolkata",
	}
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("seeding user: %w", err)
	}

	pillars := []domain.Pillar{
		{UserID: user.ID, Type: domain.PillarFasting, Name: "Fasting", Icon: "🍽️", Color: "#e94560", IsActive: true, DisplayOrder: 1},
		{UserID: user.ID, Type: domain.PillarGym, Name: "Gym", Icon: "💪", Color: "#4a9eff", IsActive: true, DisplayOrder: 2},
		{UserID: user.ID, Type: domain.PillarMeditation, Name: "Meditation", Icon: "🧘", Color: "#9b59b6", IsActive: true, DisplayOrder: 3},
		{UserID: user.ID, Type: domain.PillarRetention, Name: "Retention", Icon: "🔥", Color: "#2ecc71", IsActive: true, DisplayOrder: 4},
	}
	if err := db.Create(&pillars).Error; err != nil {
		return fmt.Errorf("seeding pillars: %w", err)
	}

	log.Println("database seeded with default user and pillars")

	if err := seedDummyData(db, user.ID); err != nil {
		return fmt.Errorf("seeding dummy data: %w", err)
	}

	log.Println("database seeded with 3 months of dummy data")
	return nil
}

func seedDummyData(db *gorm.DB, userID int64) error {
	rng := rand.New(rand.NewSource(42))
	now := time.Now().UTC()
	startDate := now.AddDate(0, -3, 0) // 3 months ago

	if err := seedFastingSessions(db, userID, rng, startDate, now); err != nil {
		return err
	}
	if err := seedGymSessions(db, userID, rng, startDate, now); err != nil {
		return err
	}
	if err := seedMeditationSessions(db, userID, rng, startDate, now); err != nil {
		return err
	}
	if err := seedRetentionStreaks(db, userID, rng, startDate, now); err != nil {
		return err
	}
	return nil
}

func seedFastingSessions(db *gorm.DB, userID int64, rng *rand.Rand, start, end time.Time) error {
	protocols := []struct {
		name  string
		hours int
	}{
		{"16:8", 16},
		{"18:6", 18},
		{"20:4", 20},
	}

	var sessions []domain.FastingSession
	current := start

	for current.Before(end.Add(-24 * time.Hour)) { // stop before today so there's no active session
		// ~75% chance of fasting on any given day
		if rng.Float64() > 0.75 {
			current = current.Add(24 * time.Hour)
			continue
		}

		proto := protocols[rng.Intn(len(protocols))]
		// Start between 7pm-10pm
		startHour := 19 + rng.Intn(3)
		sessionStart := time.Date(current.Year(), current.Month(), current.Day(), startHour, rng.Intn(60), 0, 0, time.UTC)

		// Actual duration: target +/- 2 hours variation
		actualHours := float64(proto.hours) + (rng.Float64()*4 - 2)
		if actualHours < 4 {
			actualHours = 4
		}
		sessionEnd := sessionStart.Add(time.Duration(actualHours * float64(time.Hour)))
		targetReached := actualHours >= float64(proto.hours)

		sessions = append(sessions, domain.FastingSession{
			UserID:         userID,
			Protocol:       proto.name,
			TargetHours:    proto.hours,
			StartTime:      sessionStart,
			EndTime:        &sessionEnd,
			ActualDuration: &actualHours,
			TargetReached:  targetReached,
			CreatedAt:      sessionStart,
		})

		// Skip 1-2 days after a fast sometimes
		skip := 1
		if rng.Float64() < 0.3 {
			skip = 2
		}
		current = current.Add(time.Duration(skip) * 24 * time.Hour)
	}

	if len(sessions) > 0 {
		return db.CreateInBatches(&sessions, 50).Error
	}
	return nil
}

func seedGymSessions(db *gorm.DB, userID int64, rng *rand.Rand, start, end time.Time) error {
	workoutTypes := []domain.WorkoutType{
		domain.WorkoutStrength, domain.WorkoutCardio, domain.WorkoutHIIT,
		domain.WorkoutYoga, domain.WorkoutSports, domain.WorkoutSwimming,
	}
	durations := []int{30, 45, 60, 75, 90}

	var sessions []domain.GymSession
	current := start
	lastLogged := start.Add(-48 * time.Hour) // ensure first day can log

	for current.Before(end) {
		dayOfWeek := current.Weekday()
		isToday := current.Year() == end.Year() && current.YearDay() == end.YearDay()

		// Skip today so the user can log today themselves
		if isToday {
			break
		}

		// ~70% chance on weekdays, ~40% on weekends
		chance := 0.70
		if dayOfWeek == time.Saturday || dayOfWeek == time.Sunday {
			chance = 0.40
		}

		// Only one session per day, skip if we already logged this day
		if rng.Float64() < chance && current.Sub(lastLogged) >= 20*time.Hour {
			wt := workoutTypes[rng.Intn(len(workoutTypes))]
			// Favor strength/cardio
			if rng.Float64() < 0.5 {
				if rng.Float64() < 0.6 {
					wt = domain.WorkoutStrength
				} else {
					wt = domain.WorkoutCardio
				}
			}

			dur := durations[rng.Intn(len(durations))]
			energy := 2 + rng.Intn(4) // 2-5
			loggedHour := 6 + rng.Intn(12) // 6am-6pm
			loggedAt := time.Date(current.Year(), current.Month(), current.Day(), loggedHour, rng.Intn(60), 0, 0, time.UTC)

			sessions = append(sessions, domain.GymSession{
				UserID:      userID,
				WorkoutType: wt,
				DurationMin: &dur,
				EnergyLevel: &energy,
				LoggedAt:    loggedAt,
				CreatedAt:   loggedAt,
			})
			lastLogged = current
		}

		current = current.Add(24 * time.Hour)
	}

	if len(sessions) > 0 {
		return db.CreateInBatches(&sessions, 50).Error
	}
	return nil
}

func seedMeditationSessions(db *gorm.DB, userID int64, rng *rand.Rand, start, end time.Time) error {
	targetMinutes := []int{5, 10, 15, 20}

	var sessions []domain.MeditationSession
	current := start

	for current.Before(end.Add(-24 * time.Hour)) { // stop before today
		// ~65% chance of meditating on any given day
		if rng.Float64() > 0.65 {
			current = current.Add(24 * time.Hour)
			continue
		}

		target := targetMinutes[rng.Intn(len(targetMinutes))]
		// Morning meditation: 5am-8am
		startHour := 5 + rng.Intn(3)
		sessionStart := time.Date(current.Year(), current.Month(), current.Day(), startHour, rng.Intn(60), 0, 0, time.UTC)

		// Actual: within +/- 3 min of target
		actualMin := float64(target) + (rng.Float64()*6 - 3)
		if actualMin < 2 {
			actualMin = 2
		}
		sessionEnd := sessionStart.Add(time.Duration(actualMin * float64(time.Minute)))

		moodBefore := 1 + rng.Intn(5) // 1-5
		moodAfter := moodBefore + rng.Intn(2) // mood improves or stays
		if moodAfter > 5 {
			moodAfter = 5
		}

		sessions = append(sessions, domain.MeditationSession{
			UserID:         userID,
			TargetMinutes:  target,
			StartTime:      sessionStart,
			EndTime:        &sessionEnd,
			ActualDuration: &actualMin,
			MoodBefore:     &moodBefore,
			MoodAfter:      &moodAfter,
			CreatedAt:      sessionStart,
		})

		current = current.Add(24 * time.Hour)
	}

	if len(sessions) > 0 {
		return db.CreateInBatches(&sessions, 50).Error
	}
	return nil
}

func seedRetentionStreaks(db *gorm.DB, userID int64, rng *rand.Rand, start, end time.Time) error {
	reasons := []string{
		"Lost focus",
		"Stressful week",
		"Relapsed",
	}

	var streaks []domain.RetentionStreak
	current := start

	// Create 3 past streaks of varying lengths, then one active streak
	pastStreakLengths := []int{
		12 + rng.Intn(10), // 12-21 days
		5 + rng.Intn(8),   // 5-12 days
		20 + rng.Intn(15), // 20-34 days
	}

	for _, days := range pastStreakLengths {
		if current.Add(time.Duration(days+3) * 24 * time.Hour).After(end) {
			break
		}

		streakStart := current
		streakEnd := streakStart.Add(time.Duration(days) * 24 * time.Hour)
		reason := reasons[rng.Intn(len(reasons))]

		streaks = append(streaks, domain.RetentionStreak{
			UserID:    userID,
			StartDate: streakStart,
			EndDate:   &streakEnd,
			DaysCount: days,
			Reason:    &reason,
			CreatedAt: streakStart,
		})

		// Gap of 1-3 days between streaks
		gap := 1 + rng.Intn(3)
		current = streakEnd.Add(time.Duration(gap) * 24 * time.Hour)
	}

	// Active streak: started some days ago, still going
	if current.Before(end) {
		activeDays := int(end.Sub(current).Hours() / 24)
		if activeDays > 0 {
			streaks = append(streaks, domain.RetentionStreak{
				UserID:    userID,
				StartDate: current,
				DaysCount: activeDays,
				CreatedAt: current,
			})
		}
	}

	if len(streaks) > 0 {
		return db.CreateInBatches(&streaks, 50).Error
	}
	return nil
}
