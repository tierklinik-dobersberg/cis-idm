import { AsyncPipe, NgFor, NgForOf, NgIf, NgSwitch, NgSwitchCase, NgTemplateOutlet } from "@angular/common";
import { ChangeDetectionStrategy, Component, OnInit, inject } from "@angular/core";
import { TkdButtonDirective } from "src/app/components/button";
import { ConfigService, FieldConfig } from "src/app/config.service";
import { ProfileService } from "src/services/profile.service";
import { ComplexFieldsPipe, FieldPathPipe, FieldValuePipe, SimpleFieldsPipe } from "./field.pipe";
import { TkdSettingsInputComponent } from "./settings-input/settings-input.component";
import { FormsModule } from "@angular/forms";
import { Profile } from "@tierklinik-dobersberg/apis";
import { take } from "rxjs";
import { USER_SERVICE } from "src/app/clients";
import { Value, JsonValue, Struct } from '@bufbuild/protobuf';

@Component({
  templateUrl: './settings.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  standalone: true,
  imports: [
    NgFor,
    NgForOf,
    NgIf,
    NgSwitch,
    NgSwitchCase,
    NgTemplateOutlet,
    AsyncPipe,
    TkdButtonDirective,
    SimpleFieldsPipe,
    ComplexFieldsPipe,
    FieldPathPipe,
    FieldValuePipe,
    TkdSettingsInputComponent,
    FormsModule,
  ],
})
export class SettingsPageComponent {
  readonly config = inject(ConfigService).config;
  profile$ = inject(ProfileService).profile;

  private readonly profileService = inject(ProfileService);
  private readonly userService = inject(USER_SERVICE);

  updateValue(value: JsonValue, path: string) {
    this.profileService
      .profile
      .pipe(take(1))
      .subscribe(async (profile) => {
        if (!profile?.user) {
          return
        }

        try {
          await this.userService.setUserExtraKey({
            userId: profile.user.id,
            path: path,
            value: Value.fromJson(value),
          })
        } catch(err) {
          console.error(err);
        }

        this.profileService.loadProfile();
      });
  }
}
