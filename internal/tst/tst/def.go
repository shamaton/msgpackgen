package tst2

type B struct {
}

type NotGeneratedChild struct {
	Interface interface{}
}

/*

	if v.Ints != nil{
		vp := *v.Ints
		if vp != nil {
			vpp := *vp
			if vpp != nil {
				offset = encoder.WriteSliceLength(len(vpp), offset, false)
				for _, vppv := range vpp {
					if vppv != nil {
						vppvp := *vppv
						if vppvp != nil {
							vppvpp := *vppvp
							offset = encoder.WriteInt(vppvpp, offset)
						}else {
							offset = encoder.WriteNil(offset)
						}
					}else {
						offset = encoder.WriteNil(offset)
					}
				}
			}else {
				offset = encoder.WriteNil(offset)
			}
		}else {
			offset = encoder.WriteNil(offset)
		}
	}


	if v.Ints != nil {
		vp := * v.Ints
		if vp != nil {
			vpp := * vp
			if vpp != nil {
				offset = encoder . WriteSliceLength (len(vpp),offset,false)
				for _,vppv := range vpp {
					if vppv != nil {
						vppvp := * vppv
						if vppvp != nil {
							vppvpp := * vppvp
							offset = encoder . WriteInt (vppvpp,offset)
						} else {
							offset = encoder . WriteNil (offset)
						}
					} else {
						offset = encoder . WriteNil (offset)
					}
				}
			} else {
				offset = encoder . WriteNil (offset)
			}
		} else {
			offset = encoder . WriteNil (offset)
		}
	} else {
		offset = encoder . WriteNil (offset)
	}
*/
