package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/commands"
	"github.com/zekroTJA/shinpuru/internal/core/config"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/core/lctimer"
	"github.com/zekroTJA/shinpuru/internal/core/storage"
	"github.com/zekroTJA/shinpuru/internal/core/twitchnotify"
	"github.com/zekroTJA/shinpuru/internal/util"
)

func InitCommandHandler(s *discordgo.Session, cfg *config.Config, db database.Database, st storage.Storage, tnw *twitchnotify.NotifyWorker, lct *lctimer.LifeCycleTimer) *commands.CmdHandler {
	cmdHandler := commands.NewCmdHandler(s, db, st, cfg, tnw, lct)

	cmdHandler.RegisterCommand(&commands.CmdHelp{})
	cmdHandler.RegisterCommand(&commands.CmdPrefix{})
	cmdHandler.RegisterCommand(&commands.CmdPerms{})
	cmdHandler.RegisterCommand(&commands.CmdClear{})
	cmdHandler.RegisterCommand(&commands.CmdMvall{})
	cmdHandler.RegisterCommand(&commands.CmdInfo{})
	cmdHandler.RegisterCommand(&commands.CmdSay{})
	cmdHandler.RegisterCommand(&commands.CmdQuote{})
	cmdHandler.RegisterCommand(&commands.CmdGame{})
	cmdHandler.RegisterCommand(&commands.CmdAutorole{})
	cmdHandler.RegisterCommand(&commands.CmdReport{})
	cmdHandler.RegisterCommand(&commands.CmdModlog{})
	cmdHandler.RegisterCommand(&commands.CmdKick{})
	cmdHandler.RegisterCommand(&commands.CmdBan{})
	cmdHandler.RegisterCommand(&commands.CmdVote{})
	cmdHandler.RegisterCommand(&commands.CmdProfile{})
	cmdHandler.RegisterCommand(&commands.CmdId{})
	cmdHandler.RegisterCommand(&commands.CmdMute{})
	cmdHandler.RegisterCommand(&commands.CmdMention{})
	cmdHandler.RegisterCommand(&commands.CmdNotify{})
	cmdHandler.RegisterCommand(&commands.CmdVoicelog{})
	cmdHandler.RegisterCommand(&commands.CmdBug{})
	cmdHandler.RegisterCommand(&commands.CmdStats{})
	cmdHandler.RegisterCommand(&commands.CmdTwitchNotify{})
	cmdHandler.RegisterCommand(&commands.CmdGhostping{})
	cmdHandler.RegisterCommand(&commands.CmdExec{})
	cmdHandler.RegisterCommand(&commands.CmdBackup{})
	cmdHandler.RegisterCommand(&commands.CmdInviteBlock{})
	cmdHandler.RegisterCommand(&commands.CmdTag{})
	cmdHandler.RegisterCommand(&commands.CmdJoinMsg{})
	cmdHandler.RegisterCommand(&commands.CmdLeaveMsg{})
	cmdHandler.RegisterCommand(&commands.CmdSnowflake{})
	cmdHandler.RegisterCommand(&commands.CmdChannelStats{})

	if util.Release != "TRUE" {
		cmdHandler.RegisterCommand(&commands.CmdTest{})
	}

	util.Log.Infof("%d commands registered", cmdHandler.GetCommandListLen())

	return cmdHandler
}
