import { CommonModule } from "@angular/common";
import { Component, OnInit, inject } from "@angular/core";
import { FormsModule } from "@angular/forms";
import { Router } from "@angular/router";
import { ConnectError } from "@bufbuild/connect";
import { AuthType } from "@tkd/apis/gen/es/tkd/idm/v1/auth_service_pb.js";
import { AUTH_SERVICE } from "src/app/clients";
import { ProfileService } from "src/services/profile.service";

@Component({
  standalone: true,
  templateUrl: './login.component.html',
  imports: [
    FormsModule,
    CommonModule
  ]
})
export class LoginComponent implements OnInit {
  client = inject(AUTH_SERVICE);
  profile = inject(ProfileService);
  router = inject(Router);

  username = '';
  password = '';
  loginErrorMessage = '';

  ngOnInit() {

  }

  async submit() {
    try {
      const result = await this.client.login({
        authType: AuthType.PASSWORD,
        auth: {
          case: 'password',
          value: {
            password: this.password,
            username: this.username,
          }
        }
      })

      if (result.response.case === "accessToken") {
        localStorage.setItem("access_token", result.response.value.token);
        /*
        if (result.response.value.redirectTo) {
          window.location(result.response.redirectTo)
        }
        */

      } else {
        throw new Error("unexpected response")
      }


      await this.profile.loadProfile();
      await this.router.navigate(['/profile'])

    } catch(err) {
      const connectErr = ConnectError.from(err);

      this.loginErrorMessage = connectErr.rawMessage;
    }
  }
}
