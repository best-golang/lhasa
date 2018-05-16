import { Component, OnInit } from '@angular/core';
import { ApplicationsStoreService } from './stores/applications-store.service';
import { DataApplicationService } from './services/data-application-version.service';
import { EnvironmentBean } from './models/commons/applications-bean';
import { ContentListResponse } from './models/commons/entity-bean';
import { Router } from '@angular/router';
import { EnvironmentsStoreService, LoadEnvironmentsAction } from './stores/environments-store.service';
import { DataEnvironmentService } from './services/data-environment.service';


import { TranslateService } from '@ngx-translate/core';
import { UiKitMenuItem } from './models/kit/navbar';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'app';

  public items: UiKitMenuItem[];

  constructor(
    private router: Router,
    private applicationsStoreService: ApplicationsStoreService,
    private applicationsService: DataApplicationService,
    private environmentsStoreService: EnvironmentsStoreService,
    private environmentService: DataEnvironmentService,
    private translate: TranslateService
  ) {
    // this language will be used as a fallback when a translation isn't found in the current language
    this.translate.setDefaultLang('en');

    // the lang to use, if the lang isn't available, it will use the current loader to get them
    this.translate.use('en');

    // Simple menu model
    this.items = [
      {
        id: 'domains',
        label: 'DOMAINS',
        routerLink: '/domains'
      },
      {
        id: 'applications',
        label: 'APPLICATIONS',
        routerLink: '/applications'
      }
    ];
    // the lang to use, if the lang isn't available, it will use the current loader to get them
    this.translate.use('en');
  }

  ngOnInit(): void {
    this.environmentService.GetAllFromContent('/', null).subscribe(
      (value: ContentListResponse<EnvironmentBean>) => {
        const environmentMap = new Map<string, EnvironmentBean>();
        value.content.forEach(env => {
          environmentMap[env.slug] = env;
        });
        this.environmentsStoreService.dispatch(
          new LoadEnvironmentsAction(
            environmentMap
          )
        );
      }
    );
  }
}
