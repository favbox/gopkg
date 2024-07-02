package filex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNameFromURL(t *testing.T) {
	url := "https://qhtbdoss.kujiale.com/fpimgnew/prod/bim/3FO3QX0J81W2/1/MY3AJ5YKTJX7AAABAAAAADY8.jpg"
	name, err := GetNameFromURL(url)
	assert.Nil(t, err)
	assert.Equal(t, "MY3AJ5YKTJX7AAABAAAAADY8.jpg", name)
}
