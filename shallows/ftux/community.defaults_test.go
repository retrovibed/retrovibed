package ftux_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/retrovibed/retrovibed/retroapi/userx"
	"github.com/retrovibed/retrovibed/shallows/communityapi"
	"github.com/retrovibed/retrovibed/shallows/ftux"
	"github.com/stretchr/testify/require"
)

func TestPrepareDefaultCommunities(t *testing.T) {
	t.Run("returns the curated defaults and caches them locally", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())

		suggestions, err := ftux.PrepareDefaultCommunities()
		require.NoError(t, err)
		require.Len(t, suggestions, 2)

		for _, s := range suggestions {
			require.NotEmpty(t, s.Id)
			require.NotEmpty(t, s.Description)
			require.NotEmpty(t, s.Url)
		}

		require.FileExists(t, userx.DefaultConfigDir(userx.DefaultRelRoot(), "default.communities.json"))
	})

	t.Run("reads back locally edited suggestions instead of regenerating", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())

		_, err := ftux.PrepareDefaultCommunities()
		require.NoError(t, err)

		edited, err := json.Marshal([]*communityapi.Community{
			{Id: "edited-id", Description: "edited", Url: "https://edited.example"},
		})
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(userx.DefaultConfigDir(userx.DefaultRelRoot(), "default.communities.json"), edited, 0600))

		suggestions, err := ftux.PrepareDefaultCommunities()
		require.NoError(t, err)
		require.Len(t, suggestions, 1)
		require.Equal(t, "edited-id", suggestions[0].Id)
	})
}
