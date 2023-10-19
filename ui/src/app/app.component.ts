import { Component, OnInit, inject } from '@angular/core';
import { ConfigService } from './config.service';
import { DOCUMENT } from '@angular/common';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  private readonly config = inject(ConfigService);
  private readonly document = inject(DOCUMENT);

  ngOnInit() {
    if (!!this.config.config.siteName) {
      this.document.querySelector('title')!.innerText = this.config.config.siteName;
    }
  }
}
