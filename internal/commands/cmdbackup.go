package commands

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekroTJA/shinpuru/internal/core/backup/backupmodels"
	"github.com/zekroTJA/shinpuru/internal/core/database"
	"github.com/zekroTJA/shinpuru/internal/util"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/acceptmsg"
)

const (
	timeFormat = time.RFC1123
)

type CmdBackup struct {
}

func (c *CmdBackup) GetInvokes() []string {
	return []string{"backup", "backups", "bckp", "guildbackup"}
}

func (c *CmdBackup) GetDescription() string {
	return "enable, disable and manage guild backups"
}

func (c *CmdBackup) GetHelp() string {
	return "`backup <enable|disable>` - enable or disable backups for your guild\n" +
		"`backup (list)` - list all saved backups\n" +
		"`backup restore <id>` - restore a backup"
}

func (c *CmdBackup) GetGroup() string {
	return GroupGuildAdmin
}

func (c *CmdBackup) GetDomainName() string {
	return "sp.guild.admin.backup"
}

func (c *CmdBackup) GetSubPermissionRules() []SubPermission {
	return nil
}

func (c *CmdBackup) Exec(args *CommandArgs) error {
	if len(args.Args) > 0 {
		switch strings.ToLower(args.Args[0]) {
		case "e", "enable":
			return c.switchStatus(args, true)
		case "d", "disable":
			return c.switchStatus(args, false)
		case "r", "restore":
			return c.restore(args)
		case "purge", "clear":
			return c.purgeBackupsAccept(args)
		default:
			return c.list(args)
		}
	}
	return c.list(args)
}

func (c *CmdBackup) switchStatus(args *CommandArgs, enable bool) error {
	err := args.CmdHandler.db.SetGuildBackup(args.Guild.ID, enable)
	if err != nil {
		return err
	}

	if enable {
		return util.SendEmbed(args.Session, args.Channel.ID, "Enabled backup for this guild.\nA full guild backup *(incl. Members, Roles, Channels and Guild Settings)* "+
			"will be created every 12 hours. Only 10 backups per guild will be saved, so you will habe the backup files of the last 5 days.", "", static.ColorEmbedGreen).
			DeleteAfter(15 * time.Second).Error()
	}

	return util.SendEmbed(args.Session, args.Channel.ID, "Backup creation disabled.\n"+
		"You will be still have access to created backups and be able to restore them.", "", static.ColorEmbedOrange).
		DeleteAfter(15 * time.Second).Error()
}

func (c *CmdBackup) getBackupsList(args *CommandArgs) ([]*backupmodels.Entry, string, error) {
	backups, err := args.CmdHandler.db.GetBackups(args.Guild.ID)
	if err != nil && database.IsErrDatabaseNotFound(err) {
		return nil, "", err
	}

	strBackupAll := "*no backups saved*"

	if len(backups) > 0 {
		sort.Slice(backups, func(i, j int) bool {
			return backups[i].Timestamp.Before(backups[j].Timestamp)
		})

		if len(backups) > 10 {
			backups = backups[0:10]
		}

		strBackups := make([]string, len(backups))

		for i, b := range backups {
			strBackups[i] = fmt.Sprintf("`%d` - %s - (ID: `%s`)", i, b.Timestamp.Format(timeFormat), b.FileID)
		}

		strBackupAll = strings.Join(strBackups, "\n")
	}

	return backups, strBackupAll, nil
}

func (c *CmdBackup) list(args *CommandArgs) error {
	status, err := args.CmdHandler.db.GetGuildBackup(args.Guild.ID)
	if err != nil && database.IsErrDatabaseNotFound(err) {
		return err
	}

	strStatus := ":x:  Backups **disabled**"
	if status {
		strStatus = ":white_check_mark:  Backups **enabled**"
	}

	_, strBackupAll, err := c.getBackupsList(args)
	if err != nil {
		return err
	}

	emb := &discordgo.MessageEmbed{
		Color:       static.ColorEmbedDefault,
		Title:       "Backups",
		Description: strStatus,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Saved Backups",
				Value: strBackupAll,
			},
		},
	}

	_, err = args.Session.ChannelMessageSendEmbed(args.Channel.ID, emb)
	return err
}

func (c *CmdBackup) restore(args *CommandArgs) error {
	if len(args.Args) < 2 {
		return util.SendEmbedError(args.Session, args.Channel.ID, "Please specify the index or the ID of the backup, you want to restore.").
			DeleteAfter(8 * time.Second).Error()
	}

	backups, _, err := c.getBackupsList(args)
	if err != nil {
		return err
	}

	i, err := strconv.ParseInt(args.Args[1], 10, 64)
	if err != nil {
		return err
	}

	if i < 0 {
		return util.SendEmbedError(args.Session, args.Channel.ID, "Argument must be an index between 0 and 9 or a snowflake ID.").
			DeleteAfter(8 * time.Second).Error()
	}

	var backup *backupmodels.Entry

	if i < 10 {
		if int64(len(backups)-1) < i {
			return util.SendEmbedError(args.Session, args.Channel.ID,
				fmt.Sprintf("There are only %d (index 0 to %d) backups you can chose from.", len(backups), len(backups)-1)).
				DeleteAfter(8 * time.Second).Error()
		}
		backup = backups[i]
	} else {
		for _, b := range backups {
			if b.FileID == args.Args[1] {
				backup = b
			}
		}
	}

	if backup == nil {
		return util.SendEmbedError(args.Session, args.Channel.ID,
			fmt.Sprintf("Could not find any backup by this specifier: ```\n%s\n```", args.Args[1])).
			DeleteAfter(8 * time.Second).Error()
	}

	accMsg := &acceptmsg.AcceptMessage{
		Session:        args.Session,
		DeleteMsgAfter: true,
		UserID:         args.User.ID,
		Embed: &discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: fmt.Sprintf(":warning:  **WARNING**  :warning:\n\n"+
				"By pressing :white_check_mark:, the structure of this guild will be **reset** to the selected backup:\n\n"+
				"%s - (ID: `%s`)", backup.Timestamp.Format(timeFormat), backup.FileID),
		},
		DeclineFunc: func(m *discordgo.Message) {
			util.SendEmbedError(args.Session, args.Channel.ID, "Canceled.").
				DeleteAfter(6 * time.Second).Error()
			return
		},
		AcceptFunc: func(m *discordgo.Message) {
			c.proceedRestore(args, backup.FileID)
		},
	}

	_, err = accMsg.Send(args.Channel.ID)
	return err
}

func (c *CmdBackup) proceedRestore(args *CommandArgs, fileID string) {
	statusChan := make(chan string)
	errorsChan := make(chan error)

	statusMsg, _ := args.Session.ChannelMessageSendEmbed(args.Channel.ID,
		&discordgo.MessageEmbed{
			Color:       static.ColorEmbedGray,
			Description: "initializing backup restoring...",
		})

	if statusMsg != nil {
		go func() {
			for {
				select {
				case status, ok := <-statusChan:
					if !ok {
						continue
					}
					args.Session.ChannelMessageEditEmbed(statusMsg.ChannelID, statusMsg.ID, &discordgo.MessageEmbed{
						Color:       static.ColorEmbedGray,
						Description: status + "...",
					})
				case err, ok := <-errorsChan:
					if !ok || err == nil {
						continue
					}
					util.SendEmbedError(args.Session, args.Channel.ID,
						"An unexpected error occured while restoring backup (process will not be aborted): ```\n"+err.Error()+"\n```")
				}
			}
		}()
	}

	err := args.CmdHandler.bck.RestoreBackup(args.Guild.ID, fileID, statusChan, errorsChan)
	if err != nil {
		util.SendEmbedError(args.Session, args.Channel.ID,
			fmt.Sprintf("An unexpected error occured while restoring backup: ```\n%s\n```", err.Error()))
	}
}

func (c *CmdBackup) purgeBackupsAccept(args *CommandArgs) error {
	_, err := acceptmsg.New().
		WithSession(args.Session).
		WithEmbed(&discordgo.MessageEmbed{
			Color: static.ColorEmbedOrange,
			Description: ":warning:  **WARNING**  :warning:\n\n" +
				"Do you really want to **purge __all__ backups** for this guild?",
		}).
		LockOnUser(args.User.ID).
		DeleteAfterAnswer().
		DoOnDecline(func(_ *discordgo.Message) {
			util.SendEmbedError(args.Session, args.Channel.ID, "Canceled.").
				DeleteAfter(6 * time.Second).Error()
			return
		}).
		DoOnAccept(func(_ *discordgo.Message) {
			c.purgeBackups(args)
		}).
		Send(args.Channel.ID)

	return err
}

func (c *CmdBackup) purgeBackups(args *CommandArgs) {
	backups, err := args.CmdHandler.db.GetBackups(args.Guild.ID)
	if err != nil {
		util.SendEmbedError(args.Session, args.Channel.ID,
			fmt.Sprintf("Failed getting backups: ```\n%s\n```", err.Error())).
			DeleteAfter(15 * time.Second).Error()
		return
	}

	var lnBackups = len(backups)
	if lnBackups < 1 {
		util.SendEmbedError(args.Session, args.Channel.ID,
			"There are no backups saved to be purged.").
			DeleteAfter(8 * time.Second).Error()
		return
	}

	var success int
	for _, backup := range backups {
		if err = args.CmdHandler.db.DeleteBackup(args.Guild.ID, backup.FileID); err != nil {
			continue
		}
		if err = args.CmdHandler.st.DeleteObject(static.StorageBucketBackups, backup.FileID); err != nil {
			continue
		}
		success++
	}

	if success < lnBackups {
		util.SendEmbedError(args.Session, args.Channel.ID,
			fmt.Sprintf("Successfully purged `%d` of `%d` backups.\n`%d` backup purges failed.",
				success, lnBackups, lnBackups-success)).
			DeleteAfter(8 * time.Second).Error()
		return
	}

	util.SendEmbed(args.Session, args.Channel.ID,
		fmt.Sprintf("Successfully purged `%d` of `%d` backups.",
			success, lnBackups), "", static.ColorEmbedGreen).
		DeleteAfter(8 * time.Second).Error()
	return
}
