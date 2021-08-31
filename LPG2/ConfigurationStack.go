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

	a.state_root.number = prs.getStartState()

	a.table = make([]*ConfigurationElement, a.TABLE_SIZE, a.TABLE_SIZE)
	a.configuration_stack = NewObjectTupleWithEstimate(1 << 12)
	a.max_configuration_size = 0
	a.stacks_size = 0
	return a
}

func (self ConfigurationStack) makeStateList(parent *StateElement, stack []int, index int, stack_top int) *StateElement {
	var i int
	for i = index; i <= stack_top; i++ {

		self.state_element_size++

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
func (self ConfigurationStack) findOrInsertStack(root *StateElement, stack []int, index int, stack_top int) *StateElement {
	var state_number = stack[index]
	var p *StateElement
	for p = root; p != nil; p = p.siblings {

		if p.number == state_number {
			if index == stack_top {
				return p
			} else {
				if p.children == nil {
					return self.makeStateList(p, stack, index+1, stack_top)
				} else {
					self.findOrInsertStack(p.children, stack, index+1, stack_top)
				}
			}
		}
	}

	self.state_element_size++

	var node = NewStateElement()
	node.number = state_number
	node.parent = root.parent
	node.children = nil
	node.siblings = root.siblings
	root.siblings = node

	if index == stack_top {
		return node
	} else {
		return self.makeStateList(node, stack, index+1, stack_top)
	}
}

func (self ConfigurationStack) findConfiguration(stack []int, stack_top int, curtok int) bool {

	var last_element = self.findOrInsertStack(self.state_root, stack, 0, stack_top)
	var hash_address = curtok % self.TABLE_SIZE
	var configuration *ConfigurationElement
	for configuration = self.table[hash_address]; configuration != nil; configuration = configuration.next {
		if configuration.curtok == curtok && last_element == configuration.last_element {
			return true
		}
	}
	return false
}

func (self ConfigurationStack) push(stack []int, stack_top int, conflict_index int, curtok int, action_length int) {

	var configuration = NewConfigurationElement()
	var hash_address = curtok % self.TABLE_SIZE

	configuration.next = self.table[hash_address]

	self.table[hash_address] = configuration
	self.max_configuration_size++ // keep track of int of configurations

	configuration.stack_top = stack_top
	self.stacks_size += stack_top + 1 // keep track of int of stack elements processed
	configuration.last_element = self.findOrInsertStack(self.state_root, stack, 0, stack_top)
	configuration.conflict_index = conflict_index
	configuration.curtok = curtok
	configuration.action_length = action_length

	self.configuration_stack.add(configuration)
	return
}
func (a *ConfigurationStack) cast(c interface{}) *ConfigurationElement {
	configuration, ok := c.(*ConfigurationElement)
	if ok {
		return configuration
	}
	return nil
}
func (self ConfigurationStack) pop() *ConfigurationElement {
	var configuration *ConfigurationElement = nil
	if self.configuration_stack.size() > 0 {
		var index = self.configuration_stack.size() - 1
		configuration = self.cast(self.configuration_stack.get(index))
		if configuration != nil {
			configuration.act = self.prs.baseAction(configuration.conflict_index)
			configuration.conflict_index += 1

			if self.prs.baseAction(configuration.conflict_index) == 0 {
				self.configuration_stack.resetTo(index)
			}
		} else {
			return nil
		}

	}

	return configuration
}
func (self ConfigurationStack) top() *ConfigurationElement {
	var configuration *ConfigurationElement = nil
	if self.configuration_stack.size() > 0 {
		var index = self.configuration_stack.size() - 1
		var configuration = self.cast(self.configuration_stack.get(index))
		configuration.act = self.prs.baseAction(configuration.conflict_index)
	}

	return configuration
}
func (self ConfigurationStack) size() int {
	return self.configuration_stack.size()
}
func (self ConfigurationStack) maxConfigurationSize() int {
	return self.max_configuration_size
}
func (self ConfigurationStack) numStateElements() int {
	return self.state_element_size
}
func (self ConfigurationStack) stacksSize() int {
	return self.stacks_size
}
