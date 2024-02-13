import { Injectable } from '@angular/core';

export interface PossibleValue {
  value: string;
  display_name: string;
}

export interface FieldConfig {
  type: 'string' | 'number' | 'bool' | 'object' | 'list' | 'any' | 'date' | 'time';
  name: string;
  visibility: 'public' | 'self' | 'private' | 'authenticated';
  writeable: boolean | null;
  description: string;
  display_name: string;
  property: FieldConfig[] | null;
  element_type: FieldConfig | null;
  possible_value: PossibleValue[] | null;
}

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
  customUserFields: FieldConfig[] | null;
}

@Injectable({ providedIn: 'root' })
export class ConfigService {
  static Config: RemoteConfig;

  get config(): RemoteConfig {
    return ConfigService.Config;
  }
}
