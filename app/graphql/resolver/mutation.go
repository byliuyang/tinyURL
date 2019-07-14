package resolver

type Mutation struct {
}

type UrlInput struct {
	OriginalUrl *string
	CustomAlias *string
	ExpireAt    *string
}

type CreateUrlArgs struct {
	Url       *UrlInput
	UserEmail *string
}

func (m Mutation) CreateUrl(args *CreateUrlArgs) *Url {
	return &Url{}
}
