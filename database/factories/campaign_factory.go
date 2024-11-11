package factories

import (
	"github.com/fullstackdev42/mp-emailer/campaign"
	"github.com/fullstackdev42/mp-emailer/database"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

type CampaignFactory struct {
	BaseFactory
	UserID string
}

func NewCampaignFactory(db database.Interface, userID string) *CampaignFactory {
	return &CampaignFactory{
		BaseFactory: BaseFactory{DBInterface: db},
		UserID:      userID,
	}
}

func (f *CampaignFactory) Make() interface{} {
	campaignThemes := []string{
		"Climate Action",
		"Social Justice",
		"Environmental Protection",
		"Human Rights",
		"Education Reform",
		"Healthcare Access",
		"Indigenous Rights",
		"Economic Equality",
		"Unmarked Burials",
	}

	suffixes := []string{"Initiative", "Campaign", "Movement", "Action", "Now"}

	themeIndex, err := faker.RandomInt(0, len(campaignThemes)-1)
	if err != nil {
		// Fallback to first theme if there's an error
		themeIndex = []int{0}
	}

	suffixIndex, err := faker.RandomInt(0, len(suffixes)-1)
	if err != nil {
		// Fallback to first suffix if there's an error
		suffixIndex = []int{0}
	}

	name := campaignThemes[themeIndex[0]] + " " + suffixes[suffixIndex[0]]

	ownerID, err := uuid.Parse(f.UserID)
	if err != nil {
		// If parsing fails, generate a new UUID
		ownerID = uuid.New()
	}

	c := &campaign.Campaign{
		Name:        name,
		Description: faker.Sentence(),
		Template:    generateTemplate(),
		OwnerID:     ownerID,
	}

	return c
}

func (f *CampaignFactory) MakeMany(count int) []interface{} {
	var campaigns []interface{}
	for i := 0; i < count; i++ {
		campaigns = append(campaigns, f.Make())
	}
	return campaigns
}

func generateTemplate() string {
	return `
Dear {{.RecipientName}},

I am writing to bring your attention to an important issue that affects our community.

{{.CampaignDescription}}

We urgently need your support to make a difference. Here's how you can help:
1. {{.ActionItem1}}
2. {{.ActionItem2}}
3. {{.ActionItem3}}

Your voice matters. Together, we can create positive change.

Best regards,
{{.SenderName}}
`
}
