package lpg2

type ConfigurationStack struct {
	TABLE_SIZE             int
	table                  []*ConfigurationElement
	configuration_stack    *ObjectTuple
	state_root             *StateElement
	max_configuration_size int
	stacks_size            int
	state_element_size     int
	prs                    ParseTable
}

func NewConfigurationStack(prs ParseTable) *ConfigurationStack {
	a := new(ConfigurationStack)
	a.TABLE_SIZE = 1021 // 1021 is a prime
	a.prs = prs
	a.state_element_size = 1
	a.state_root = NewStateElement()

	a.state_root.number = prs.GetStartState()

	a.table = make([]*ConfigurationElement, a.TABLE_SIZE, a.TABLE_SIZE)
	a.configuration_stack = NewObjectTupleWithEstimate(1 << 12)
	a.max_configuration_size = 0
	a.stacks_size = 0
	return a
}

func (my *ConfigurationStack) MakeStateList(parent *StateElement, stack []int, index int, stack_top int) *StateElement {
	var i int
	for i = index; i <= stack_top; i++ {

		my.state_element_size++

		var state = NewStateElement()

		state.number = stack[i]
		state.parent = parent

		//state.children = undefined
		//state.siblings = undefined

		parent.children = state
		parent = state
	}
	return parent
}
func (my *ConfigurationStack) FindOrInsertStack(root *StateElement, stack []int, index int, stack_top int) *StateElement {
	var state_number = stack[index]
	var p *StateElement
	for p = root; p != nil; p = p.siblings {

		if p.number == state_number {
			if index == stack_top {
				return p
			} else {
				if p.children == nil {
					return my.MakeStateList(p, stack, index+1, stack_top)
				} else {
					my.FindOrInsertStack(p.children, stack, index+1, stack_top)
				}
			}
		}
	}

	my.state_element_size++

	var node = NewStateElement()
	node.number = state_number
	node.parent = root.parent
	node.children = nil
	node.siblings = root.siblings
	root.siblings = node

	if index == stack_top {
		return node
	} else {
		return my.MakeStateList(node, stack, index+1, stack_top)
	}
}

func (my *ConfigurationStack) FindConfiguration(stack []int, stack_top int, curtok int) bool {

	var last_element = my.FindOrInsertStack(my.state_root, stack, 0, stack_top)
	var hash_address = curtok % my.TABLE_SIZE
	var configuration *ConfigurationElement
	for configuration = my.table[hash_address]; configuration != nil; configuration = configuration.next {
		if configuration.curtok == curtok && last_element == configuration.last_element {
			return true
		}
	}
	return false
}

func (my *ConfigurationStack) Push(stack []int, stack_top int, conflict_index int, curtok int, action_length int) {

	var configuration = NewConfigurationElement()
	var hash_address = curtok % my.TABLE_SIZE

	configuration.next = my.table[hash_address]

	my.table[hash_address] = configuration
	my.max_configuration_size++ // keep track of int of configurations

	configuration.stack_top = stack_top
	my.stacks_size += stack_top + 1 // keep track of int of stack elements processed
	configuration.last_element = my.FindOrInsertStack(my.state_root, stack, 0, stack_top)
	configuration.conflict_index = conflict_index
	configuration.curtok = curtok
	configuration.action_length = action_length

	my.configuration_stack.Add(configuration)
	return
}
func (my *ConfigurationStack) cast(c interface{}) *ConfigurationElement {
	configuration, ok := c.(*ConfigurationElement)
	if ok {
		return configuration
	}
	return nil
}
func (my *ConfigurationStack) Pop() *ConfigurationElement {
	var configuration *ConfigurationElement = nil
	if my.configuration_stack.Size() > 0 {
		var index = my.configuration_stack.Size() - 1
		configuration = my.cast(my.configuration_stack.Get(index))
		if configuration != nil {
			configuration.act = my.prs.BaseAction(configuration.conflict_index)
			configuration.conflict_index += 1

			if my.prs.BaseAction(configuration.conflict_index) == 0 {
				my.configuration_stack.ResetTo(index)
			}
		} else {
			return nil
		}

	}

	return configuration
}
func (my *ConfigurationStack) Top() *ConfigurationElement {
	var configuration *ConfigurationElement = nil
	if my.configuration_stack.Size() > 0 {
		var index = my.configuration_stack.Size() - 1
		var configuration = my.cast(my.configuration_stack.Get(index))
		configuration.act = my.prs.BaseAction(configuration.conflict_index)
	}

	return configuration
}
func (my *ConfigurationStack) size() int {
	return my.configuration_stack.Size()
}
func (my *ConfigurationStack) MaxConfigurationSize() int {
	return my.max_configuration_size
}
func (my *ConfigurationStack) NumStateElements() int {
	return my.state_element_size
}
func (my *ConfigurationStack) StacksSize() int {
	return my.stacks_size
}
