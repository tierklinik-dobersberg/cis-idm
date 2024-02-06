import { APP_INITIALIZER, NgModule, isDevMode } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import {
  AUTH_SERVICE,
  NOTIFY_SERVICE,
  SELF_SERVICE,
  TRANSPORT,
  authServiceFactory,
  notifyServiceFactory,
  selfServiceFactory,
  transportFactory,
} from './clients';

import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { map } from 'rxjs';
import { ConfigService, RemoteConfig } from './config.service';
import { ServiceWorkerModule } from '@angular/service-worker';
import { NavigationComponent } from './navigation';
import { TkdLayoutModule } from '@tierklinik-dobersberg/angular/layout';

const loadConfigFactory = (client: HttpClient) => {
  return () =>
    client.get<RemoteConfig>(`/config.json`).pipe(
      map((response) => {
        ConfigService.Config = response;
        console.log('remote configuration loaded', response);
      })
    );
};

@NgModule({
  declarations: [AppComponent, NavigationComponent],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    TkdLayoutModule,
    ServiceWorkerModule.register('ngsw-worker.js', {
      enabled: !isDevMode(),
      // Register the ServiceWorker as soon as the application is stable
      // or after 30 seconds (whichever comes first).
      registrationStrategy: 'registerWhenStable:30000',
    }),
  ],
  providers: [
    {
      provide: APP_INITIALIZER,
      multi: true,
      useFactory: loadConfigFactory,
      deps: [HttpClient],
    },
    {
      provide: TRANSPORT,
      useFactory: transportFactory,
      deps: [ActivatedRoute, Router],
    },
    {
      provide: AUTH_SERVICE,
      useFactory: authServiceFactory,
      deps: [TRANSPORT],
    },
    {
      provide: SELF_SERVICE,
      useFactory: selfServiceFactory,
      deps: [TRANSPORT],
    },
    {
      provide: NOTIFY_SERVICE,
      useFactory: notifyServiceFactory,
      deps: [TRANSPORT],
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
