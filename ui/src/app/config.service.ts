import { Injectable } from "@angular/core";

export enum FeatureFlags {
	FeatureAll         = "all",
	FeatureAddresses   = "addresses",
	FeatureEMails      = "emails",
	FeaturePhoneNumbers= "phoneNumbers",
	FeatureEMailInvite = "emailInvite",
	FeatureLoginByMail = "loginByMail",
    FeatureAllowUsernameChange = "allowUsernameChange"
}

export interface RemoteConfig {
    domain: string;
    loginURL: string;
    siteName: string;
    siteNameUrl: string;
    features: {
        [key in FeatureFlags]: boolean
    };
}

@Injectable({providedIn: 'root'})
export class ConfigService {
    static Config: RemoteConfig;

    get config(): RemoteConfig {
        return ConfigService.Config;
    }
}