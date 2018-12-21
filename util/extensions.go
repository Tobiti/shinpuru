package util

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func EnsureNotEmpty(str, def string) string {
	if str == "" {
		return def
	}
	return str
}

func DeleteMessageLater(s *discordgo.Session, msg *discordgo.Message, duration time.Duration) {
	if msg == nil {
		return
	}
	time.AfterFunc(duration, func() {
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	})
}

func FetchRole(s *discordgo.Session, guildID, resolvable string) (*discordgo.Role, error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	rx := regexp.MustCompile("<@&|>")
	resolvable = rx.ReplaceAllString(resolvable, "")

	checkFuncs := []func(*discordgo.Role, string) bool{
		func(r *discordgo.Role, resolvable string) bool {
			return r.ID == resolvable
		},
		func(r *discordgo.Role, resolvable string) bool {
			return r.Name == resolvable
		},
		func(r *discordgo.Role, resolvable string) bool {
			return strings.ToLower(r.Name) == strings.ToLower(resolvable)
		},
		func(r *discordgo.Role, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
		func(r *discordgo.Role, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.Name), strings.ToLower(resolvable))
		},
	}

	for _, checkFunc := range checkFuncs {
		for _, r := range guild.Roles {
			if checkFunc(r, resolvable) {
				return r, nil
			}
		}
	}

	return nil, errors.New("could not be fetched")
}

func FetchMember(s *discordgo.Session, guildID, resolvable string) (*discordgo.Member, error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, err
	}
	rx := regexp.MustCompile("<@|!|>")
	resolvable = rx.ReplaceAllString(resolvable, "")

	checkFuncs := []func(*discordgo.Member, string) bool{
		func(r *discordgo.Member, resolvable string) bool {
			return r.User.ID == resolvable
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.User.Username == resolvable
		},
		func(r *discordgo.Member, resolvable string) bool {
			return strings.ToLower(r.User.Username) == strings.ToLower(resolvable)
		},
		func(r *discordgo.Member, resolvable string) bool {
			return strings.HasPrefix(strings.ToLower(r.User.Username), strings.ToLower(resolvable))
		},
		func(r *discordgo.Member, resolvable string) bool {
			return strings.Contains(strings.ToLower(r.User.Username), strings.ToLower(resolvable))
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick == resolvable
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.ToLower(r.Nick) == strings.ToLower(resolvable)
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.HasPrefix(strings.ToLower(r.Nick), strings.ToLower(resolvable))
		},
		func(r *discordgo.Member, resolvable string) bool {
			return r.Nick != "" && strings.Contains(strings.ToLower(r.Nick), strings.ToLower(resolvable))
		},
	}

	for _, checkFunc := range checkFuncs {
		for _, m := range guild.Members {
			if checkFunc(m, resolvable) {
				return m, nil
			}
		}
	}

	return nil, errors.New("could not be fetched")
}