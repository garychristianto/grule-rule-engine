package examples

import (
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/context"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/model"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"testing"
)

const (
	Rule = `
rule UserTestRule3 "test 3"  salience 10{
	when
	  User.GetName() == "Watson"
	then
	  User.SetName("FromRuleScope3");
	  Retract("UserTestRule3");
}
`
)

type AUser struct {
	Name string
}

func (u *AUser) GetName() string {
	return u.Name
}

func (u *AUser) SetName(name string) {
	u.Name = name
}

func TestMethodCall_Issue5(t *testing.T) {
	user := &AUser{
		Name: "Watson",
	}

	dataContext := context.NewDataContext()
	err := dataContext.Add("User", user)
	if err != nil {
		t.Fatal(err)
	}

	knowledgeBase := model.NewKnowledgeBase()
	ruleBuilder := builder.NewRuleBuilder(knowledgeBase)

	err = ruleBuilder.BuildRuleFromResource(pkg.NewBytesResource([]byte(Rule)))
	if err != nil {
		t.Log(err)
	} else {
		eng1 := &engine.GruleEngine{MaxCycle: 5}
		err := eng1.Execute(dataContext, knowledgeBase)
		if err != nil {
			t.Fatal(err)
		}
		if user.GetName() != "FromRuleScope3" {
			t.Errorf("User should be FromRuleScope3 but %s", user.GetName())
		}
	}
}
