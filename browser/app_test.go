package browser

import (
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
	bt.app.Character = &models.Character{
		Skills: []models.Skill{
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
	bt.app.Character = &models.Character{
		Resources: []models.Resource{
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
