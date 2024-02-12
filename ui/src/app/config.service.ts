import { Injectable } from '@angular/core';

export interface RemoteConfig {
  domain: string;
  loginURL: string;
  siteName: string;
  siteNameUrl: string;
  logoURL: string;
  registration: 'public' | 'token' | 'disabled';
  userAddresses: boolean;
  phoneNumbers: boolean;
  userNameChange: boolean;
}

@Injectable({ providedIn: 'root' })
export class ConfigService {
  static Config: RemoteConfig;

  get config(): RemoteConfig {
    return ConfigService.Config;
  }
}
