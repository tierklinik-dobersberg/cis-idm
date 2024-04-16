import { Injectable } from "@angular/core";
import { L10nConfig, L10nProvider, L10nTranslationLoader } from "angular-l10n";
import { Observable, from } from "rxjs";

export const l10nConfig: L10nConfig = {
  format: 'language-region',
  providers: [
    { name: 'app', asset: 'app' }
  ],
  cache: true,
  keySeparator: '.',
  defaultLocale: {
    language: 'de-AT', currency: 'EUR', timeZone: 'Europe/Vienna'
  },
  schema: [
    {
      locale: {
        language: 'de-AT', currency: 'EUR', timeZone: 'Europe/Vienna'
      }
    }
  ]
}

@Injectable()
export class TranslationLoader implements L10nTranslationLoader {
  get(language: string, provider: L10nProvider): Observable<{ [key: string]: any; }> {
    const data = import(`../i18n/${language}/${provider.asset}.json`)

    return from(data);
  }
}

