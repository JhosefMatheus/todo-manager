package project_dto

type CreateProjectDTO struct {
	Name            string
	ParentProjectId *int
}

func (this *CreateProjectDTO) IsInvalid() bool {
	isNameInvalid := this.IsNameInvalid()

	isInvalid := isNameInvalid

	return isInvalid
}

func (this *CreateProjectDTO) IsNameInvalid() bool {
	return len(this.Name) == 0
}
