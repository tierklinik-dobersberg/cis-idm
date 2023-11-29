import { Injectable } from "@angular/core";

export enum FeatureFlags {
	FeatureAll         = "all",
	FeatureAddresses   = "addresses",
	FeatureEMails      = "emails",
	FeaturePhoneNumbers= "phoneNumbers",
	FeatureEMailInvite = "emailInvite",
	FeatureLoginByMail = "loginByMail",
  FeatureAllowUsernameChange = "allowUsernameChange",
  FeatureSelfRegistration = "registration"
}

export interface RemoteConfig {
    domain: string;
    loginURL: string;
    siteName: string;
    siteNameUrl: string;
    registrationRequiresToken: boolean;
    features: {
        [key in FeatureFlags]: boolean
    };
    logoURL: string;
}

@Injectable({providedIn: 'root'})
export class ConfigService {
    static Config: RemoteConfig;

    get config(): RemoteConfig {
        return ConfigService.Config;
    }
}
