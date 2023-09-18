import { APP_INITIALIZER, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { AUTH_SERVICE, SELF_SERVICE, TRANSPORT, authServiceFactory, selfServiceFactory, transportFactory } from './clients';

import { HttpClient, HttpClientModule } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { map } from 'rxjs';
import { ConfigService, RemoteConfig } from './config.service';

const loadConfigFactory = (client: HttpClient) => {
  return () => client.get<RemoteConfig>(`/config.json`)
    .pipe(
      map(response => {
        ConfigService.Config = response;
        console.log('remote configuration loaded', response)
      })
    )
}

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule
  ],
  providers: [
    {
      provide: APP_INITIALIZER,
      multi: true,
      useFactory: loadConfigFactory,
      deps: [HttpClient]
    },
    {
      provide: TRANSPORT,
      useFactory: transportFactory,
      deps: [ActivatedRoute, Router]
    },
    {
      provide: AUTH_SERVICE,
      useFactory: authServiceFactory,
      deps: [TRANSPORT]
    },
    {
      provide: SELF_SERVICE,
      useFactory: selfServiceFactory,
      deps: [TRANSPORT]
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
