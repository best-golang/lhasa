import { Injectable } from '@angular/core';
import { ConfigurationService } from './configuration.service';
import { DefaultResource } from '../interfaces/default-resources.interface';
import { DataCoreResource } from './data-core-resources.service';
import { BitbucketBean } from '../models/commons/applications-bean';
import { HttpClient } from '@angular/common/http';

/**
 * data model
 */

@Injectable()
export class BitbucketService extends DataCoreResource<BitbucketBean> implements DefaultResource<BitbucketBean> {
  constructor(
    private _http: HttpClient,
    private _configuration: ConfigurationService
  ) {
    super(_configuration, _configuration.ServerCatalogCompanionWithApiUrl + 'v1/bitbucket', _http);
  }
}
