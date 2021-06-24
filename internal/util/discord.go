package util

import (
	"fmt"

	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekrotja/discordgo"
)

// GetInviteLink returns the invite link for the bot's
// account with the specified permissions.
func GetInviteLink(s *discordgo.Session) string {
	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=%d",
		s.State.User.ID, static.InvitePermission)
}
