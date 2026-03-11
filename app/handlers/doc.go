package handlers

import (
	"context"
	"github.com/desygner/desygner-toolchain/app/models"
)

type Doc struct {
}

func (receiver *Doc) Search(ctx context.Context, keyword string) (res *models.Doc, err error) {
	
}
