/** @format */

import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AntiraidSettings, JoinlogEntry } from 'src/app/api/api.models';
import { APIService } from 'src/app/api/api.service';
import { ToastService } from 'src/app/components/toast/toast.service';
import { format } from 'date-fns';
import { TIME_FORMAT } from 'src/app/utils/consts';

@Component({
  selector: 'app-ga-antiraid',
  templateUrl: './ga-antiraid.component.html',
  styleUrls: ['./ga-antiraid.component.scss'],
})
export class GuildAdminAntiraidComponent implements OnInit {
  public antiraidSettings: AntiraidSettings;
  public joinlog: JoinlogEntry[] = [];
  private guildID: string;

  public dateFormat = (d: string | Date, f = TIME_FORMAT) =>
    format(new Date(d), f);

  constructor(
    private route: ActivatedRoute,
    private api: APIService,
    private toasts: ToastService
  ) {}

  ngOnInit() {
    this.route.params.subscribe(async (params) => {
      this.guildID = params.guildid;

      this.antiraidSettings = await this.api
        .getGuildSettingsAntiraid(this.guildID)
        .toPromise();

      this.fetchJoinlog();
    });
  }

  public async fetchJoinlog() {
    try {
      const res = await this.api
        .getGuildAntiraidJoinlog(this.guildID)
        .toPromise();
      this.joinlog = res.data;
    } catch {}
  }

  public async onSave() {
    try {
      await this.api
        .postGuildSettingsAntiraid(this.guildID, this.antiraidSettings)
        .toPromise();
      this.toasts.push(
        'Antiraid settings saved.',
        'Settings saved',
        'green',
        4000
      );
    } catch {}
  }

  public onDownloadJoinlog() {
    const element = document.createElement('a');
    element.setAttribute('href', this.api.rcGuildAntiraidJoinlog(this.guildID));
    element.setAttribute('download', 'joinlog_export');

    element.style.display = 'none';
    document.body.appendChild(element);

    element.click();

    document.body.removeChild(element);
  }

  public async flushJoinlog() {
    try {
      await this.api.deleteGuildAntiraidJoinlog(this.guildID).toPromise();
      this.joinlog = [];
    } catch {}
  }
}
