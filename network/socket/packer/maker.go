package packer

type Maker interface {
	New() Packer
}

func NewSizeMaker() Maker {
	return &sizePackerMaker{packer: &sizePacker{}}
}

type sizePackerMaker struct {
	packer Packer
}

func (m *sizePackerMaker) New() Packer {
	return m.packer
}
