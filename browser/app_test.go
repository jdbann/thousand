package browser

import (
	"context"
	"fmt"
	"testing"
)

func TestAppTitle(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	bt.Run(
		bt.Navigate("/"),
		bt.Text("h1").Equals("Thousand"),
	)
}

func TestShowVampires(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	if _, err := bt.Models().CreateVampire(context.Background(), "Gruffudd"); err != nil {
		t.Fatal(err)
	}

	bt.Run(
		bt.Navigate("/vampires"),
		bt.WaitVisible(`#vampires`),
		bt.Text(`#vampires`).Contains("Gruffudd"),
	)
}

func TestCreateVampire(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	newLinkSelector := `#newVampire a`
	nameFieldSelector := `#newVampire input[name="name"]`

	bt.Run(
		bt.Navigate("/vampires"),
		bt.WaitVisible(newLinkSelector),
		bt.Click(newLinkSelector),
		bt.WaitVisible(nameFieldSelector),
		bt.SendKeys(nameFieldSelector, "Gruffudd"),
		bt.Submit(nameFieldSelector),
		bt.Text(`#details`).Equals("Gruffudd"),
	)
}

func TestAddExperience(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	newExperienceLinkSelector := `#memories a[href$="/experiences/new"]`
	experienceFieldSelector := `#memories input[name="description"]`
	expectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, expectedExperience),
		bt.Submit(experienceFieldSelector),
		bt.Text(`#memories`).Contains(expectedExperience),
	)
}

func TestAddExperienceFormDismisses(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	memory := vampire.Memories[0]

	newExperienceLinkSelector := fmt.Sprintf(`#memory-%s a[href$="/experiences/new"]`, memory.ID.String())
	experienceFieldSelector := fmt.Sprintf(`#memory-%s input[name="description"]`, memory.ID.String())

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),

		// Clicking outside the form should dismiss it
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.Click("h1"),
		bt.WaitNotPresent(experienceFieldSelector),

		// Clicking outside the form should not dismiss it if user has started
		// writing
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, "Sta"),
		bt.Click("h1"),
		bt.WaitVisible(experienceFieldSelector),

		// Clicking outside the form should dismiss it if user has cleared input
		bt.SendKeys(experienceFieldSelector, "\b\b\b"), // \b => Backspace
		bt.Click("h1"),
		bt.WaitNotPresent(experienceFieldSelector),
	)
}

func TestCannotAddFourExperiences(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	memory := vampire.Memories[0]

	if _, err := bt.Models().CreateExperience(context.Background(), vampire.ID, memory.ID, "Experience #1"); err != nil {
		t.Fatal(err)
	}

	if _, err := bt.Models().CreateExperience(context.Background(), vampire.ID, memory.ID, "Experience #2"); err != nil {
		t.Fatal(err)
	}

	newExperienceLinkSelector := fmt.Sprintf(`#memory-%s a[href$="/experiences/new"]`, memory.ID.String())
	experienceFieldSelector := fmt.Sprintf(`#memory-%s input[name="description"]`, memory.ID.String())
	expectedExperience := "I am Gruffudd, a Welsh farmer in the valleys of Pembroke; I am a recluse, fond of nature and withdrawn from the village."

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newExperienceLinkSelector),
		bt.Click(newExperienceLinkSelector),
		bt.WaitVisible(experienceFieldSelector),
		bt.SendKeys(experienceFieldSelector, expectedExperience),
		bt.Submit(experienceFieldSelector),
		bt.Text(fmt.Sprintf("#memory-%s", memory.ID.String())).Contains(expectedExperience),
		bt.Text(fmt.Sprintf("#memory-%s", memory.ID.String())).Not().Contains("New Experience"),
	)
}

func TestAddSkill(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	newSkillLinkSelector := `#skills a[href$="/skills/new"]`
	skillFieldSelector := `#skills input[name="description"]`
	expectedSkill := "Navigating forests"

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newSkillLinkSelector),
		bt.Click(newSkillLinkSelector),
		bt.WaitVisible(skillFieldSelector),
		bt.SendKeys(skillFieldSelector, expectedSkill),
		bt.Submit(skillFieldSelector),
		bt.Text(`#skills`).Contains(expectedSkill),
	)
}

func TestAddResource(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	newResourceLinkSelector := `#resources a[href$="/resources/new"]`
	resourceDescriptionSelector := `#resources input[name="description"]`
	resourceStationarySelector := `#resources input[name="stationary"]`
	expectedResource := "Calweddyn Farm, rich but challenging soils"

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newResourceLinkSelector),
		bt.Click(newResourceLinkSelector),
		bt.WaitVisible(resourceDescriptionSelector),
		bt.SendKeys(resourceDescriptionSelector, expectedResource),
		bt.Click(resourceStationarySelector),
		bt.Submit(resourceDescriptionSelector),
		bt.Text(`#resources`).Contains(expectedResource),
		bt.Text(`#resources`).Contains("Stationary"),
	)
}

func TestAddCharacter(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	newCharacterLinkSelector := `#characters a[href$="/characters/new"]`
	characterNameSelector := `#characters input[name="name"]`
	characterTypeSelector := `#characters select[name="type"]`
	expectedCharacter := "Lord Othian, English gentry visiting a cathedral in St. Davids."

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newCharacterLinkSelector),
		bt.Click(newCharacterLinkSelector),
		bt.WaitVisible(characterNameSelector),
		bt.SendKeys(characterNameSelector, expectedCharacter),
		bt.SendKeys(characterTypeSelector, "Immortal"),
		bt.Submit(characterNameSelector),
		bt.Text(`#characters`).Contains(expectedCharacter),
		bt.Text(`#characters`).Contains("Immortal"),
	)
}

func TestAddMark(t *testing.T) {
	t.Parallel()

	bt := NewBrowserTest(t)

	vampire, err := bt.Models().CreateVampire(context.Background(), "Gruffudd")
	if err != nil {
		t.Fatal(err)
	}

	newMarkLinkSelector := `#marks a[href$="/marks/new"]`
	markDescriptionSelector := `#marks input[name="description"]`
	expectedMark := "Muddy footprints, muddy handprints, muddy sheets - I leave a trail of dirt wherever I travel."

	bt.Run(
		bt.Navigate(fmt.Sprintf("/vampires/%s", vampire.ID.String())),
		bt.WaitForTurbo(),
		bt.WaitVisible(newMarkLinkSelector),
		bt.Click(newMarkLinkSelector),
		bt.WaitVisible(markDescriptionSelector),
		bt.SendKeys(markDescriptionSelector, expectedMark),
		bt.Submit(markDescriptionSelector),
		bt.Text(`#marks`).Contains(expectedMark),
	)
}
