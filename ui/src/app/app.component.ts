import { ApplicationRef, Component, OnInit, inject } from '@angular/core';
import { ConfigService } from './config.service';
import { DOCUMENT } from '@angular/common';
import { SwUpdate } from '@angular/service-worker';
import { concat, first, interval, startWith } from 'rxjs';
import { ProfileService } from 'src/services/profile.service';
import { Tooltip, Sidenav, Datepicker, Input, initTE } from 'tw-elements';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit {
  private readonly config = inject(ConfigService);
  private readonly document = inject(DOCUMENT);
  private readonly updates = inject(SwUpdate);
  private readonly appRef = inject(ApplicationRef);
  readonly profileService = inject(ProfileService);

  ngOnInit() {
    initTE({ Sidenav, Tooltip, Datepicker, Input });

    if (!!this.config.config.siteName) {
      this.document.querySelector('title')!.innerText =
        this.config.config.siteName;
    }

    if (this.updates.isEnabled) {
      // version updates
      this.updates.versionUpdates.subscribe((evt) => {
        switch (evt.type) {
          case 'VERSION_READY':
            // this.nzMessage.info(`Eine neue Version von CIS wurde installiert. Bitte lade die Seite neu`);
            break;

          case 'VERSION_INSTALLATION_FAILED':
            // this.nzMessage.error(`Failed to install app version '${evt.version.hash}': ${evt.error}`);
            break;
        }
      });

      // Allow the app to stabilize first, before starting
      // polling for updates with `interval()`.
      const appIsStable$ = this.appRef.isStable.pipe(
        first((isStable) => isStable === true)
      );
      const everySixHours$ = interval(60 * 60 * 1000).pipe(startWith(-1));

      const everyHourOnceAppIsStable$ = concat(appIsStable$, everySixHours$);

      everyHourOnceAppIsStable$.subscribe(async () => {
        try {
          const updateFound = await this.updates.checkForUpdate();
          console.log(
            updateFound
              ? 'A new version is available.'
              : 'Already on the latest version.'
          );
        } catch (err) {
          console.error('Failed to check for updates:', err);
        }
      });

      this.updates.unrecoverable.subscribe((event) => {
        /*
        this.nzMessage.error(
          'An error occurred that we cannot recover from:' +
          event.reason +
          ' Please reload the page.'
        );
        */
      });
    }
  }
}
