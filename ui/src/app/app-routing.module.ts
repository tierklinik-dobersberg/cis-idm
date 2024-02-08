import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard, notLoggedInGuard } from 'src/guards/auth.guard';

const routes: Routes = [
  { path: "login", canActivate: [notLoggedInGuard], loadComponent: () => import('./pages/login/login.component').then(m => m.LoginComponent) },
  { path: "refresh", loadComponent: () => import('./pages/refresh/refresh.component').then(m => m.RefreshComponent) },
  { path: "registration", canActivate: [notLoggedInGuard], loadComponent: () => import('./pages/registration/registration.component').then(m => m.RegistrationComponent) },
  { path: "password/request-reset", canActivate: [notLoggedInGuard], loadComponent: () => import('./pages/reset-password/reset-password.component').then(m => m.ResetPasswordComponent) },
  { path: "password/reset", canActivate: [notLoggedInGuard], loadComponent: () => import('./pages/reset-password/reset-password.component').then(m => m.ResetPasswordComponent) },
  { path: "profile", canActivate: [authGuard], loadComponent: () => import('./pages/profile/profile.component').then(m => m.ProfileComponent )},
  { path: "profile/change-password", canActivate: [authGuard], loadComponent: () => import('./pages/change-password/change-password.component').then(m => m.ChangePasswordComponent)},
  { path: "profile/edit", canActivate: [authGuard], loadComponent: () => import('./pages/edit-profile/edit-profile.component').then(m => m.EditProfileComponent)},
  { path: "profile/edit-address", canActivate: [authGuard], loadComponent: () => import('./pages/add-edit-address/add-edit-address.component').then(m => m.AddEditAddressComponent)},
  { path: "profile/edit-address/:id", canActivate: [authGuard], loadComponent: () => import('./pages/add-edit-address/add-edit-address.component').then(m => m.AddEditAddressComponent)},
  { path: "profile/edit-mail", canActivate: [authGuard], loadComponent: () => import('./pages/add-edit-mail/add-edit-mail.component').then(m => m.AddEditMailComponent)},
  { path: "profile/edit-mail/:id", canActivate: [authGuard], loadComponent: () => import('./pages/add-edit-mail/add-edit-mail.component').then(m => m.AddEditMailComponent)},
  { path: "profile/edit-phone", canActivate: [authGuard], loadComponent: () => import('./pages/add-edit-phone/add-edit-phone.component').then(m => m.AddEditPhoneComponent)},
  { path: "profile/edit-phone/:id", canActivate: [authGuard], loadComponent: () => import('./pages/add-edit-phone/add-edit-phone.component').then(m => m.AddEditPhoneComponent)},
  { path: "profile/verify-phone/:id", canActivate: [authGuard], loadComponent: () => import('./pages/verify-phone/verify-phone.component').then(m => m.VerifyPhoneComponent)},
  { path: "profile/verify-mail", canActivate: [authGuard], loadComponent: () => import('./pages/verify-mail/verify-mail.component').then(m => m.VerifyMailComponent) },
  { path: "security", canActivate: [authGuard], loadComponent: () => import('./pages/security-overview/security-overview.component').then(m => m.SecurityOverviewComponent)},
  { path: "profile/edit-avatar", canActivate: [authGuard], loadComponent: () => import('./pages/edit-avatar/edit-avatar.component').then(m => m.EditAvatarComponent)},
  { path: "**", redirectTo: '/login' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

