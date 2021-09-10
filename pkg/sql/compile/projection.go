package compile

import (
	"matrixone/pkg/sql/colexec/extend"
	vprojection "matrixone/pkg/sql/colexec/projection"
	"matrixone/pkg/sql/op/projection"
	"matrixone/pkg/vm"
)

func (c *compile) compileOutput(o *projection.Projection, mp map[string]uint64) ([]*Scope, error) {
	refer := make(map[string]uint64)
	{
		for i, e := range o.Es {
			if name, ok := e.E.(*extend.Attribute); ok {
				mp[name.Name]++
				continue
			}
			attr := o.As[i]
			if v, ok := mp[attr]; ok {
				refer[attr] = v
				delete(mp, attr)
			} else {
				refer[attr]++
				IncRef(e.E, mp)
			}
		}
	}
	ss, err := c.compile(o.Prev, mp)
	if err != nil {
		return nil, err
	}
	arg := &vprojection.Argument{
		Attrs: make([]string, len(o.Es)),
		Es:    make([]extend.Extend, len(o.Es)),
	}
	{
		for i, e := range o.Es {
			arg.Es[i] = e.E
			arg.Attrs[i] = e.Alias
		}
	}
	arg.Refer = refer
	if o.IsPD {
		for i, s := range ss {
			ss[i] = pushProjection(s, arg)
		}
	} else {
		for i, s := range ss {
			ss[i].Ins = append(s.Ins, vm.Instruction{
				Arg: arg,
				Op:  vm.Projection,
			})
		}
	}
	return ss, nil
}

func (c *compile) compileProjection(o *projection.Projection, mp map[string]uint64) ([]*Scope, error) {
	refer := make(map[string]uint64)
	{
		for i, e := range o.Es {
			if _, ok := e.E.(*extend.Attribute); ok {
				continue
			}
			attr := o.As[i]
			if v, ok := mp[attr]; ok {
				refer[attr] = v
				delete(mp, attr)
			} else {
				refer[attr]++
				IncRef(e.E, mp)
			}
		}
	}
	ss, err := c.compile(o.Prev, mp)
	if err != nil {
		return nil, err
	}
	arg := &vprojection.Argument{
		Attrs: make([]string, len(o.Es)),
		Es:    make([]extend.Extend, len(o.Es)),
	}
	{
		for i, e := range o.Es {
			arg.Es[i] = e.E
			arg.Attrs[i] = e.Alias
		}
	}
	arg.Refer = refer
	if o.IsPD {
		for i, s := range ss {
			ss[i] = pushProjection(s, arg)
		}
	} else {
		for i, s := range ss {
			ss[i].Ins = append(s.Ins, vm.Instruction{
				Arg: arg,
				Op:  vm.Projection,
			})
		}
	}
	return ss, nil
}

func pushProjection(s *Scope, arg *vprojection.Argument) *Scope {
	if s.Magic == Merge || s.Magic == Remote {
		for i := range s.Ss {
			s.Ss[i] = pushProjection(s.Ss[i], arg)
		}
	} else {
		n := len(s.Ins) - 1
		s.Ins = append(s.Ins, vm.Instruction{
			Arg: arg,
			Op:  vm.Projection,
		})
		s.Ins[n], s.Ins[n+1] = s.Ins[n+1], s.Ins[n]
	}
	return s
}
