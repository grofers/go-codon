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
	pongo2 "github.com/flosch/pongo2"
	shared "github.com/grofers/go-codon/shared"
	conv "github.com/cstockton/go-conv"
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
	ErrorOutput	interface{}						`yaml:"output-on-error"`
	References	map[string]map[string]string	`yaml:"references"`
}

type Task struct {
	Action				string					`yaml:"action"`
	Input				map[string]string		`yaml:"input"`
	PublishRaw			interface{}				`yaml:"publish"`
	Publish				map[string]string		`yaml:"-"`
	PublishList			[]PublishObj			`yaml:"-"`
	ErrorPublishRaw		interface{}				`yaml:"publish-on-error"`
	ErrorPublish		map[string]string		`yaml:"-"`
	ErrorPublishList	[]PublishObj			`yaml:"-"`
	OnError				[]map[string]string		`yaml:"on-error"`
	OnErrorList			[]TodoObj				`yaml:"-"`
	OnSuccess			[]map[string]string		`yaml:"on-success"`
	OnSuccessList		[]TodoObj				`yaml:"-"`
	OnComplete			[]map[string]string		`yaml:"on-complete"`
	OnCompleteList		[]TodoObj				`yaml:"-"`
	Timeout				int64					`yaml:"timeout"`
	WithItems			string					`yaml:"with-items"`
	Loop				LoopInfo				`yaml:"loop"`
}

type LoopInfo struct {
	TaskName			string				`yaml:"task"`
	Input				map[string]string	`yaml:"input"`
	PublishRaw			interface{}			`yaml:"publish"`
	Publish				map[string]string	`yaml:"-"`
	PublishList			[]PublishObj		`yaml:"-"`
	ErrorPublishRaw		interface{}			`yaml:"publish-on-error"`
	ErrorPublish		map[string]string	`yaml:"-"`
	ErrorPublishList	[]PublishObj		`yaml:"-"`
}

type TodoObj struct {
	TaskName		string
	ExpressionName	string
	Srno			int
}

type PublishObj struct {
	VariableName	string
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

	err = s.setLists()
	if err != nil {
		return PostSpec {}, err
	}

	ps.ExpressionMap, err = s.getExpressionMap()
	if err != nil {
		return PostSpec {}, err
	}

	ps.ActionMap, err = s.getActionMap()
	if err != nil {
		return PostSpec {}, err
	}

	return ps, nil
}

func (s *Spec) setLists() error {
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
		task_obj.PublishList, err = task_obj.getPublishList()
		if err != nil {
			return err
		}
		task_obj.Publish = task_obj.tryPublishMap()
		task_obj.ErrorPublishList, err = task_obj.getErrorPublishList()
		if err != nil {
			return err
		}
		task_obj.ErrorPublish = task_obj.tryErrorPublishMap()
		if task_obj.WithItems != "" {
			task_obj.Loop.PublishList, err = task_obj.Loop.getPublishList()
			if err != nil {
				return err
			}
			task_obj.Loop.Publish = task_obj.Loop.tryPublishMap()
			task_obj.Loop.ErrorPublishList, err = task_obj.Loop.getErrorPublishList()
			if err != nil {
				return err
			}
			task_obj.Loop.ErrorPublish = task_obj.Loop.tryErrorPublishMap()
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

func (s *Spec) appendExpression(all_exprs map[string]Expression, expr *string, counter *int) error {
	if _, ok := all_exprs[*expr]; ok {
		return nil
	}
	expr_obj, err := s.processExpression(*expr)
	if err != nil {
		log.Println("Unable to process expression:", *expr)
		return err
	}
	expr_obj.Srno = *counter
	(*counter)++
	all_exprs[*expr] = expr_obj
	return nil
}

func (s *Spec) getExpressionMap() (map[string]Expression, error) {
	all_exprs := make(map[string]Expression)
	counter := 1

	for _, task := range s.Tasks {
		for _, expr := range task.Input {
			err := s.appendExpression(all_exprs, &expr, &counter)
			if err != nil {
				return nil, err
			}
		}
		for _, publish_obj := range task.PublishList {
			expr := publish_obj.ExpressionName
			err := s.appendExpression(all_exprs, &expr, &counter)
			if err != nil {
				return nil, err
			}
		}
		for _, publish_obj := range task.ErrorPublishList {
			expr := publish_obj.ExpressionName
			err := s.appendExpression(all_exprs, &expr, &counter)
			if err != nil {
				return nil, err
			}
		}
		for _, task_obj := range task.OnErrorList {
			expr := task_obj.ExpressionName
			err := s.appendExpression(all_exprs, &expr, &counter)
			if err != nil {
				return nil, err
			}
		}
		for _, task_obj := range task.OnSuccessList {
			expr := task_obj.ExpressionName
			err := s.appendExpression(all_exprs, &expr, &counter)
			if err != nil {
				return nil, err
			}
		}
		for _, task_obj := range task.OnCompleteList {
			expr := task_obj.ExpressionName
			err := s.appendExpression(all_exprs, &expr, &counter)
			if err != nil {
				return nil, err
			}
		}
		if task.WithItems != "" {
			expr := task.WithItems
			err := s.appendExpression(all_exprs, &expr, &counter)
			if err != nil {
				return nil, err
			}
			for _, expr := range task.Loop.Input {
				err := s.appendExpression(all_exprs, &expr, &counter)
				if err != nil {
					return nil, err
				}
			}
			for _, task_obj := range task.Loop.PublishList {
				expr := task_obj.ExpressionName
				err := s.appendExpression(all_exprs, &expr, &counter)
				if err != nil {
					return nil, err
				}
			}
			for _, task_obj := range task.Loop.ErrorPublishList {
				expr := task_obj.ExpressionName
				err := s.appendExpression(all_exprs, &expr, &counter)
				if err != nil {
					return nil, err
				}
			}
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

	exprs = make([]string, 0)
	exprs = extractExpressionsFromMap(s.ErrorOutput, exprs)
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

func (t Task) getPublishList() (publish_list []PublishObj, err error) {
	publish_list, err = createPublishList(t.PublishRaw)
	return
}

func (t Task) tryPublishMap() (publish_map map[string]string) {
	publish_map = createPublishMap(t.PublishList)
	return
}

func (t Task) getErrorPublishList() (publish_list []PublishObj, err error) {
	publish_list, err = createPublishList(t.ErrorPublishRaw)
	return
}

func (t Task) tryErrorPublishMap() (publish_map map[string]string) {
	publish_map = createPublishMap(t.ErrorPublishList)
	return
}

func (t LoopInfo) getPublishList() (publish_list []PublishObj, err error) {
	publish_list, err = createPublishList(t.PublishRaw)
	return
}

func (t LoopInfo) tryPublishMap() (publish_map map[string]string) {
	publish_map = createPublishMap(t.PublishList)
	return
}

func (t LoopInfo) getErrorPublishList() (publish_list []PublishObj, err error) {
	publish_list, err = createPublishList(t.ErrorPublishRaw)
	return
}

func (t LoopInfo) tryErrorPublishMap() (publish_map map[string]string) {
	publish_map = createPublishMap(t.ErrorPublishList)
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
		// Just testing compile
		case "jmes":
			_, err := jmespath.Compile(ret_expr.Raw)
			if err != nil {
				return Expression{}, err
			}
		case "jngo":
			_, err := pongo2.FromString(ret_expr.Raw)
			if err != nil {
				return Expression{}, err
			}
		default:
			return Expression{}, fmt.Errorf("Expression of this type is not supported: %v", expr)
		}
		return ret_expr, nil
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

func createPublishMap(publish_list []PublishObj) map[string]string {
	publish_map := make(map[string]string)
	for _, publish_obj := range publish_list {
		publish_map[publish_obj.VariableName] = publish_obj.ExpressionName
	}
	return publish_map
}

func createPublishList(map_list interface{}) ([]PublishObj, error) {
	if map_list == nil {
		return []PublishObj {}, nil
	}

	counter := 1
	switch map_list_v := map_list.(type) {
	case map[interface{}]interface{}:
		retlist := make([]PublishObj, len(map_list_v))
		for ct, ce_i := range map_list_v {
			ce, err := conv.String(ce_i)
			if err != nil {
				return nil, fmt.Errorf("Invalid expression value in publis list: %v", map_list_v)
			}
			retlist[counter-1] = PublishObj{
				Srno: counter,
				VariableName: ct.(string),
				ExpressionName: ce,
			}
			counter++
		}
		return retlist, nil
	case []interface{}:
		retlist := make([]PublishObj, len(map_list_v))
		for _, task_map_i := range map_list_v {
			task_map := task_map_i.(map[interface{}]interface{})
			if len(task_map) != 1 {
				return nil, fmt.Errorf("Each entry in todo must have only one key-value pair: %v", map_list_v)
			}
			for ct, ce_i := range task_map {
				ce, err := conv.String(ce_i)
				if err != nil {
					return nil, fmt.Errorf("Invalid expression value in publis list: %v", map_list_v)
				}
				retlist[counter-1] = PublishObj{
					Srno: counter,
					VariableName: ct.(string),
					ExpressionName: ce,
				}
				counter++
				break
			}
		}
		return retlist, nil
	default:
		return nil, fmt.Errorf("Publish list must be a list of variables: %v, Type: %T", map_list, map_list)
	}
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
