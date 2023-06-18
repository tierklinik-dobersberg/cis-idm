import { Injectable, inject } from "@angular/core";
import { ConnectError } from "@bufbuild/connect";
import { Profile } from "@tkd/apis/gen/es/tkd/idm/v1/user_pb";
import { BehaviorSubject, Observable, filter } from "rxjs";
import { AUTH_SERVICE } from "src/app/clients";

@Injectable({providedIn: 'root'})
export class ProfileService {
  client = inject(AUTH_SERVICE);
  _ready$ = new BehaviorSubject<boolean |null>(null);
  profile = new BehaviorSubject<Profile | null>(null);

  get ready(): Observable<boolean> {
    return this._ready$.asObservable()
      .pipe(
        filter(value => value !== null)
      ) as any;
  }

  constructor() {
    console.log("loading profile")
    this.loadProfile()
  }

  async loadProfile() {
    try {
      const result = await this.client.introspect({})
      this.profile.next(result.profile!);
      console.log("got response", result)
    } catch(err) {
      console.error(ConnectError.from(err))
      this.profile.next(null)
    }

    this._ready$.next(true);
  }
}
