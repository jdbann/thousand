package browser

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"emailaddress.horse/thousand/app/models"
)

func TestAppTitle(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var title string

	bt.Run(
		bt.Navigate("/"),
		bt.Text("h1", &title),
	)

	if title != "Thousand" {
		t.Errorf("expected %q; got %q", "Thousand", title)
	}
}

func TestShowVampires(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	if _, err := bt.Models().CreateVampire(context.Background(), "Gruffudd"); err != nil {
		t.Fatal(err)
	}

	var vampires string

	bt.Run(
		bt.Navigate("/vampires"),
		bt.WaitVisible(`#vampires`),
		bt.Text(`#vampires`, &vampires),
	)

	if !strings.Contains(vampires, "Gruffudd") {
		t.Errorf("expected %q to contain %q", "Gruffudd", vampires)
	}
}

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var name string

	newLinkSelector := `#newVampire a`
	nameFieldSelector := `#newVampire input[name="name"]`

	bt.Run(
		bt.Navigate("/vampires"),
		bt.WaitVisible(newLinkSelector),
		bt.Click(newLinkSelector),
		bt.WaitVisible(nameFieldSelector),
		bt.SendKeys(nameFieldSelector, "Gruffudd"),
		bt.Submit(nameFieldSelector),
		bt.Text(`#details`, &name),
	)

	if strings.TrimSpace(name) != "Gruffudd" {
		t.Errorf("expected %q; got %q", "Gruffudd", name)
	}
}

func TestSetName(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var name string

	nameFieldSelector := `input[name="name"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(nameFieldSelector),
		bt.SendKeys(nameFieldSelector, "Gruffudd"),
		bt.Submit(nameFieldSelector),
		bt.Text(`#details`, &name),
	)

	if strings.TrimSpace(name) != "Gruffudd" {
		t.Errorf("expected %q; got %q", "Gruffudd", name)
	}
}

func TestAddExperience(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var memories string
	expectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	experienceFieldSelector := `input[name="experience"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, expectedExperience),
		bt.Submit(experienceFieldSelector),
		bt.Text("#memories", &memories),
	)

	if !strings.Contains(memories, expectedExperience) {
		t.Errorf("expected %q to contain %q", memories, expectedExperience)
	}
}

func TestForgetExperience(t *testing.T) {
	t.Parallel()

	notExpectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.Vampire{
		Memories: []*models.Memory{
			{
				ID:          1,
				Experiences: []models.Experience{models.Experience(notExpectedExperience)},
			},
		},
	}

	var memories string

	memoryForgetSelector := `#memories input[name="_method"][value="DELETE"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(memoryForgetSelector),
		bt.Submit(memoryForgetSelector),
		bt.InnerHTML("#memories", &memories),
	)

	if strings.Contains(memories, notExpectedExperience) {
		t.Errorf("expected %q not to contain %q", memories, notExpectedExperience)
	}
}

func TestAddSkill(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var skills string
	expectedSkill := "Basic agricultural practices"

	skillFieldSelector := `#skills input[name="description"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(skillFieldSelector),
		bt.SendKeys(skillFieldSelector, expectedSkill),
		bt.Submit(skillFieldSelector),
		bt.Text("#skills", &skills),
	)

	if !strings.Contains(skills, expectedSkill) {
		t.Errorf("expected %q to contain %q", skills, expectedSkill)
	}
}

func TestCheckSkill(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.Vampire{
		Skills: []*models.Skill{
			{
				ID:          1,
				Description: "Basic agricultural practices",
			},
		},
	}

	var skills string
	expectedSkill := "<s>Basic agricultural practices</s>"

	skillCheckSelector := `input[name="checked"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(skillCheckSelector),
		bt.Submit(skillCheckSelector),
		bt.InnerHTML("#skills", &skills),
	)

	if !strings.Contains(skills, expectedSkill) {
		t.Errorf("expected %q to contain %q", skills, expectedSkill)
	}
}

func TestAddResource(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var resources string
	expectedResource := "Calweddyn Farm, rich but challenging soils"

	descriptionFieldSelector := `#resources input[name="description"]`
	stationaryFieldSelector := `#resources input[name="stationary"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(descriptionFieldSelector),
		bt.SendKeys(descriptionFieldSelector, expectedResource),
		bt.Click(stationaryFieldSelector),
		bt.Submit(descriptionFieldSelector),
		bt.Text("#resources", &resources),
	)

	if !strings.Contains(resources, expectedResource+" (Stationary)") {
		t.Errorf("expected %q to contain %q", resources, expectedResource)
	}
}

func TestLoseResource(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.Vampire{
		Resources: []*models.Resource{
			{
				ID:          1,
				Description: "Calweddyn Farm, rich but challenging soils",
				Stationary:  true,
			},
		},
	}

	var resources string
	expectedResource := "<s>Calweddyn Farm, rich but challenging soils (Stationary)</s>"

	resourceLoseSelector := `input[name="lost"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(resourceLoseSelector),
		bt.Submit(resourceLoseSelector),
		bt.InnerHTML("#resources", &resources),
	)

	if !strings.Contains(resources, expectedResource) {
		t.Errorf("expected %q to contain %q", resources, expectedResource)
	}
}

func TestAddCharacter(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var characters string
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
		bt.Text("#characters", &characters),
	)

	if !strings.Contains(characters, expectedCharacter) {
		t.Errorf("expected %q to contain %q", characters, expectedCharacter)
	}
}

func TestAddCharacterDescriptor(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.Vampire{
		Characters: []*models.Character{
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

	var characters string
	descriptor := "brought violence upon me as I fled wolves in the forest"

	expectedCharacter := "Lord Othian, English gentry visiting a cathedral in St. Davids, brought violence upon me as I fled wolves in the forest. (Immortal)"

	descriptorFieldSelector := `#character-1 input[name="descriptor"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(descriptorFieldSelector),
		bt.SendKeys(descriptorFieldSelector, descriptor),
		bt.Click(descriptorFieldSelector),
		bt.Submit(descriptorFieldSelector),
		bt.Text("#characters", &characters),
	)

	if !strings.Contains(characters, expectedCharacter) {
		t.Errorf("expected %q to contain %q", characters, expectedCharacter)
	}
}

func TestKillCharacter(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)
	bt.app.Vampire = &models.Vampire{
		Characters: []*models.Character{
			{
				ID:       1,
				Name:     "Lord Othian",
				Deceased: false,
			},
		},
	}

	var characters string
	notExpectedCharacter := "Lord Othian"

	characterDeceasedSelector := `input[name="deceased"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitReady(characterDeceasedSelector),
		bt.Submit(characterDeceasedSelector),
		bt.InnerHTML("#characters", &characters),
	)

	if strings.Contains(characters, notExpectedCharacter) {
		t.Errorf("expected %q not to contain %q", characters, notExpectedCharacter)
	}
}

func TestAddMark(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	var marks string
	expectedMark := "Muddy footprints, muddy handprints, muddy sheets - I leave a trail of dirt wherever I travel."

	descriptionFieldSelector := `#marks input[name="description"]`

	bt.Run(
		bt.Navigate("/"),
		bt.WaitVisible(descriptionFieldSelector),
		bt.SendKeys(descriptionFieldSelector, expectedMark),
		bt.Submit(descriptionFieldSelector),
		bt.Text("#marks", &marks),
	)

	if !strings.Contains(marks, expectedMark) {
		t.Errorf("expected %q to contain %q", marks, expectedMark)
	}
}
