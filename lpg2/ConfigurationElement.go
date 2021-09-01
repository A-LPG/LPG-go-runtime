package lpg2



type ConfigurationElement struct {
     next *ConfigurationElement
     last_element *StateElement
     stack_top int
     action_length int
     conflict_index int
     curtok int
     act int
}


 func  NewConfigurationElement() *ConfigurationElement{
        a := new(ConfigurationElement)
        a.next  = nil
        a.last_element  = nil
        a.stack_top  = 0
        a.action_length  = 0
        a.conflict_index  = 0
        a.curtok  = 0
        a.act  = 0
        return a
}
func (a *ConfigurationElement) RetrieveStack(stack []int) {
    var tail = a.last_element
    var i  int
    for  i  = a.stack_top;i >= 0; i-- {
        if nil == tail {
            return
        }
        stack[i] = tail.number
        tail = tail.parent
    }
    return
}
