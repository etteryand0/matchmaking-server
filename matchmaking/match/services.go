package match

import (
	"encoding/json"
	"errors"
	"etteryand0/matchmaking/server/common"
	"etteryand0/matchmaking/server/models"
	"fmt"
	"math"
)

type TeamData struct {
	MMRMedian float64
	Users     map[string]UserData
}

type UserData struct {
	ID   string
	Role string
	MMR  int
}

func (m *Match) CalculateScore() (float64, error) {
	if len(m.Teams) != 2 {
		return 0, errors.New("invalid teams count")
	}
	if len(m.Teams[0].Users) != 5 || len(m.Teams[1].Users) != 5 {
		return 0, errors.New("invalid team members count")
	}

	var teams [2]TeamData
	rolesMMRDifferences := 0
	score := 0.0
	for teamIndex, team := range m.Teams {
		teamData := TeamData{
			MMRMedian: 0,
			Users:     make(map[string]UserData),
		}

		var hasTop, hasMid, hasBot, hasSup, hasJungle bool

		var userIds []string
		for _, user := range team.Users {
			userIds = append(userIds, user.ID)
			teamData.Users[user.ID] = UserData{
				ID:   user.ID,
				Role: user.Role,
			}
			switch user.Role {
			case "top":
				hasTop = true
			case "mid":
				hasMid = true
			case "bot":
				hasBot = true
			case "sup":
				hasSup = true
			case "jungle":
				hasJungle = true
			default:
				return 0, errors.New("invalid role")
			}
		}
		if !hasTop || !hasMid || !hasBot || !hasSup || !hasJungle {
			return 0, errors.New("invalid roles distribution")
		}

		var users []models.User
		result := common.DB.Select("id", "mmr", "roles", "epoch_position").Find(&users, userIds)
		if result.Error != nil {
			return 0, errors.New(fmt.Sprint("DB error:", result.Error.Error()))
		}
		if result.RowsAffected != 5 {
			return 0, errors.New("invalid team users data")
		}

		var mmrs []int
		for _, user := range users {
			// Calculate fairness score
			if teamIndex == 0 {
				rolesMMRDifferences += user.MMR
			} else {
				rolesMMRDifferences -= user.MMR
			}

			mmrs = append(mmrs, user.MMR)
			userData := teamData.Users[user.ID]
			userData.MMR = user.MMR
			teamData.Users[user.ID] = userData

			// Calculate roles score
			var preferredRoles []string
			err := json.Unmarshal([]byte(user.Roles), &preferredRoles)
			if err != nil {
				return 0, errors.New("couldn't unmarshal roles json")
			}
			roleIndex := common.Find(
				len(preferredRoles), func(i int) bool { return preferredRoles[i] == userData.Role },
			)
			if roleIndex == -1 {
				return 0, errors.New("unknown role")
			}
			switch roleIndex {
			case 0:
				score += 3
			case 1:
				score += 5
			case 2:
				score += 8
			case 3:
				score += 13
			case 4:
				score += 21
			}

			// Calculate time score

		}
		teamData.MMRMedian = common.Median(mmrs)
		teams[teamIndex] = teamData
	}

	score += math.Abs(teams[0].MMRMedian-teams[1].MMRMedian) + math.Abs(float64(rolesMMRDifferences))

	return score, nil
}
