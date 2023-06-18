import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { AUTH_SERVICE, SELF_SERVICE, TRANSPORT, authServiceFactory, selfServiceFactory, transportFactory } from './clients';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [
    {
      provide: TRANSPORT,
      useFactory: transportFactory,
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
