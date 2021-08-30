package lpg2

const TABLE_SIZE int = 1021 // 1021 is a prime
type ConfigurationStack struct {
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
	a.prs = prs
	a.state_element_size = 1
	a.state_root = NewStateElement()

	a.state_root.number = prs.getStartState()

	a.table = make([]*ConfigurationElement, TABLE_SIZE, TABLE_SIZE)
	a.max_configuration_size = 0
	a.stacks_size = 0
	return a
}


func (this ConfigurationStack) makeStateList(parent *StateElement, stack []int, index int, stack_top int) *StateElement {
	var i int
	for i = index; i <= stack_top; i++ {

		this.state_element_size++

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
func (this ConfigurationStack) findOrInsertStack(root *StateElement, stack []int, index int, stack_top int) *StateElement {
	var state_number = stack[index]
	var p *StateElement
	for p = root; p != nil; p = p.siblings {

		if p.number == state_number {
			if index == stack_top {
				return p
			}
		} else {
			if p.children == nil {
				return this.makeStateList(p, stack, index+1, stack_top)
			} else {
				this.findOrInsertStack(p.children, stack, index+1, stack_top)
			}
		}
	}

	this.state_element_size++

	var node = NewStateElement()
	node.number = state_number
	node.parent = root.parent
	node.children = nil
	node.siblings = root.siblings
	root.siblings = node

	if index == stack_top {
		return node
	} else {
		return this.makeStateList(node, stack, index+1, stack_top)
	}
}

func (this ConfigurationStack) findConfiguration(stack []int, stack_top int, curtok int) bool {

	var last_element = this.findOrInsertStack(this.state_root, stack, 0, stack_top)
	var hash_address = curtok % TABLE_SIZE
	var configuration *ConfigurationElement
	for configuration = this.table[hash_address]; configuration != nil; configuration = configuration.next {
		if configuration.curtok == curtok && last_element == configuration.last_element {
			return true
		}
	}
	return false
}

func (this ConfigurationStack) push(stack []int, stack_top int, conflict_index int, curtok int, action_length int) {

	var configuration = NewConfigurationElement()
	var hash_address = curtok % TABLE_SIZE

	configuration.next = this.table[hash_address]

	this.table[hash_address] = configuration
	this.max_configuration_size++ // keep track of int of configurations

	configuration.stack_top = stack_top
	this.stacks_size += stack_top + 1 // keep track of int of stack elements processed
	configuration.last_element = this.findOrInsertStack(this.state_root, stack, 0, stack_top)
	configuration.conflict_index = conflict_index
	configuration.curtok = curtok
	configuration.action_length = action_length

	this.configuration_stack.add(configuration)
	return
}
func (a *ConfigurationStack) cast(c interface{}) *ConfigurationElement {
	configuration, ok := c.(*ConfigurationElement)
	if ok {
		return configuration
	}
	return nil
}
func (this ConfigurationStack) pop() *ConfigurationElement {
	var configuration *ConfigurationElement = nil
	if this.configuration_stack.size() > 0 {
		var index = this.configuration_stack.size() - 1
		configuration = this.cast(this.configuration_stack.get(index))
		if configuration != nil {
			configuration.act = this.prs.baseAction(configuration.conflict_index)
			configuration.conflict_index += 1
			if this.prs.baseAction(configuration.conflict_index) == 0 {
				this.configuration_stack.reset(index)
			}
		} else {
			return nil
		}

	}

	return configuration
}
func (this ConfigurationStack) top() *ConfigurationElement {
	var configuration *ConfigurationElement = nil
	if this.configuration_stack.size() > 0 {
		var index = this.configuration_stack.size() - 1
		var configuration = this.cast(this.configuration_stack.get(index))
		configuration.act = this.prs.baseAction(configuration.conflict_index)
	}

	return configuration
}
func (this ConfigurationStack) size() int {
	return this.configuration_stack.size()
}
func (this ConfigurationStack) maxConfigurationSize() int {
	return this.max_configuration_size
}
func (this ConfigurationStack) numStateElements() int {
	return this.state_element_size
}
func (this ConfigurationStack) stacksSize() int {
	return this.stacks_size
}
