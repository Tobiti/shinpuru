/** @format */

import { Component } from '@angular/core';
import { APIService } from 'src/app/api/api.service';
import { SpinnerService } from 'src/app/components/spinner/spinner.service';
import { ActivatedRoute } from '@angular/router';
import { Guild, Role, Member } from 'src/app/api/api.models';
import { utils } from 'protractor';
import { permLvlColor } from 'src/app/utils/utils';

@Component({
  selector: 'app-guild',
  templateUrl: './guild.component.html',
  styleUrls: ['./guild.component.sass'],
})
export class GuildComponent {
  public guild: Guild;
  public permLvl: number;
  public members: Member[];

  public permLvlColor = permLvlColor;

  constructor(
    public api: APIService,
    public spinner: SpinnerService,
    private route: ActivatedRoute
  ) {
    const guildID = this.route.snapshot.paramMap.get('id');
    this.api.getGuild(guildID).subscribe((guild) => {
      this.guild = guild;
      this.members = this.guild.members.filter(
        (m) => m.user.id !== this.guild.self_member.user.id
      );
      this.api
        .getPermissionLvl(guildID, guild.self_member.user.id)
        .subscribe((lvl) => {
          this.permLvl = lvl;
        });
      this.spinner.stop('spinner-load-guild');
    });
  }

  public get userRoles(): Role[] {
    const userRoleIDs = this.guild.self_member.roles;
    return this.guild.roles
      .filter((r) => userRoleIDs.includes(r.id))
      .sort((a, b) => b.position - a.position);
  }

  public searchInput(e: any) {
    const val = e.target.value.toLowerCase();

    if (val === '') {
      this.members = this.guild.members.filter(
        (m) => m.user.id !== this.guild.self_member.user.id
      );
    } else {
      this.members = this.guild.members.filter(
        (m) =>
          m.user.id !== this.guild.self_member.user.id &&
          ((m.nick && m.nick.toLowerCase().includes(val)) ||
            m.user.username.toLowerCase().includes(val) ||
            m.user.id.includes(val))
      );
    }
  }
}