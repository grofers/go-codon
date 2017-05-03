package shared

import (
	"log"
	"fmt"
	"strings"
	"encoding/json"
	"regexp"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	jmespath "github.com/jmespath/go-jmespath"
	shared "github.com/grofers/go-codon/shared"
)

type PostSpec struct {
	//The original spec read from yaml
	OrigSpec		*Spec
	// Compiled list of all the expressions in spec
	ExpressionMap	map[string]Expression
	// Compiles list of actions
	ActionMap		map[string]Action
	//Variables added by language level processing
	LanguageSpec	map[string]interface{}
}

type Spec struct {
	Name		string							`yaml:"name"`
	Start		[]string						`yaml:"start"`
	Tasks		map[string]Task					`yaml:"tasks"`
	Output		interface{}						`yaml:"output"`
	ErrorOutput	map[string]string				`yaml:"output-on-error"`
	References	map[string]map[string]string	`yaml:"references"`
}

type Task struct {
	Action			string						`yaml:"action"`
	Input			map[string]string			`yaml:"input"`
	Publish			map[string]string			`yaml:"publish"`
	ErrorPublish	map[string]string			`yaml:"publish-on-error"`
	OnError			[]map[string]string			`yaml:"on-error"`
	OnErrorList		[]TodoObj					`yaml:"-"`
	OnSuccess		[]map[string]string			`yaml:"on-success"`
	OnSuccessList	[]TodoObj					`yaml:"-"`
	OnComplete		[]map[string]string			`yaml:"on-complete"`
	OnCompleteList	[]TodoObj					`yaml:"-"`
	Timeout			int64						`yaml:"timeout"`
}

type TodoObj struct {
	TaskName		string
	ExpressionName	string
	Srno			int
}

type Expression struct {
	Type			string
	Raw				string
	Srno			int
}

type Action struct {
	Type			string
	Raw				string
	Pascalized		string
}

func ReadSpec(filename string) (Spec, error) {
	filename_data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Spec{}, err
	}

	var t Spec

	err = yaml.Unmarshal([]byte(filename_data), &t)
	if err != nil {
		return Spec{}, err
	}

	return t, nil
}

func (s *Spec) Process() (PostSpec, error) {
	var err error
	ps := PostSpec {}
	ps.OrigSpec = s
	if s.References == nil {
		s.References = make(map[string]map[string]string)
	}
	ps.LanguageSpec = make(map[string]interface{})

	ps.ExpressionMap, err = s.getExpressionMap()
	if err != nil {
		return PostSpec {}, err
	}

	ps.ActionMap, err = s.getActionMap()
	if err != nil {
		return PostSpec {}, err
	}

	err = s.setTodoLists()
	if err != nil {
		return PostSpec {}, err
	}

	return ps, nil
}

func (s *Spec) setTodoLists() error {
	var err error
	for task_name, _ := range s.Tasks {
		task_obj := s.Tasks[task_name]
		task_obj.OnErrorList, err = task_obj.getOnErrorList()
		if err != nil {
			return err
		}
		task_obj.OnSuccessList, err = task_obj.getOnSuccessList()
		if err != nil {
			return err
		}
		task_obj.OnCompleteList, err = task_obj.getOnCompleteList()
		if err != nil {
			return err
		}
		s.Tasks[task_name] = task_obj
	}
	return nil
}

func (s *Spec) getActionMap() (map[string]Action, error) {
	all_actions := make(map[string]Action)

	for _, task_obj := range s.Tasks {
		action_obj, err := processAction(task_obj.Action)
		if err != nil {
			log.Println("Unable to process action:", task_obj.Action)
			return nil, err
		}
		all_actions[task_obj.Action] = action_obj
	}

	return all_actions, nil
}

func (s *Spec) getExpressionMap() (map[string]Expression, error) {
	all_exprs := make(map[string]Expression)
	counter := 1

	for _, task := range s.Tasks {
		for _, expr := range task.Input {
			if _, ok := all_exprs[expr]; ok {
				continue
			}
			expr_obj, err := s.processExpression(expr)
			if err != nil {
				log.Println("Unable to process expression:", expr)
				return nil, err
			}
			expr_obj.Srno = counter
			counter++
			all_exprs[expr] = expr_obj
		}
		for _, expr := range task.Publish {
			if _, ok := all_exprs[expr]; ok {
				continue
			}
			expr_obj, err := s.processExpression(expr)
			if err != nil {
				log.Println("Unable to process expression:", expr)
				return nil, err
			}
			expr_obj.Srno = counter
			counter++
			all_exprs[expr] = expr_obj
		}
		for _, expr := range task.ErrorPublish {
			if _, ok := all_exprs[expr]; ok {
				continue
			}
			expr_obj, err := s.processExpression(expr)
			if err != nil {
				log.Println("Unable to process expression:", expr)
				return nil, err
			}
			expr_obj.Srno = counter
			counter++
			all_exprs[expr] = expr_obj
		}
		for _, error_map := range task.OnError {
			if len(error_map) != 1 {
				return nil, fmt.Errorf("Each entry in on-error must have only one key-value pair: %v", task.OnError)
			}
			var expr string
			for _, ce := range error_map {expr = ce;break;}
			if _, ok := all_exprs[expr]; ok {
				continue
			}
			expr_obj, err := s.processExpression(expr)
			if err != nil {
				log.Println("Unable to process expression:", expr)
				return nil, err
			}
			expr_obj.Srno = counter
			counter++
			all_exprs[expr] = expr_obj
		}
		for _, success_map := range task.OnSuccess {
			if len(success_map) != 1 {
				return nil, fmt.Errorf("Each entry in on-error must have only one key-value pair: %v", task.OnError)
			}
			var expr string
			for _, ce := range success_map {expr = ce;break;}
			if _, ok := all_exprs[expr]; ok {
				continue
			}
			expr_obj, err := s.processExpression(expr)
			if err != nil {
				log.Println("Unable to process expression:", expr)
				return nil, err
			}
			expr_obj.Srno = counter
			counter++
			all_exprs[expr] = expr_obj
		}
		for _, complete_map := range task.OnComplete {
			if len(complete_map) != 1 {
				return nil, fmt.Errorf("Each entry in on-error must have only one key-value pair: %v", task.OnError)
			}
			var expr string
			for _, ce := range complete_map {expr = ce;break;}
			if _, ok := all_exprs[expr]; ok {
				continue
			}
			expr_obj, err := s.processExpression(expr)
			if err != nil {
				log.Println("Unable to process expression:", expr)
				return nil, err
			}
			expr_obj.Srno = counter
			counter++
			all_exprs[expr] = expr_obj
		}
	}

	var exprs []string = make([]string, 0)
	exprs = extractExpressionsFromMap(s.Output, exprs)
	for _, expr := range exprs {
		if _, ok := all_exprs[expr]; ok {
			continue
		}
		expr_obj, err := s.processExpression(expr)
		if err != nil {
			log.Println("Unable to process expression:", expr)
			return nil, err
		}
		expr_obj.Srno = counter
		counter++
		all_exprs[expr] = expr_obj
	}

	return all_exprs, nil
}

func extractExpressionsFromMap(expr_imap interface{}, exprs []string) []string {
	switch expr_map := expr_imap.(type) {
	case map[string]interface{}:
		for _, val := range expr_map {
			switch val_type := val.(type) {
			case map[interface{}]interface{}:
				exprs = extractExpressionsFromMap(val_type, exprs)
			default:
				val_type_str := fmt.Sprintf("%v", val_type)
				exprs = append(exprs, val_type_str)
			}
		}
	case map[interface{}]interface{}:
		for _, val := range expr_map {
			switch val_type := val.(type) {
			case map[interface{}]interface{}:
				exprs = extractExpressionsFromMap(val_type, exprs)
			default:
				val_type_str := fmt.Sprintf("%v", val_type)
				exprs = append(exprs, val_type_str)
			}
		}
	case string:
		val_type_str := fmt.Sprintf("%v", expr_map)
		exprs = append(exprs, val_type_str)
	}
	return exprs
}

func (t Task) getOnErrorList() (todo_list []TodoObj, err error) {
	todo_list, err = createTodoList(t.OnError)
	return
}

func (t Task) getOnSuccessList() (todo_list []TodoObj, err error) {
	todo_list, err = createTodoList(t.OnSuccess)
	return
}

func (t Task) getOnCompleteList() (todo_list []TodoObj, err error) {
	todo_list, err = createTodoList(t.OnComplete)
	return
}

var expr_regex = regexp.MustCompile("<%(?P<type>[^ ]*) (?P<expr>[^ ].*) %>")

func (s *Spec) processExpression(expr string) (Expression, error) {
	type_expr := expr_regex.FindStringSubmatch(expr)
	ret_expr := Expression{}
	if type_expr != nil {
		ret_expr.Type = type_expr[1]
		ret_expr.Raw = type_expr[2]
		if ret_expr.Type == "" {
			ret_expr.Type = "yaql"
		}
		switch ret_expr.Type {
		case "jmes":
			// Just testing compile
			_, err := jmespath.Compile(ret_expr.Raw)
			if err != nil {
				return Expression{}, err
			}
			return ret_expr, nil
		default:
			return Expression{}, fmt.Errorf("Expression of this type is not supported: %v", expr)
		}
	} else {
		ret_expr.Type = "json"
		ret_expr.Raw = expr
		var res interface{}
		err := json.Unmarshal([]byte(expr), &res)
		if err != nil {
			return Expression{}, err
		}
		return ret_expr, nil
	}
}

func processAction(action string) (Action, error) {
	ret_action := Action{}
	
	action_elems := strings.Split(action, ".")
	for i := 1; i < len(action_elems); i++ {
		action_elems[i] = shared.Pascalize(action_elems[i])
	}
	ret_action.Pascalized = strings.Join(action_elems, ".")
	ret_action.Type = action_elems[0]
	ret_action.Raw = action

	return ret_action, nil
}

func createTodoList(map_list []map[string]string) ([]TodoObj, error) {
	counter := 1
	retlist := make([]TodoObj, len(map_list))
	for _, task_map := range map_list {
		if len(task_map) != 1 {
			return nil, fmt.Errorf("Each entry in todo must have only one key-value pair: %v", map_list)
		}
		for ct, ce := range task_map {
			retlist[counter-1] = TodoObj{
				Srno: counter,
				TaskName: ct,
				ExpressionName: ce,
			}
			counter++
			break
		}
	}
	return retlist, nil
}
