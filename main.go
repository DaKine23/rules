package main

import (
	"fmt"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

type Input struct {
	Value  string
	Audits []Audit
}

func (i *Input) NewAudit(when time.Time, what string) {

	i.Audits = append(i.Audits, Audit{When: when, What: what})

}

type Audit struct {
	When time.Time
	What string
}

type ruleEngine struct {
	knowledgeBase    *ast.KnowledgeBase
	knowledgeLibrary *ast.KnowledgeLibrary
	engine           *engine.GruleEngine

	rulesFilePath string
	name          string
	version       string
}

func NewRuleEngine(name, version, filepath string) *ruleEngine {

	re := ruleEngine{}

	//create a knowledge lib and build rules from source file

	re.name = name
	re.rulesFilePath = filepath
	re.version = version

	re.knowledgeLibrary = ast.NewKnowledgeLibrary()

	ruleBuilder := builder.NewRuleBuilder(re.knowledgeLibrary)
	err := ruleBuilder.BuildRuleFromResource(name, version, pkg.NewFileResource(filepath))

	re.knowledgeBase = re.knowledgeLibrary.NewKnowledgeBaseInstance(name, version)

	if err != nil {
		panic(err)
	}
	// create an engine
	re.engine = &engine.GruleEngine{MaxCycle: 50}
	return &re
}

func (re *ruleEngine) Execute(data ast.IDataContext) error {

	return re.engine.Execute(data, re.knowledgeBase)

}

func main() {

	in := &Input{Value: "foo"}
	in.Audits = make([]Audit, 0)
	fmt.Println(in)

	re := NewRuleEngine("PrototypeRules", "0.0.1", "rules.grl")

	//setup a context (data to be processed)
	dataContext := ast.NewDataContext()
	err := dataContext.Add("Input", in)
	if err != nil {
		panic(err)
	}

	err = re.Execute(dataContext)
	if err != nil {
		panic(err)
	}

	fmt.Println(in)
}
