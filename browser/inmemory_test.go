package browser

import (
	"fmt"
	"testing"

	"emailaddress.horse/thousand/app/models"
)

func Test_InMemory_SetName(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	nameFieldSelector := `input[name="name"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(nameFieldSelector),
		bt.SendKeys(nameFieldSelector, "Gruffudd"),
		bt.Submit(nameFieldSelector),
		bt.Text(`#details`).Equals("Gruffudd"),
	)
}

func Test_InMemory_AddExperience(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	expectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	experienceFieldSelector := `input[name="experience"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, expectedExperience),
		bt.Submit(experienceFieldSelector),
		bt.Text("#memories").Contains(expectedExperience),
	)
}

func Test_InMemory_ForgetExperience(t *testing.T) {
	t.Parallel()

	notExpectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.OldVampire{
		Memories: []*models.OldMemory{
			{
				ID:          1,
				Experiences: []models.OldExperience{models.OldExperience(notExpectedExperience)},
			},
		},
	}

	memoryForgetSelector := `#memories input[name="_method"][value="DELETE"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(memoryForgetSelector),
		bt.Submit(memoryForgetSelector),
		bt.InnerHTML("#memories").Not().Contains(notExpectedExperience),
	)
}

func Test_InMemory_AddSkill(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	expectedSkill := "Basic agricultural practices"

	skillFieldSelector := `#skills input[name="description"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(skillFieldSelector),
		bt.SendKeys(skillFieldSelector, expectedSkill),
		bt.Submit(skillFieldSelector),
		bt.Text("#skills").Contains(expectedSkill),
	)
}

func Test_InMemory_CheckSkill(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.OldVampire{
		Skills: []*models.OldSkill{
			{
				ID:          1,
				Description: "Basic agricultural practices",
			},
		},
	}

	expectedSkill := "<s>Basic agricultural practices</s>"

	skillCheckSelector := `input[name="checked"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(skillCheckSelector),
		bt.Submit(skillCheckSelector),
		bt.InnerHTML("#skills").Contains(expectedSkill),
	)
}

func Test_InMemory_AddResource(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	expectedResource := "Calweddyn Farm, rich but challenging soils"

	descriptionFieldSelector := `#resources input[name="description"]`
	stationaryFieldSelector := `#resources input[name="stationary"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(descriptionFieldSelector),
		bt.SendKeys(descriptionFieldSelector, expectedResource),
		bt.Click(stationaryFieldSelector),
		bt.Submit(descriptionFieldSelector),
		bt.Text("#resources").Contains(expectedResource+" (Stationary)"),
	)
}

func Test_InMemory_LoseResource(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.OldVampire{
		Resources: []*models.OldResource{
			{
				ID:          1,
				Description: "Calweddyn Farm, rich but challenging soils",
				Stationary:  true,
			},
		},
	}

	expectedResource := "<s>Calweddyn Farm, rich but challenging soils (Stationary)</s>"

	resourceLoseSelector := `input[name="lost"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(resourceLoseSelector),
		bt.Submit(resourceLoseSelector),
		bt.InnerHTML("#resources").Contains(expectedResource),
	)
}

func Test_InMemory_AddCharacter(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	name := "Lord Othian"
	descriptor := "English gentry visiting a cathedral in St. Davids"

	expectedCharacter := fmt.Sprintf("%s, %s. (Immortal)", name, descriptor)

	nameFieldSelector := `#characters input[name="name"]`
	descriptorFieldSelector := `#characters input[name="descriptor[]"]`
	typeFieldSelector := `#characters input[name="type"][value="immortal"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(nameFieldSelector),
		bt.SendKeys(nameFieldSelector, name),
		bt.SendKeys(descriptorFieldSelector, descriptor),
		bt.Click(typeFieldSelector),
		bt.Submit(nameFieldSelector),
		bt.Text("#characters").Contains(expectedCharacter),
	)
}

func Test_InMemory_AddCharacterDescriptor(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.OldVampire{
		Characters: []*models.OldCharacter{
			{
				ID:   1,
				Name: "Lord Othian",
				Descriptors: []string{
					"English gentry visiting a cathedral in St. Davids",
				},
				Type:     "immortal",
				Deceased: false,
			},
		},
	}

	descriptor := "brought violence upon me as I fled wolves in the forest"

	expectedCharacter := "Lord Othian, English gentry visiting a cathedral in St. Davids, brought violence upon me as I fled wolves in the forest. (Immortal)"

	descriptorFieldSelector := `#character-1 input[name="descriptor"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(descriptorFieldSelector),
		bt.SendKeys(descriptorFieldSelector, descriptor),
		bt.Click(descriptorFieldSelector),
		bt.Submit(descriptorFieldSelector),
		bt.Text("#characters").Contains(expectedCharacter),
	)
}

func Test_InMemory_KillCharacter(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.OldVampire{
		Characters: []*models.OldCharacter{
			{
				ID:       1,
				Name:     "Lord Othian",
				Deceased: false,
			},
		},
	}

	notExpectedCharacter := "Lord Othian"

	characterDeceasedSelector := `input[name="deceased"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(characterDeceasedSelector),
		bt.Submit(characterDeceasedSelector),
		bt.InnerHTML("#characters").Not().Contains(notExpectedCharacter),
	)
}

func Test_InMemory_AddMark(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	expectedMark := "Muddy footprints, muddy handprints, muddy sheets - I leave a trail of dirt wherever I travel."

	descriptionFieldSelector := `#marks input[name="description"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(descriptionFieldSelector),
		bt.SendKeys(descriptionFieldSelector, expectedMark),
		bt.Submit(descriptionFieldSelector),
		bt.Text("#marks").Contains(expectedMark),
	)
}
