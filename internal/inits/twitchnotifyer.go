package inits

import (
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/shinpuru/internal/listeners"
	"github.com/zekroTJA/shinpuru/internal/services/config"
	"github.com/zekroTJA/shinpuru/internal/services/database"
	"github.com/zekroTJA/shinpuru/internal/util/static"
	"github.com/zekroTJA/shinpuru/pkg/twitchnotify"
)

func InitTwitchNotifyWorker(container di.Container) *twitchnotify.NotifyWorker {

	listener := container.Get(static.DiTwitchNotifyListener).(*listeners.ListenerTwitchNotify)
	cfg := container.Get(static.DiConfig).(config.Provider)
	db := container.Get(static.DiDatabase).(database.Database)

	if cfg.Config().TwitchApp.ClientID == "" || cfg.Config().TwitchApp.ClientSecret == "" {
		logrus.Info("twitch app credentials are empty")
		return nil
	}

	tnw, err := twitchnotify.New(
		twitchnotify.Credentials{
			ClientID:     cfg.Config().TwitchApp.ClientID,
			ClientSecret: cfg.Config().TwitchApp.ClientSecret,
		},
		listener.HandlerWentOnline,
		listener.HandlerWentOffline,
		twitchnotify.Config{
			TimerDelay: 0,
		},
	)

	if err != nil {
		logrus.WithError(err).Fatal("twitch app credentials are invalid")
	}

	notifies, err := db.GetAllTwitchNotifies("")
	if err == nil {
		for _, notify := range notifies {
			if u, err := tnw.GetUser(notify.TwitchUserID, twitchnotify.IdentID); err == nil {
				tnw.AddUser(u)
			}
		}
	} else {
		logrus.WithError(err).Fatal("failed getting Twitch notify entreis")
	}

	return tnw
}
